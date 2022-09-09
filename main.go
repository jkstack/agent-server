package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	rt "runtime"
	"server/docs"
	"server/internal/app"
	"server/internal/conf"

	runtime "github.com/jkstack/jkframe/utils"
	"github.com/lwch/service"
)

var (
	version      string = "0.0.0"
	gitBranch    string = "<branch>"
	gitHash      string = "<hash>"
	gitReversion string = "0"
	buildTime    string = "0000-00-00 00:00:00"
)

func showVersion() {
	fmt.Printf("version: %s\ncode version: %s.%s.%s\nbuild time: %s\ngo version: %s\n",
		version,
		gitBranch, gitHash, gitReversion,
		buildTime,
		rt.Version())
}

func main() {
	cf := flag.String("conf", "", "config file dir")
	ver := flag.Bool("version", false, "show version info")
	act := flag.String("action", "", "install or uninstall")
	flag.Parse()

	if *ver {
		showVersion()
		return
	}

	if len(*cf) == 0 {
		fmt.Println("缺少-conf参数")
		os.Exit(1)
	}

	var user string
	var depends []string
	if rt.GOOS != "windows" {
		user = "root"
		depends = append(depends, "After=network.target")
	}

	dir, err := filepath.Abs(*cf)
	runtime.Assert(err)

	opt := make(service.KeyValue)
	opt["LimitNOFILE"] = 65535

	appCfg := &service.Config{
		Name:         "agent-server",
		DisplayName:  "agent-server",
		Description:  "agent server",
		UserName:     user,
		Arguments:    []string{"-conf", dir},
		Dependencies: depends,
		Option:       opt,
	}

	dir, err = os.Executable()
	runtime.Assert(err)

	docs.SwaggerInfo.Version = version

	var sv service.Service
	if *act == "install" || *act == "uninstall" {
		sv, err = service.New(&dummy{}, appCfg)
		runtime.Assert(err)
	} else {
		cfg := conf.Load(*cf, filepath.Join(filepath.Dir(dir), "/../"))

		app := app.New(cfg, version)
		sv, err = service.New(app, appCfg)
		runtime.Assert(err)
	}

	switch *act {
	case "install":
		fmt.Printf("service name: %s\n", "agent-server")
		fmt.Printf("platform: %s\n", sv.Platform())
		err := sv.Install()
		if err != nil {
			fmt.Printf("Install failed: %v\n", err)
		}
	case "uninstall":
		sv.Stop()
		err := sv.Uninstall()
		if err != nil {
			fmt.Printf("Uninstall failed: %v\n", err)
		}
	default:
		runtime.Assert(sv.Run())
	}
}

type dummy struct{}

func (*dummy) Start(s service.Service) error {
	return nil
}

func (*dummy) Stop(s service.Service) error {
	return nil
}
