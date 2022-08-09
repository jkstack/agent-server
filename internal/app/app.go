package app

import (
	"context"
	"fmt"
	"net/http"
	"server/internal/agent"
	"server/internal/conf"
	"server/internal/utils"
	"sync"
	"time"

	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/api"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/stat"
	runtime "github.com/jkstack/jkframe/utils"
	"github.com/kardianos/service"
	"github.com/shirou/gopsutil/v3/disk"
	"golang.org/x/time/rate"
)

type handler interface {
	Init(*conf.Configure, *stat.Mgr)
	HandleFuncs() map[string]func(*agent.Agents, *api.Context)
	OnConnect(*agent.Agent)
	OnClose(string)
	OnMessage(*agent.Agent, *anet.Msg)
}

type App struct {
	cfg    *conf.Configure
	stats  *stat.Mgr
	agents *agent.Agents

	// runtime
	connectLock  sync.Mutex
	blocked      bool
	stAgentCount *stat.Counter
	connectLimit *rate.Limiter
}

// New new app
func New(cfg *conf.Configure, version string) *App {
	st := stat.New(5 * time.Second)
	app := &App{
		cfg:    cfg,
		agents: agent.NewAgents(st),
		stats:  st,
		// runtime
		stAgentCount: st.NewCounter("agent_count"),
		connectLimit: rate.NewLimiter(
			rate.Limit(time.Second/time.Duration(cfg.ConnectLimit)), 1),
	}
	go app.limit()
	return app
}

// Start start app
func (app *App) Start(s service.Service) error {
	go func() {
		logging.SetSizeRotate(logging.SizeRotateConfig{
			Dir:         app.cfg.LogDir,
			Name:        "agent-server",
			Size:        int64(app.cfg.LogSize.Bytes()),
			Rotate:      app.cfg.LogRotate,
			WriteStdout: true,
			WriteFile:   true,
		})
		defer logging.Flush()

		defer utils.Recover("service")

		var mods []handler

		for _, mod := range mods {
			mod.Init(app.cfg, app.stats)
			for uri, cb := range mod.HandleFuncs() {
				app.reg(uri, cb)
			}
		}

		http.HandleFunc("/metrics", app.stats.ServeHTTP)
		http.HandleFunc("/ws/agent", func(w http.ResponseWriter, r *http.Request) {
			if !app.connectLimit.Allow() {
				http.Error(w, "rate limit", http.StatusServiceUnavailable)
				return
			}
			onConnect := make(chan *agent.Agent)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				select {
				case cli := <-onConnect:
					for _, mod := range mods {
						mod.OnConnect(cli)
					}
				case <-ctx.Done():
					return
				}
			}()
			cli := app.agent(w, r, onConnect, cancel)
			go func() {
				for {
					select {
					case msg := <-cli.Unknown():
						if msg == nil {
							return
						}
						for _, mod := range mods {
							mod.OnMessage(cli, msg)
						}
					case <-ctx.Done():
						return
					}
				}
			}()
			if cli != nil {
				<-ctx.Done()
				app.stAgentCount.Dec()
				logging.Info("agent %s connection closed", cli.ID())
				for _, mod := range mods {
					mod.OnClose(cli.ID())
				}
			}
		})

		logging.Info("http listen on %d", app.cfg.Listen)
		runtime.Assert(http.ListenAndServe(fmt.Sprintf(":%d", app.cfg.Listen), nil))
	}()
	return nil
}

func (app *App) Stop(s service.Service) error {
	return nil
}

func (app *App) limit() {
	for {
		usage, err := disk.Usage(app.cfg.CacheDir)
		if err == nil {
			if usage.UsedPercent > float64(app.cfg.CacheThreshold) {
				app.blocked = true
			} else {
				app.blocked = false
			}
		}
		time.Sleep(time.Second)
	}
}
