package cmd

import (
	"kubazulo/pkg/token"

	"github.com/spf13/cobra"
)

func GetToken() *cobra.Command {
	o := Options()
	cmd := &cobra.Command{
		Use:   "get-token",
		Short: "Attempts to retrieve the token from AzureAD",
		Long:  `Attempts to retrieve the token from AzureAD`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			token.InvokeTokenProcess(cmd.Flags())
		},
	}

	o.AddFlags(cmd.Flags())
	for _, rf := range RequiredFlags() {
		cmd.MarkFlagRequired(rf)
	}

	cmd.Flags().SortFlags = false
	cmd.MarkFlagsRequiredTogether(DependendFlags()...)

	return cmd
}
