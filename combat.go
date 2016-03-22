package main

import (
	"Combat/CLIParser"
	"os"
	"Combat/SerialRunner"
)

func main() {
	action := CLIParser.GetAction() //"run" action by default
	if action == "" {
		action = "run"
	}

	if action == "help" {
		println("Help")
		os.Exit(0)
	}

	var testManager TestManager

	testManager.Init("Tests_Examples/Tests")

	switch action {
	case "list":
		testManager.PrintListOrderedByNames()
	case "tags":
		testManager.PrintListOrderedByTag()
	case "params":
		testManager.PrintListOrderedByParameter()
	case "cases":
		testManager.PrintCases()
	case "run":
		testManager.PrintCases()
		totalFailed := SerialRunner.RunCasesSerial(testManager.AllCases(),"Tests_Examples/Tests")
		os.Exit(totalFailed)
	default:
		println("Incorrect action. Please run \"Combat help\" for find available actions.")
		os.Exit(1)
	}
	os.Exit(0)
}
