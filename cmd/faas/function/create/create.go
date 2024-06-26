package create

import (
	"cwc/client"
	"cwc/handlers/user"
	"cwc/utils"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	function    client.Function
	interactive bool = false
	pretty      bool = false
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a function in the cloud",
	Long:  "This command lets you create a function in the cloud.",
	Run: func(cmd *cobra.Command, args []string) {
		created_function, err := user.PrepareAddFunction(&function, &interactive)
		utils.ExitIfError(err)
		user.HandleAddFunction(created_function, &pretty)
	},
}

func init() {
	CreateCmd.Flags().StringVarP(&function.Content.Name, "name", "n", "", "Name of the function")
	CreateCmd.Flags().StringVarP(&function.Content.Language, "language", "l", "", "Language of the function")
	CreateCmd.Flags().StringSliceVarP(&function.Content.Args, "args", "g", []string{}, "Arguments of the function")
	CreateCmd.Flags().BoolVar(&function.Is_public, "is_public", false, "Is the function public? (optional)")
	CreateCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode (optional)")
	CreateCmd.Flags().BoolVarP(&pretty, "pretty", "p", false, "Pretty print the output (optional)")
	CreateCmd.Flags().StringVarP(&function.Content.Regexp, "regexp", "r", "", "Arguments matching regexp (optional)")
	CreateCmd.Flags().StringVarP(&function.Content.Code, "code", "c", "", "Code of the function (optional)")

	err := CreateCmd.MarkFlagRequired("name")
	if nil != err {
		fmt.Println(err)
	}

	err = CreateCmd.MarkFlagRequired("language")
	if nil != err {
		fmt.Println(err)
	}
}
