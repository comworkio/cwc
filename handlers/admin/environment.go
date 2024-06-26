package admin

import (
	"cwc/admin"
	"cwc/config"
	"cwc/utils"
	"fmt"
	"os"
	"os/exec"

	"github.com/olekukonko/tablewriter"
)

func HandleAddEnvironment(name *string, path *string, roles *string, is_private *bool, description *string, subdomains *string, logo_url *string) {
	added_environment := &admin.Environment{
		Name:                *name,
		Path:                *path,
		Roles:               *roles,
		IsPrivate:           *is_private,
		Description:         *description,
		SubDomains:          *subdomains,
		LogUrl:              *logo_url,
		EnvironmentTemplate: "",
		DocTemplate:         "",
	}
	// prompt for environment template
	fmt.Print("Do you want to add environment template? [Y/N]: ")
	var add_env_template string
	fmt.Scanln(&add_env_template)

	if add_env_template == "Y" || add_env_template == "y" {
		editorCommand := utils.GetSystemEditor()

		// Create a temporary file with a specific name and path
		tempFileName := "temp-code-editor.txt"
		_, err := os.Create(tempFileName)
		utils.ExitIfErrorWithMsg("Error creating temporary file", err)

		defer os.Remove(tempFileName)

		// Prompt the user to write code in the editor
		fmt.Printf("Please write your code in the text editor that opens. Save and close the editor when done.\n")

		cmd := exec.Command(editorCommand, tempFileName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		utils.ExitIfErrorWithMsg("Error opening the text editor", err)

		// Read the code from the temporary file
		codeBytes, err := os.ReadFile(tempFileName)
		utils.ExitIfErrorWithMsg("Error reading code from the text editor", err)

		added_environment.EnvironmentTemplate = string(codeBytes)
	}

	// prompt for doc template
	fmt.Print("Do you want to add doc template? [Y/N]: ")
	var add_doc_template string
	fmt.Scanln(&add_doc_template)

	if add_doc_template == "Y" || add_doc_template == "y" {
		editorCommand := utils.GetSystemEditor()

		// Create a temporary file with a specific name and path
		tempFileName := "temp-code-editor.txt"
		_, err := os.Create(tempFileName)
		utils.ExitIfErrorWithMsg("Error creating temporary file", err)

		defer os.Remove(tempFileName)

		// Prompt the user to write code in the editor
		fmt.Printf("Please write your code in the text editor that opens. Save and close the editor when done.\n")

		cmd := exec.Command(editorCommand, tempFileName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		utils.ExitIfErrorWithMsg("Error opening the text editor", err)

		// Read the code from the temporary file
		codeBytes, err := os.ReadFile(tempFileName)
		utils.ExitIfErrorWithMsg("Error reading code from the text editor", err)

		added_environment.DocTemplate = string(codeBytes)
	}

	c, err := admin.NewClient()
	utils.ExitIfError(err)

	created_env, err := c.AdminAddEnvironment(*added_environment)
	utils.ExitIfError(err)

	if config.GetDefaultFormat() == "json" {
		utils.PrintJson(*created_env)
	} else {
		utils.PrintPretty(fmt.Sprintf("Environment %s successfully created", created_env.Name), *created_env)
	}
}

func HandleDeleteEnvironment(id *string) {
	c, err := admin.NewClient()
	utils.ExitIfError(err)

	err = c.AdminDeleteEnvironment(*id)
	utils.ExitIfError(err)

	fmt.Printf("Environment %v successfully deleted\n", *id)
}

func HandleGetEnvironments(environments *[]admin.Environment, pretty *bool) {
	if config.IsPrettyFormatExpected(pretty) {
		displayEnvironmentsAsTable(*environments)
	} else if config.GetDefaultFormat() == "json" {
		utils.PrintJson(environments)
	} else {
		utils.PrintMultiRow(admin.Environment{}, *environments)
	}
}

func HandleGetEnvironment(environment *admin.Environment, pretty *bool) {
	if config.IsPrettyFormatExpected(pretty) {
		utils.PrintPretty("Environment found", environment)
	} else if config.GetDefaultFormat() == "json" {
		utils.PrintJson(environment)
	} else {
		utils.PrintRow(environment)
	}
}

func displayEnvironmentsAsTable(environments []admin.Environment) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Path", "Description", "Subdomains", "Is Private"})

	if len(environments) == 0 {
		fmt.Println("No environments found")
	} else {
		for _, environment := range environments {
			table.Append([]string{
				fmt.Sprintf("%d", environment.Id),
				environment.Name,
				environment.Path,
				environment.Description,
				environment.SubDomains,
				fmt.Sprintf("%t", environment.IsPrivate),
			})
		}
		table.Render()
	}
}
