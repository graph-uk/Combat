package main

import (
	"os"

	"github.com/graph-uk/combat/CLIParser"
	"github.com/graph-uk/combat/Manual"
	"github.com/graph-uk/combat/SerialRunner"
)

func main() {
	action := CLIParser.GetAction() //"run" action by default
	if action == "" {
		action = "run"
	}

	if action == "help" {
		Manual.PrintManual()
		os.Exit(0)
	}

	var testManager TestManager

	curDirectory, _ := os.Getwd()
	testManager.Init(curDirectory)

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
		totalFailed := SerialRunner.RunCasesSerial(testManager.AllCases(), curDirectory)
		os.Chdir(curDirectory)
		os.Exit(totalFailed)
	default:
		println("Incorrect action. Please run \"Combat help\" for find available actions.")
		os.Exit(1)
	}
	os.Exit(0)
}
