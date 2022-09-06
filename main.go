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
	"github.com/kardianos/service"
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

	cfg := conf.Load(*cf, filepath.Join(filepath.Dir(dir), "/../"))

	app := app.New(cfg, version)
	sv, err := service.New(app, appCfg)
	runtime.Assert(err)

	docs.SwaggerInfo.Version = version

	switch *act {
	case "install":
		fmt.Printf("platform: %s\n", sv.Platform())
		runtime.Assert(sv.Install())
	case "uninstall":
		runtime.Assert(sv.Uninstall())
	default:
		runtime.Assert(sv.Run())
	}
}
