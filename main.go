package main

import (
	"kubazulo/pkg/cmd"
	"kubazulo/pkg/utils"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "kubazulo",
		Short: "Kubeconfig Authentication Helper for Kubernetes API-Server in cunjunction with kubectl",
		Long: `Kubeconfig Authentication Helper for Kubernetes API-Server in cunjunction with kubectl
Kubazulo is a client-go credential (exec) plugin implementing Azure authentication. 
It plugs in seemless into the process of communicating to the kubernetes API-Server.`,
		Version: utils.PrintAppVersion(),
	}
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(cmd.GetToken())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
