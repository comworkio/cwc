package admin

import (
	"cwc/admin"
	"cwc/config"
	"cwc/utils"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func HandleGetMonitors(monitors *[]admin.Monitor, pretty *bool) {
	if config.IsPrettyFormatExpected(pretty) {
		displayMonitorsAsTable(*monitors)
	} else if config.GetDefaultFormat() == "json" {
		utils.PrintJson(monitors)
	} else {
		var monitorsDisplay []admin.Monitor
		for i, monitor := range *monitors {
			monitorsDisplay = append(monitorsDisplay, admin.Monitor{
				Id:            monitor.Id,
				Name:          monitor.Name,
				Family:        monitor.Family,
				Url:           monitor.Url,
				Method:        monitor.Method,
				Timeout:       monitor.Timeout,
				Updated_at:    monitor.Updated_at,
				Status:        monitor.Status,
				Response_time: monitor.Response_time,
			})
			monitorsDisplay[i].Id = monitor.Id
		}
		utils.PrintMultiRow(admin.Monitor{}, monitorsDisplay)
	}
}

func HandleGetMonitor(monitor *admin.Monitor, pretty *bool) {
	var monitorDisplay admin.Monitor
	monitorDisplay.Id = monitor.Id
	monitorDisplay.Type = monitor.Type
	monitorDisplay.Name = monitor.Name
	monitorDisplay.Family = monitor.Family
	monitorDisplay.Url = monitor.Url
	monitorDisplay.Method = monitor.Method
	monitorDisplay.Expected_http_code = monitor.Expected_http_code
	if monitor.Method == "POST" || monitor.Method == "PUT" {
		monitorDisplay.Body = monitor.Body
	}
	monitorDisplay.Timeout = monitor.Timeout
	monitorDisplay.Username = monitor.Username
	monitorDisplay.Password = monitor.Password
	monitorDisplay.Headers = monitor.Headers
	monitorDisplay.Status = monitor.Status
	monitorDisplay.Response_time = monitor.Response_time
	monitorDisplay.Updated_at = monitor.Updated_at

	if config.IsPrettyFormatExpected(pretty) {
		utils.PrintPretty("Found monitor", monitorDisplay)
	} else if config.GetDefaultFormat() == "json" {
		utils.PrintJson(monitor)
	} else {
		utils.PrintRow(monitorDisplay)
	}
}

func PrepareAddMonitor(monitor *admin.Monitor) (admin.Monitor, error) {
	c, err := admin.NewClient()
	utils.ExitIfError(err)

	created_monitor, err := c.AddMonitor(*monitor)
	utils.ExitIfError(err)
	return *created_monitor, err
}

func HandleAddMonitor(createdMonitor *admin.Monitor, pretty *bool) {
	if config.IsPrettyFormatExpected(pretty) {
		utils.PrintPretty("Monitor successfully created", *createdMonitor)
	} else if config.GetDefaultFormat() == "json" {
		utils.PrintJson(createdMonitor)
	} else {
		utils.PrintRow(*createdMonitor)
	}
}

func HandleUpdateMonitor(monitorId *string, updatedMonitor *admin.Monitor) {
	c, err := admin.NewClient()
	utils.ExitIfError(err)

	monitor, err := c.GetMonitorById(*monitorId)
	utils.ExitIfError(err)

	if utils.IsNotBlank(updatedMonitor.Type) {
		monitor.Type = updatedMonitor.Type
	}

	if utils.IsNotBlank(updatedMonitor.Name) {
		monitor.Name = updatedMonitor.Name
	}

	if utils.IsNotBlank(updatedMonitor.Family) {
		monitor.Family = updatedMonitor.Family
	}

	if utils.IsNotBlank(updatedMonitor.Url) {
		monitor.Url = updatedMonitor.Url
	}

	if utils.IsNotBlank(updatedMonitor.Method) {
		monitor.Method = updatedMonitor.Method
	}

	if utils.IsNotBlank(updatedMonitor.Body) {
		monitor.Body = updatedMonitor.Body
	}

	if utils.IsNotBlank(updatedMonitor.Expected_contain) {
		monitor.Expected_contain = updatedMonitor.Expected_contain
	}

	if utils.IsNotBlank(updatedMonitor.Username) {
		monitor.Username = updatedMonitor.Username
	}

	if utils.IsNotBlank(updatedMonitor.Password) {
		monitor.Password = updatedMonitor.Password
	}

	if len(updatedMonitor.Headers) > 0 {
		monitor.Headers = updatedMonitor.Headers
	}

	if updatedMonitor.User_id != 0 {
		monitor.User_id = updatedMonitor.User_id
	}

	_, updateError := c.UpdateMonitorById(*monitorId, *monitor)
	utils.ExitIfError(updateError)

	fmt.Println("Monitor successfully updated")
}

func HandleDeleteMonitor(monitorId *string) {
	c, err := admin.NewClient()
	utils.ExitIfError(err)

	err = c.DeleteMonitorById(*monitorId)
	utils.ExitIfError(err)

	fmt.Println("Monitor successfully deleted")
}

func displayMonitorsAsTable(monitors []admin.Monitor) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name", "Family", "Method", "Url", "Updated_at", "Status", "Response_time"})

	if len(monitors) == 0 {
		table.Append([]string{"No monitors available", "404", "404", "404", "404", "404", "404", "404"})
	} else {
		for _, monitor := range monitors {
			table.Append([]string{
				monitor.Id,
				monitor.Name,
				monitor.Family,
				monitor.Method,
				monitor.Url,
				monitor.Updated_at,
				monitor.Status,
				monitor.Response_time,
			})
		}
		table.Render()
	}
}