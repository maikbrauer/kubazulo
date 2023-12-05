package cmd

import (
	"github.com/spf13/cobra"
	"kubazulo/pkg/token"
)

func GetToken() *cobra.Command {
	o := Options()
	cmd := &cobra.Command{
		Use:   "get-token",
		Short: "Tries to fetch the Token from AzureAD",
		Long:  `Tries to fetch the Token from AzureAD`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			token.GetTokenProcess(cmd.Flags())
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
