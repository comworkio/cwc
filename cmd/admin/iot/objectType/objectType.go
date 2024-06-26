package objectType

import (
	"cwc/cmd/admin/iot/objectType/create"
	"cwc/cmd/admin/iot/objectType/delete"
	"cwc/cmd/admin/iot/objectType/ls"
	"cwc/cmd/admin/iot/objectType/update"

	"github.com/spf13/cobra"
)

var ObjectTypeCmd = &cobra.Command{
	Use:   "objectType",
	Short: "Manage your object types in the cloud",
	Long: `This command lets you Manage your object types in the cloud.
Several actions are associated with this command such as creating an object type`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	ObjectTypeCmd.DisableFlagsInUseLine = true
	ObjectTypeCmd.AddCommand(create.CreateCmd)
	ObjectTypeCmd.AddCommand(ls.LsCmd)
	ObjectTypeCmd.AddCommand(delete.DeleteCmd)
	ObjectTypeCmd.AddCommand(update.UpdateCmd)
}
