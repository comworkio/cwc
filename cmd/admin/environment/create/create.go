package create

import (
	"cwc/handlers/admin"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	name        string
	path        string
	roles       string
	privacy     bool = true
	description string
	subdomains  string
	logo_url    string
)

// createCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an environment in the cloud",
	Long:  `This command lets you create an environment in the cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		admin.HandleAddEnvironment(&name, &path, &roles, &privacy, &description, &subdomains, &logo_url)
	},
}

func init() {
	CreateCmd.Flags().StringVarP(&name, "name", "n", "", "The environment name")
	CreateCmd.Flags().StringVarP(&path, "path", "p", "", "The environment path")
	CreateCmd.Flags().StringVarP(&roles, "roles", "r", "", "The environment roles")
	CreateCmd.Flags().BoolVarP(&privacy, "private", "P", false, "The environment privacy")
	CreateCmd.Flags().StringVarP(&description, "description", "d", "", "The environment description")
	CreateCmd.Flags().StringVarP(&subdomains, "subdomains", "s", "", "The environment subdomains")
	CreateCmd.Flags().StringVarP(&logo_url, "logo_url", "l", "", "The environment logo url")

	err := CreateCmd.MarkFlagRequired("name")
	if nil != err {
		fmt.Println(err)
	}

	err = CreateCmd.MarkFlagRequired("path")
	if nil != err {
		fmt.Println(err)
	}

	err = CreateCmd.MarkFlagRequired("roles")
	if nil != err {
		fmt.Println(err)
	}
}
