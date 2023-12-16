package function

import (
	"cwc/cmd/admin/faas/function/ls"

	"github.com/spf13/cobra"
)

var FunctionCmd = &cobra.Command{
	Use:   "function",
	Short: "Manage your functions in the cloud",
	Long: `This command lets you manage your functions in the cloud.
Several actions are associated with this command such listing your available functions`,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	FunctionCmd.DisableFlagsInUseLine = true
	FunctionCmd.AddCommand(ls.LsCmd)
}