package user

import (
	"cwc/client"
	"cwc/config"
	"cwc/utils"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func HandleGetEnvironments(environments *[]client.Environment, pretty *bool) {
	if config.IsPrettyFormatExpected(pretty) {
		displayEnvironmentsAsTable(*environments)
	} else if config.GetDefaultFormat() == "json" {
		utils.PrintJson(environments)
	} else {
		utils.PrintMultiRow(client.Environment{}, *environments)
	}
}

func HandleGetEnvironment(environment *client.Environment, pretty *bool) {
	if config.IsPrettyFormatExpected(pretty) {
		utils.PrintPretty("Environment found", *environment)
	} else if config.GetDefaultFormat() == "json" {
		utils.PrintJson(environment)
	} else {
		utils.PrintRow(environment)
	}
}

func displayEnvironmentsAsTable(environments []client.Environment) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Path", "Description"})

	if len(environments) == 0 {
		fmt.Println("No environments found")
		return
	} else {
		for _, environment := range environments {
			table.Append([]string{
				fmt.Sprintf("%d", environment.Id),
				environment.Name,
				environment.Path,
				environment.Description,
			})
		}
		table.Render()
	}
}
