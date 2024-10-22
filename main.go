package main

import (
	"kubazulo/pkg/cmd"
	"kubazulo/pkg/utils"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "kubazulo",
		Short: "Kubeconfig Authentication Helper for Kubernetes API-Server with kubectl",
		Long: `Kubazulo is a client-go credential (exec) plugin that implements Azure authentication. 
It seamlessly integrates into the process of communicating with the Kubernetes API-Server.`,
		Version: utils.PrintAppVersion() + "\ngo-runtime: " + runtime.Version(),
	}
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(cmd.GetToken())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
