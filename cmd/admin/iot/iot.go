package iot

import (
	"cwc/cmd/admin/iot/data"
	"cwc/cmd/admin/iot/device"
	"cwc/cmd/admin/iot/objectType"

	"github.com/spf13/cobra"
)

var IotCmd = &cobra.Command{
	Use:   "iot",
	Short: "Manage your internet of things in the cloud",
	Long: `This command lets you manage your internet of things in the cloud.
Several actions are associated with this command such as managing your object types`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	IotCmd.DisableFlagsInUseLine = true
	IotCmd.AddCommand(objectType.ObjectTypeCmd)
	IotCmd.AddCommand(device.DeviceCmd)
	IotCmd.AddCommand(data.DataCmd)
}
