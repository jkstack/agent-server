package app

import (
	"fmt"
	"io"
	"net"
	"server/internal/agent"
	"server/internal/api"
	"server/internal/api/agents"
	"server/internal/api/foo"
	"server/internal/conf"
	"server/internal/utils"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/stat"
	runtime "github.com/jkstack/jkframe/utils"
	"github.com/kardianos/service"
	"github.com/shirou/gopsutil/v3/disk"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

type apiHandler interface {
	Module() string
	Init(*conf.Configure, *stat.Mgr)
	HandleFuncs() map[api.Route]func(*gin.Context)
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

		gin.SetMode(gin.ReleaseMode)
		g := gin.New()

		apiGroup := g.Group("/api")
		apiGroup.Use(
			gin.RecoveryWithWriter(io.Discard, handleRecovery),
			app.ratelimit,
			app.point,
			app.bind,
		)

		var apis []apiHandler
		apis = append(apis, agents.New())
		apis = append(apis, foo.New())

		for _, api := range apis {
			api.Init(app.cfg, app.stats)
			g := apiGroup.Group("/" + api.Module())
			for route, cb := range api.HandleFuncs() {
				logging.Info("route => %s /api/%s%s", route.Method, api.Module(), route.Uri)
				g.Handle(route.Method, route.Uri, cb)
			}
		}

		logging.Info("route => GET /metrics")
		g.GET("/metrics", func(g *gin.Context) {
			app.stats.ServeHTTP(g.Writer, g.Request)
		})
		logging.Info("route => GET /ws/agent")
		g.GET("/ws/agent", func(g *gin.Context) {
			app.handleWS(g, apis)
		})
		logging.Info("route => GET /docs/*any")
		g.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		logging.Info("http listen on %d", app.cfg.Listen)
		addrs, _ := net.InterfaceAddrs()
		for _, addr := range addrs {
			a, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			address := a.IP.String()
			if address == "127.0.0.1" {
				continue
			}
			logging.Info("  - docs url maybe in http://%s:%d/docs/index.html",
				address, app.cfg.Listen)
		}
		runtime.Assert(g.Run(fmt.Sprintf(":%d", app.cfg.Listen)))
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
