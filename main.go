package main

import (
	"fmt"
	rt "runtime"
	"server/docs"
	"server/internal"

	"github.com/jkstack/jkframe/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "metrics-agent",
	Long: "jkstack metrics agent",
	Run:  internal.Run,
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "register service",
	Run:   internal.Install,
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "unregister service",
	Run:   internal.Uninstall,
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run program",
	Run:   internal.Run,
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start service",
	Run:   internal.Start,
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop service",
	Run:   internal.Stop,
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart service",
	Run:   internal.Restart,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "show service status",
	Run:   internal.Status,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version info",
	Run:   showVersion,
}

var (
	version      string = "0.0.0"
	gitBranch    string = "<branch>"
	gitHash      string = "<hash>"
	gitReversion string = "0"
	buildTime    string = "0000-00-00 00:00:00"
)

func showVersion(*cobra.Command, []string) {
	fmt.Printf("version: %s\ncode version: %s.%s.%s\nbuild time: %s\ngo version: %s\n",
		version,
		gitBranch, gitHash, gitReversion,
		buildTime,
		rt.Version())
}

func main() {
	internal.Version = version
	docs.SwaggerInfo.Version = version

	installCmd.Flags().StringVar(&internal.ConfDir, "conf", "", "configure file dir")
	runCmd.Flags().StringVar(&internal.ConfDir, "conf", "", "configure file dir")
	rootCmd.AddCommand(installCmd, uninstallCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(startCmd, stopCmd, restartCmd, statusCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Flags().StringVar(&internal.ConfDir, "conf", "", "configure file dir")
	utils.Assert(rootCmd.Execute())
}
