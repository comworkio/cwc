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
	Short: "List available environments",
	Long: `This command lets you list the available environment in the cloud that can be associeted to an instance
This command takes no arguments`,
	Run: func(cmd *cobra.Command, args []string) {
		handlers.HandleGetEnvironment()
	},
}

func init() {

}