package ls

import (
	adminClient "cwc/admin"
	"cwc/handlers/admin"
	"cwc/utils"

	"github.com/spf13/cobra"
)

var (
	environmentId string
	pretty        bool = false
)

// lsCmd represents the ls command
var LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List available environments",
	Long: `This command lets you list the available environment in the cloud that can be associeted to an instance
This command takes no arguments`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := adminClient.NewClient()
		utils.ExitIfError(err)
		if utils.IsBlank(environmentId) {
			environments, err := c.GetAllEnvironments()
			utils.ExitIfError(err)
			admin.HandleGetEnvironments(environments, &pretty)
		} else {
			environment, err := c.GetEnvironment(environmentId)
			utils.ExitIfError(err)
			admin.HandleGetEnvironment(environment, &pretty)
		}
	},
}

func init() {
	LsCmd.Flags().StringVarP(&environmentId, "environment", "e", "", "The environment id")
	LsCmd.Flags().BoolVarP(&pretty, "pretty", "p", false, "Pretty print the output (optional)")
}
