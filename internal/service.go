package internal

import (
	"fmt"
	"os"
	"path/filepath"
	rt "runtime"
	"server/internal/app"
	"server/internal/conf"

	"github.com/jkstack/jkframe/utils"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

var ConfDir string
var Version string

type dummy struct{}

func (*dummy) Start(s service.Service) error {
	return nil
}

func (*dummy) Stop(s service.Service) error {
	return nil
}

func newService(app service.Interface) (service.Service, error) {
	var user string
	var depends []string
	if rt.GOOS != "windows" {
		user = "root"
		depends = append(depends, "After=network.target")
	}

	dir, err := filepath.Abs(ConfDir)
	if err != nil {
		return nil, err
	}

	opt := make(service.KeyValue)
	opt["LimitNOFILE"] = 65535

	return service.New(app, &service.Config{
		Name:         "agent-server",
		DisplayName:  "agent-server",
		Description:  "agent server",
		UserName:     user,
		Arguments:    []string{"--conf", dir},
		Dependencies: depends,
		Option:       opt,
	})
}

func dummyService() (service.Service, error) {
	return newService(&dummy{})
}

func newApp() (service.Service, error) {
	dir, err := os.Executable()
	utils.Assert(err)
	cfg := conf.Load(ConfDir, filepath.Join(filepath.Dir(dir), "/../"))
	return newService(app.New(cfg, Version))
}

func Install(*cobra.Command, []string) {
	if len(ConfDir) == 0 {
		fmt.Println("missing --conf argument")
		os.Exit(1)
	}

	svc, err := dummyService()
	if err != nil {
		fmt.Printf("can not create service: %v\n", err)
		return
	}
	err = svc.Install()
	if err != nil {
		fmt.Printf("can not register service: %v\n", err)
		return
	}
	fmt.Println("register service success")
}

func Uninstall(*cobra.Command, []string) {
	svc, err := dummyService()
	if err != nil {
		fmt.Printf("can not create service: %v\n", err)
		return
	}
	err = svc.Uninstall()
	if err != nil {
		fmt.Printf("can not unregister service: %v\n", err)
		return
	}
	fmt.Println("unregister service success")
}

func Start(*cobra.Command, []string) {
	svc, err := dummyService()
	if err != nil {
		fmt.Printf("can not create service: %v\n", err)
		return
	}
	err = svc.Start()
	if err != nil {
		fmt.Printf("can not start service: %v\n", err)
		return
	}
	fmt.Println("start service success")
}

func Stop(*cobra.Command, []string) {
	svc, err := dummyService()
	if err != nil {
		fmt.Printf("can not create service: %v\n", err)
		return
	}
	err = svc.Stop()
	if err != nil {
		fmt.Printf("can not stop service: %v\n", err)
		return
	}
	fmt.Println("stop service success")
}

func Restart(*cobra.Command, []string) {
	svc, err := dummyService()
	if err != nil {
		fmt.Printf("can not create service: %v\n", err)
		return
	}
	err = svc.Stop()
	if err != nil {
		fmt.Printf("can not restart service: %v\n", err)
		return
	}
	fmt.Println("restart service success")
}

func Status(*cobra.Command, []string) {
	svc, err := dummyService()
	if err != nil {
		fmt.Printf("can not create service: %v\n", err)
		return
	}
	status, err := svc.Status()
	if err != nil {
		fmt.Printf("can not get service status: %v\n", err)
		return
	}
	switch status {
	case service.StatusRunning:
		fmt.Println("service is running")
	case service.StatusStopped:
		fmt.Println("service is stopped")
	case service.StatusUnknown:
		fmt.Println("service status unknown")
	}
}

func Run(*cobra.Command, []string) {
	svc, err := newApp()
	if err != nil {
		fmt.Printf("can not create service: %v\n", err)
		return
	}
	utils.Assert(svc.Run())
}
