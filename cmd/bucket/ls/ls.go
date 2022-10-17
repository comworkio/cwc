/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package ls

import (
	"cwc/handlers"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List available buckets",
	Long: `This command lets you list your available buckets in the cloud
This command takes no arguments`,
	Run: func(cmd *cobra.Command, args []string) {
		handlers.HandleGetBucket()
	},
}

func init() {

}