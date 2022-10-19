/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package get

import (
	"cwc/handlers"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var GetProviderCmd = &cobra.Command{
	Use:   "provider",
	Short: "Get the default provider",
	Long: `This command lets you retrieve the default provider`,
	Run: func(cmd *cobra.Command, args []string) {
			handlers.HandlerGetDefaultProvider()

	},
}

func init() {

}
