package app

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"server/internal/agent"
	"server/internal/api"
	"server/internal/api/agents"
	"server/internal/api/exec"
	"server/internal/api/file"
	"server/internal/api/foo"
	"server/internal/api/info"
	"server/internal/api/layout"
	"server/internal/api/metrics"
	"server/internal/api/rpa"
	"server/internal/api/script"
	"server/internal/conf"
	"server/internal/utils"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"github.com/jkstack/jkframe/stat"
	runtime "github.com/jkstack/jkframe/utils"
	"github.com/kardianos/service"
	"github.com/shirou/gopsutil/v3/disk"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type apiHandler interface {
	Module() string
	Init(*conf.Configure, *stat.Mgr)
	HandleFuncs() map[api.Route]func(*gin.Context)
	OnConnect(*agent.Agent)
	OnClose(string)
	OnMessage(*agent.Agent, *anet.Msg)
}

// App application
type App struct {
	cfg     *conf.Configure
	version string
	stats   *stat.Mgr
	agents  *agent.Agents

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
		cfg:     cfg,
		version: version,
		agents:  agent.NewAgents(st),
		stats:   st,
		// runtime
		stAgentCount: st.NewCounter("agent_count"),
		connectLimit: rate.NewLimiter(rate.Limit(cfg.ConnectLimit), 1),
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
			app.context,
			gin.RecoveryWithWriter(io.Discard, handleRecovery),
			app.ratelimit,
			app.point,
		)

		var apis []apiHandler
		apis = append(apis, agents.New())
		apis = append(apis, foo.New())
		apis = append(apis, info.New(app.version, &app.blocked))
		apis = append(apis, metrics.New())
		apis = append(apis, exec.New())
		apis = append(apis, file.New())
		apis = append(apis, script.New())
		apis = append(apis, layout.New())

		for _, api := range apis {
			api.Init(app.cfg, app.stats)
			g := apiGroup.Group("/" + api.Module())
			for route, cb := range api.HandleFuncs() {
				logging.Info("route => %6s /api/%s%s", route.Method, api.Module(), route.URI)
				g.Handle(route.Method, route.URI, cb)
			}
		}

		logging.Info("route => %6s /metrics", http.MethodGet)
		g.GET("/metrics", func(g *gin.Context) {
			app.stats.ServeHTTP(g.Writer, g.Request)
		})
		logging.Info("route => %6s /ws/agent", http.MethodGet)
		g.GET("/ws/agent", func(g *gin.Context) {
			app.handleWS(g, apis)
		})
		logging.Info("route => %6s /docs/*any", http.MethodGet)
		g.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		go app.listenGRPC(app.cfg.GrpcListen)

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

// Stop stop service
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

func (app *App) listenGRPC(port uint16) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	runtime.Assert(err)
	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: 30 * time.Second,
		}),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	rpa.RegisterRpaServer(s, rpa.NewGRPC(app.agents))
	logging.Info("grpc listen on %d", port)
	runtime.Assert(s.Serve(lis))
}
