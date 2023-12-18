package ls

import (
	"cwc/handlers/user"

	"github.com/spf13/cobra"
)

var (
	pretty bool
)

// lsCmd represents the ls command
var LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List availble instances types",
	Long:  `List availble instances types`,
	Run: func(cmd *cobra.Command, args []string) {
		user.HandleListInstancesTypes(&pretty)
	},
}

func init() {
	LsCmd.Flags().BoolVarP(&pretty, "pretty", "p", false, "Pretty print the output (optional)")
}
