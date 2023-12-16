package user

import (
	"bufio"
	"cwc/client"
	"cwc/utils"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func HandleGetLanguages(pretty *bool) {
	languages, err := client.GetLanguages()
	if nil != err {
		fmt.Printf("failed: %s\n", err)
		os.Exit(1)
	}

	if client.GetDefaultFormat() == "json" {
		utils.PrintJson(languages)
	} else if *pretty {
		utils.PrintPrettyArray("Available languages", languages.Languages)
	} else {
		utils.PrintArray(languages.Languages)
	}
}

func HandleGetTriggerKinds(pretty *bool) {
	triggerKinds, err := client.GetTriggerKinds()
	if nil != err {
		fmt.Printf("failed: %s\n", err)
		os.Exit(1)
	}

	if client.GetDefaultFormat() == "json" {
		utils.PrintJson(triggerKinds)
	} else if *pretty {
		utils.PrintPrettyArray("Available trigger kinds", triggerKinds.TriggerKinds)
	} else {
		utils.PrintArray(triggerKinds.TriggerKinds)
	}
}

func HandleGetFunctions(pretty *bool) {
	c, _ := client.NewClient()
	functions, err := c.GetAllFunctions()
	if nil != err {
		fmt.Printf("failed: %s\n", err)
		os.Exit(1)
	}

	if client.GetDefaultFormat() == "json" {
		utils.PrintJson(functions)
	} else if *pretty {
		displayFunctionsAsTable(*functions)
	} else {
		var functionsDisplay []client.FunctionDisplay
		for i, function := range *functions {
			functionsDisplay = append(functionsDisplay, client.FunctionDisplay{
				Id:         function.Id,
				Owner_id:   function.Owner_id,
				Is_public:  function.Is_public,
				Name:       function.Content.Name,
				Language:   function.Content.Language,
				Created_at: function.Created_at,
				Updated_at: function.Updated_at,
			})
			functionsDisplay[i].Id = function.Id
		}

		utils.PrintMultiRow(client.FunctionDisplay{}, functionsDisplay)
	}
}

func HandleGetFunction(id *string, pretty *bool) {
	c, _ := client.NewClient()
	function, err := c.GetFunctionById(*id)
	if nil != err {
		fmt.Printf("failed: %s\n", err)
		os.Exit(1)
	}

	var functionDisplay client.FunctionDisplay
	functionDisplay.Id = function.Id
	functionDisplay.Owner_id = function.Owner_id
	functionDisplay.Is_public = function.Is_public
	functionDisplay.Name = function.Content.Name
	functionDisplay.Language = function.Content.Language
	functionDisplay.Created_at = function.Created_at
	functionDisplay.Updated_at = function.Updated_at

	if client.GetDefaultFormat() == "json" {
		utils.PrintJson(function)
	} else if *pretty {
		utils.PrintPretty("Found function", functionDisplay)
	} else {
		utils.PrintRow(functionDisplay)
	}
}

func displayFunctionsAsTable(functions []client.Function) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Language", "Public", "Created At", "Updated At"})

	if len(functions) == 0 {
		// If there are no functions available, display a message in a single cell.
		table.Append([]string{"No functions available", "404", "404", "404", "404", "404"})
	} else {
		for _, function := range functions {
			table.Append([]string{
				function.Id,
				function.Content.Name,
				function.Content.Language,
				fmt.Sprintf("%t", function.Is_public),
				function.Created_at,
				function.Updated_at,
			})
		}
	}

	table.Render() // Render the table
}

func HandleDeleteFunction(id *string) {
	c, _ := client.NewClient()
	delete_err := c.DeleteFunctionById(*id)
	if nil != delete_err {
		fmt.Printf("failed: %s\n", delete_err)
		os.Exit(1)
	}

	fmt.Printf("Function successfully deleted\n")
}

func HandleAddFunction(function *client.Function, interactive *bool, pretty *bool) {
	language_response, err := client.GetLanguages()
	if nil != err {
		fmt.Printf("failed: %s\n", err)
		os.Exit(1)
	}

	isLanguageAllowed := false
	for _, allowedLang := range language_response.Languages {
		if function.Content.Language == allowedLang {
			isLanguageAllowed = true
			break
		}
	}

	if !isLanguageAllowed {
		fmt.Printf("Invalid language. Allowed languages are: %s\n", strings.Join(language_response.Languages, ", "))
		os.Exit(1)
	}

	if *interactive {
		// Prompt for Regexp
		fmt.Print("Enter Regexp (or press Enter for empty): ")
		fmt.Scanln(&function.Content.Regexp)

		// Prompt for Callback URL
		fmt.Print("Enter Callback URL (or press Enter for empty): ")
		fmt.Scanln(&function.Content.Callback_url)

		// Prompt for Callback Authorization Header
		fmt.Print("Enter Callback Authorization Header (or press Enter for empty): ")
		fmt.Scanln(&function.Content.Callback_authorization_header)

		// Prompt for Args array
		fmt.Println("Enter Args (one per line, press Enter for each entry; leave an empty line to finish):")
		for {
			var arg string
			_, err := fmt.Scanln(&arg)
			if nil != err {
				break // Exit the loop if an error occurs (e.g., empty line)
			}
			function.Content.Args = append(function.Content.Args, arg)
		}
		if len(function.Content.Args) == 0 {
			function.Content.Args = []string{}
		}

		c, _ := client.NewClient()

		// assign the code template after choosing the language
		code_template, err := c.GetFunctionCodeTemplate(function.Content.Args, function.Content.Language)
		if nil != err {
			fmt.Printf("failed: %s\n", err)
			os.Exit(1)
		}

		fmt.Print("Do you want to add code? [Y/N]: ")
		var addCode string
		fmt.Scanln(&addCode)

		if addCode == "y" || addCode == "Y" {
			var editorCommand string
			editorCommand = os.Getenv("EDITOR")
			if editorCommand == "" {
				editorCommand = "vi" // Use 'vi' as the default editor
			}

			// Create a temporary file with a specific name and path
			tempFileName := "temp-code-editor.txt"
			tempFile, err := os.Create(tempFileName)
			if nil != err {
				fmt.Printf("Error creating temporary file: %s\n", err)
				os.Exit(1)
			}
			defer tempFile.Close()
			defer os.Remove(tempFileName)

			// Write the code_template to the temporary file
			_, err = tempFile.WriteString(*code_template)
			if nil != err {
				fmt.Printf("Error writing code_template to the temporary file: %s\n", err)
				os.Exit(1)
			}

			// Prompt the user to write code in the editor
			fmt.Printf("Please write your code in the text editor that opens. Save and close the editor when done.\n")

			cmd := exec.Command(editorCommand, tempFileName)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if nil != err {
				fmt.Printf("Error opening the text editor: %s\n", err)
				os.Exit(1)
			}

			// Read the code from the temporary file
			codeBytes, err := ioutil.ReadFile(tempFileName)
			if nil != err {
				fmt.Printf("Error reading code from the text editor: %s\n", err)
				os.Exit(1)
			}

			function.Content.Code = string(codeBytes)
		}

		fmt.Printf("code: %s\n", function.Content.Code)
	}

	c, _ := client.NewClient()
	created_function, err := c.AddFunction(*function)
	if nil != err {
		fmt.Printf("failed: %s\n", err)
		os.Exit(1)
	}

	if client.GetDefaultFormat() == "json" {
		utils.PrintJson(created_function)
	} else if *pretty {
		utils.PrintPretty("Function successfully created", *created_function)
	} else {
		utils.PrintRow(*created_function)
	}
}

func HandleUpdateFunction(id *string, updated_function *client.Function, interactive *bool) {
	language_response, err := client.GetLanguages()
	if nil != err {
		fmt.Printf("failed: %s\n", err)
		os.Exit(1)
	}

	c, _ := client.NewClient()
	function, err := c.GetFunctionById(*id)
	if nil != err {
		fmt.Printf("failed: %s\n", err)
		os.Exit(1)
	}

	isLanguageAllowed := false

	if *interactive {
		// prompt to update language
		fmt.Printf("Current language: %s\n", function.Content.Language)
		fmt.Printf("Available languages are: %s\n", strings.Join(language_response.Languages, ", "))
		fmt.Printf("Enter new language (or press Enter to keep the current one): ")
		fmt.Scanln(&function.Content.Language)

		for _, allowedLang := range language_response.Languages {
			if function.Content.Language == allowedLang {
				isLanguageAllowed = true
				break
			}
		}

		if !isLanguageAllowed {
			fmt.Printf("Invalid language. Allowed languages are: %s\n", strings.Join(language_response.Languages, ", "))
			os.Exit(1)
		}

		// Prompt for Regexp
		fmt.Printf("Current regexp: %s\n", function.Content.Regexp)
		fmt.Print("Enter new regexp (or press Enter to keep the current one): ")
		fmt.Scanln(&function.Content.Regexp)

		// Prompt for Callback URL
		fmt.Printf("Current callback URL: %s\n", function.Content.Callback_url)
		fmt.Print("Enter new callback URL (or press Enter to keep the current one): ")
		fmt.Scanln(&function.Content.Callback_url)

		// Prompt for Callback Authorization Header
		fmt.Printf("Current callback authorization header: %s\n", function.Content.Callback_authorization_header)
		fmt.Print("Enter new callback authorization header (or press Enter to keep the current one): ")
		fmt.Scanln(&function.Content.Callback_authorization_header)

		// Prompt for function name
		fmt.Printf("Current name: %s\n", function.Content.Name)
		fmt.Print("Enter new name (or press Enter to keep the current one): ")
		fmt.Scanln(&function.Content.Name)

		// Ask if the function should be public
		if function.Is_public {
			fmt.Print("Current function status: Public\n")
		} else {
			fmt.Print("Current function status: Private\n")
		}
		fmt.Print("Do you want to make the change the function status? [Y/N]: ")
		var answer string
		fmt.Scanln(&answer)
		if answer == "y" || answer == "Y" {
			function.Is_public = !function.Is_public
		}

		// Prompt for new Args array
		utils.PrintPrettyArray("Current args", function.Content.Args)
		fmt.Println("Enter new Args (one per line, press Enter for each entry; leave an empty line to finish):")
		for {
			var arg string
			_, err := fmt.Scanln(&arg)
			if nil != err {
				break // Exit the loop if an error occurs (e.g., empty line)
			}
			function.Content.Args = append(function.Content.Args, arg)
		}

		// Prompt for new code
		fmt.Print("Do you want to update code? [Y/N]: ")
		var updateCode string
		fmt.Scanln(&updateCode)

		if updateCode == "y" || updateCode == "Y" {
			var editorCommand string
			editorCommand = os.Getenv("EDITOR")
			if editorCommand == "" {
				editorCommand = "vi" // Use 'vi' as the default editor
			}

			// Create a temporary file with a specific name and path
			tempFileName := "temp-code-editor-update.txt"
			tempFile, err := os.Create(tempFileName)
			if nil != err {
				fmt.Printf("Error creating temporary file: %s\n", err)
				os.Exit(1)
			}

			defer os.Remove(tempFileName)

			// Write the current code to the temporary file
			_, err = tempFile.WriteString(function.Content.Code)
			if nil != err {
				fmt.Printf("Error writing current code to the temporary file: %s\n", err)
				os.Exit(1)
			}

			// Prompt the user to edit the code in the editor
			fmt.Printf("Please update your code in the text editor that opens. Save and close the editor when done.\n")

			cmd := exec.Command(editorCommand, tempFileName)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if nil != err {
				fmt.Printf("Error opening the text editor: %s\n", err)
				os.Exit(1)
			}

			// Read the updated code from the temporary file
			updatedCodeBytes, err := ioutil.ReadFile(tempFileName)
			if nil != err {
				fmt.Printf("Error reading updated code from the text editor: %s\n", err)
				os.Exit(1)
			}

			// Update the function's code with the edited code
			function.Content.Code = string(updatedCodeBytes)
		}
	} else {
		// If interactive mode is not enabled, update only the fields that are not empty
		if updated_function.Content.Language != "" {
			function.Content.Language = updated_function.Content.Language
		}

		if updated_function.Content.Regexp != "" {
			function.Content.Regexp = updated_function.Content.Regexp
		}

		if updated_function.Content.Callback_url != "" {
			function.Content.Callback_url = updated_function.Content.Callback_url
		}

		if updated_function.Content.Callback_authorization_header != "" {
			function.Content.Callback_authorization_header = updated_function.Content.Callback_authorization_header
		}

		if updated_function.Content.Name != "" {
			function.Content.Name = updated_function.Content.Name
		}

		if updated_function.Is_public {
			function.Is_public = updated_function.Is_public
		}

		if len(updated_function.Content.Args) > 0 {
			function.Content.Args = updated_function.Content.Args
		}

		if updated_function.Content.Code != "" {
			function.Content.Code = updated_function.Content.Code
		}
	}

	_, err = c.UpdateFunction(*function)
	if nil != err {
		fmt.Printf("failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Function successfully updated\n")
}

func HandleGetInvocations(pretty *bool) {
	c, _ := client.NewClient()

	invocations, err := c.GetAllInvocations()
	if nil != err {
		fmt.Printf("failed: %s\n", err)
		os.Exit(1)
	}

	if client.GetDefaultFormat() == "json" {
		utils.PrintJson(invocations)
	} else if *pretty {
		displayInvocationsAsTable(*invocations)
	} else {
		var invocationsDisplay []client.InvocationDisplay
		for i, invocation := range *invocations {
			invocationsDisplay = append(invocationsDisplay, client.InvocationDisplay{
				Id:          invocation.Id,
				Invoker_id:  invocation.Invoker_id,
				Function_id: invocation.Content.Function_id,
				State:       invocation.Content.State,
				Created_at:  invocation.Created_at,
				Updated_at:  invocation.Updated_at,
			})
			invocationsDisplay[i].Id = invocation.Id
		}

		utils.PrintMultiRow(client.InvocationDisplay{}, invocationsDisplay)
	}
}

func HandleGetInvocation(id *string, pretty *bool) {
	c, _ := client.NewClient()
	invocation, err := c.GetInvocationById(*id)
	if nil != err {
		fmt.Printf("failed: %s\n", err)
		os.Exit(1)
	}

	var invocationDisplay client.InvocationDisplay
	invocationDisplay.Id = invocation.Id
	invocationDisplay.Invoker_id = invocation.Invoker_id
	invocationDisplay.Function_id = invocation.Content.Function_id
	invocationDisplay.State = invocation.Content.State
	invocationDisplay.Created_at = invocation.Created_at
	invocationDisplay.Updated_at = invocation.Updated_at

	if client.GetDefaultFormat() == "json" {
		utils.PrintJson(invocation)
	} else if *pretty {
		utils.PrintPretty("Found invocation", invocationDisplay)
	} else {
		utils.PrintRow(invocationDisplay)
	}
}

func displayInvocationsAsTable(invocations []client.Invocation) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Function ID", "State", "Created At", "Updated At"})
	if len(invocations) == 0 {
		// If there are no invocations available, display a message in a single cell.
		table.Append([]string{"No invocations available", "404", "404", "404", "404", "404"})
	} else {
		for _, invocation := range invocations {
			table.Append([]string{
				invocation.Id,
				invocation.Content.Function_id,
				invocation.Content.State,
				invocation.Created_at,
				invocation.Updated_at,
			})
		}
	}
	table.Render() // Render the table
}

func HandleAddInvocation(content *client.InvocationAddContent, argument_values *[]string, interactive *bool, pretty *bool, synchronous *bool) {
	c, _ := client.NewClient()
	var id = &content.Function_id
	function, _ := c.GetFunctionById(*id)
	args := function.Content.Args

	if *interactive {
		// Prompt values of the existing arguments
		if len(args) > 0 {
			fmt.Println("Enter Args (one per line, press Enter for each entry; leave an empty line to finish):")
			for _, arg := range args {
				fmt.Printf("  ➤ %s: ", arg)
				var value string
				fmt.Scanln(&value)
				content.Args = append(content.Args, client.Argument{Key: arg, Value: value})
			}
		}
	} else {
		if len(*argument_values) != len(args) {
			fmt.Printf("Invalid number of arguments. Expected %d arguments, got %d\n", len(args), len(*argument_values))
			os.Exit(1)
		}

		if len(args) > 0 {
			for i, arg := range args {
				content.Args = append(content.Args, client.Argument{Key: arg, Value: (*argument_values)[i]})
			}
		}
	}

	created_invocation, err := c.AddInvocation(*content, *synchronous)
	if nil != err {
		fmt.Printf("failed: %s\n", err)
		os.Exit(1)
	}

	if client.GetDefaultFormat() == "json" {
		utils.PrintJson(*created_invocation)
	} else if *pretty {
		utils.PrintPretty("Invocation successfully created", *created_invocation)
	} else {
		utils.PrintRow(*created_invocation)
	}
}

func HandleDeleteInvocation(id *string) {
	c, _ := client.NewClient()
	delete_err := c.DeleteInvocationById(*id)
	if nil != delete_err {
		fmt.Printf("failed: %s\n", delete_err)
		os.Exit(1)
	}

	fmt.Printf("Invocation successfully deleted\n")
}

func HandleTruncateInvocations() {
	c, _ := client.NewClient()
	truncate_err := c.TruncateInvocations()
	if nil != truncate_err {
		fmt.Printf("failed: %s\n", truncate_err)
		os.Exit(1)
	}

	fmt.Printf("Invocations successfully truncated\n")
}

func HandleGetTriggers(pretty *bool) {
	c, _ := client.NewClient()
	triggers, err := c.GetAllTriggers()
	if nil != err {
		fmt.Printf("failed: %s\n", err)
		os.Exit(1)
	}

	if client.GetDefaultFormat() == "json" {
		utils.PrintJson(triggers)
	} else if *pretty {
		displayTriggersAsTable(*triggers)
	} else {
		var triggersDisplay []client.TriggerDisplay
		for i, trigger := range *triggers {
			triggersDisplay = append(triggersDisplay, client.TriggerDisplay{
				Id:          trigger.Id,
				Function_id: trigger.Content.Function_id,
				Kind:        trigger.Kind,
				Name:        trigger.Content.Name,
				Cron_expr:   trigger.Content.Cron_expr,
				Created_at:  trigger.Created_at,
				Updated_at:  trigger.Updated_at,
			})
			triggersDisplay[i].Id = trigger.Id
		}
		utils.PrintMultiRow(client.TriggerDisplay{}, triggersDisplay)
	}
}

func HandleGetTrigger(id *string, pretty *bool) {
	c, _ := client.NewClient()
	trigger, err := c.GetTriggerById(*id)
	if nil != err {
		fmt.Printf("failed: %s\n", err)
		os.Exit(1)
	}

	if client.GetDefaultFormat() == "json" {
		utils.PrintJson(trigger)
	} else {
		var triggerDisplay client.TriggerDisplay
		triggerDisplay.Id = trigger.Id
		triggerDisplay.Function_id = trigger.Content.Function_id
		triggerDisplay.Kind = trigger.Kind
		triggerDisplay.Name = trigger.Content.Name
		triggerDisplay.Cron_expr = trigger.Content.Cron_expr
		triggerDisplay.Created_at = trigger.Created_at
		triggerDisplay.Updated_at = trigger.Updated_at

		if client.GetDefaultFormat() == "json" {
			utils.PrintJson(trigger)
		} else if *pretty {
			utils.PrintPretty("Found trigger", triggerDisplay)
		} else {
			utils.PrintRow(triggerDisplay)
		}
	}
}

func displayTriggersAsTable(triggers []client.Trigger) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Function ID", "Kind", "Name", "Cron Expr", "Created At", "Updated At"})

	if len(triggers) == 0 {
		// If there are no functions available, display a message in a single cell.
		table.Append([]string{"No triggers available", "404", "404", "404", "404"})
	} else {
		for _, trigger := range triggers {
			table.Append([]string{
				trigger.Id,
				trigger.Content.Function_id,
				trigger.Kind,
				trigger.Content.Name,
				trigger.Content.Cron_expr,
				trigger.Created_at,
				trigger.Updated_at,
			})
		}
	}

	table.Render()
}

func HandleAddTrigger(trigger *client.Trigger, argument_values *[]string, interactive *bool, pretty *bool) {
	triggerKinds, _ := client.GetKinds()
	isTriggerKindAllowed := false
	var id = &trigger.Content.Function_id
	args, _ := client.GetFunctionByIdArgs(*id)

	if *interactive {
		// Create a scanner to read user input
		scanner := bufio.NewScanner(os.Stdin)

		// Prompt for trigger kind
		fmt.Printf("Enter one of these Trigger Kinds:\n")
		for _, available_triggerKind := range triggerKinds.TriggerKinds {
			fmt.Printf("  ➤ %v\n", available_triggerKind)
		}
		fmt.Printf("Trigger kind: ")
		scanner.Scan()
		trigger.Kind = scanner.Text()
		for _, available_triggerKind := range triggerKinds.TriggerKinds {
			if trigger.Kind == available_triggerKind {
				isTriggerKindAllowed = true
				break
			}
		}
		if !isTriggerKindAllowed {
			fmt.Printf("Invalid trigger kind. Allowed trigger kinds are: %s\n", strings.Join(triggerKinds.TriggerKinds, ", "))
			os.Exit(1)
		}

		// Prompt for trigger name
		fmt.Printf("Enter Trigger name: ")
		scanner.Scan()
		trigger.Content.Name = scanner.Text()

		// Prompt for trigger cron expression
		fmt.Printf("Enter Trigger cron expression: ")
		scanner.Scan()
		trigger.Content.Cron_expr = scanner.Text()

		// Prompt values of the existing arguments
		if len(args) > 0 {
			fmt.Println("Enter Args (one per line, press Enter for each entry; leave an empty line to finish):")
			for _, arg := range args {
				fmt.Printf("  ➤ %s: ", arg)
				scanner.Scan()
				value := scanner.Text()
				trigger.Content.Args = append(trigger.Content.Args, client.Argument{Key: arg, Value: value})
			}
		}
	} else {
		if len(*argument_values) != len(args) {
			fmt.Printf("Invalid number of arguments. Expected %d arguments, got %d\n", len(args), len(*argument_values))
			os.Exit(1)
		}
		if len(*argument_values) > 0 {
			for i, arg := range args {
				trigger.Content.Args = append(trigger.Content.Args, client.Argument{Key: arg, Value: (*argument_values)[i]})
			}
		}
	}

	c, _ := client.NewClient()
	created_trigger, err := c.AddTrigger(*trigger)
	if nil != err {
		fmt.Printf("failed: %s\n", err)
		os.Exit(1)
	}

	if client.GetDefaultFormat() == "json" {
		utils.PrintJson(created_trigger)
	} else if *pretty {
		utils.PrintPretty("Trigger successfully created", *created_trigger)
	} else {
		utils.PrintRow(*created_trigger)
	}
}

func HandleDeleteTrigger(id *string) {
	c, _ := client.NewClient()
	delete_err := c.DeleteTriggerById(*id)
	if nil != delete_err {
		fmt.Printf("failed: %s\n", delete_err)
		os.Exit(1)
	}

	fmt.Printf("Trigger successfully deleted\n")
}

func HandleTruncateTriggers() {
	c, _ := client.NewClient()
	truncate_err := c.TruncateTriggers()
	if nil != truncate_err {
		fmt.Printf("failed: %s\n", truncate_err)
		os.Exit(1)
	}

	fmt.Printf("Triggers successfully truncated\n")
}
