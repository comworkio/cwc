package set

import (
	"cwc/handlers/user"
	"cwc/utils"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var SetRegionCmd = &cobra.Command{
	Use:   "region",
	Short: "Set the default endpoint",
	Long:  `This command lets you update the default endpoint`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.ExitIfNeeded("You have to provide a value", len(args) == 0)

		value := args[0]
		user.HandlerSetDefaultRegion(value)
	},
}

func init() {
}
