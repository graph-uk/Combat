package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// This is the base struct contain all required in all test fields
type aTestParams struct {
	Name              string
	ParamsJSON        []byte
	paramsUnmarshaled UnmarshaledTestParams
}

type param struct {
	Name     string
	Type     string
	Variants []string
}

type CLIFlag struct {
	Name  string
	Value string
}

type UnmarshaledTestParams struct {
	Params []param
	Tags   []string
}

func listTestsOrderByName(allTestsWithParams []aTestParams) {
	//println(len(allTestsWithParams))
	for _, curTestParams := range allTestsWithParams {
		println(curTestParams.Name)
		var params UnmarshaledTestParams
		if err := json.Unmarshal(curTestParams.ParamsJSON, &params); err != nil {
			log.Println("Cannot parse json for test: " + curTestParams.Name)
			panic(err)
		}
		for _, curParam := range params.Params {
			fmt.Printf("%-20s %-20s", curParam.Name, curParam.Type)
			if curParam.Type == "EnumParam" {
				for _, curEnumVariant := range curParam.Variants {
					print(curEnumVariant + " ")
				}
			}
			println()
		}
		println()
		//curTestParams.paramsUnmarshaled = params
	}
}

func loadTestParams(testNames []string) []aTestParams {
	// collect all test params
	var allTestsWithParams []aTestParams
	for _, curTestFile := range testNames {
		cmd := exec.Command("go", "run", "./Tests/Tests/"+curTestFile+`/`+"main.go", "paramsJSON")
		//cmd := exec.Command("cmd", "-c", "go run ./Tests/Tests/"+curTestFile+`/`+"main.go paramsJSON")
		cmd.Env = os.Environ()
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()

		var curTestParams aTestParams
		curTestParams.Name = curTestFile
		curTestParams.ParamsJSON = out.Bytes()
		allTestsWithParams = append(allTestsWithParams, curTestParams)
	}

	for _, curTestParams := range allTestsWithParams {
		var params UnmarshaledTestParams
		if err := json.Unmarshal(curTestParams.ParamsJSON, &params); err != nil {
			log.Println("Cannot parse json for test: " + curTestParams.Name)
			panic(err)
		}
		curTestParams.paramsUnmarshaled = params
		//println(params.Params[0].Name)
	}

	for curTestParamsKey, curTestParams := range allTestsWithParams {
		var params UnmarshaledTestParams
		if err := json.Unmarshal(curTestParams.ParamsJSON, &params); err != nil {
			log.Println("Cannot parse json for test: " + curTestParams.Name)
			panic(err)
		}
		allTestsWithParams[curTestParamsKey].paramsUnmarshaled = params
	}

	return allTestsWithParams
}

func getFullTestList() []string {
	testsFileList, err := ioutil.ReadDir("./Tests/Tests")
	if err != nil {
		log.Println("Error: cannot list test's directory.")
		log.Fatal(err)
	}

	var fullTestList []string
	// check that all files is folders
	for _, curTestFile := range testsFileList {
		if !curTestFile.IsDir() {
			log.Fatal("File in tests directory. There is should exist folders only")
		}
		fullTestList = append(fullTestList, curTestFile.Name())
	}
	return fullTestList
}

func selectTestsByName(fullTestList []string, name string) []string {
	var selectedTests []string
	for _, curTest := range fullTestList {
		match, err := regexp.MatchString(name, curTest)
		if err != nil {
			log.Fatal("Incorrect regexp for name")
		}
		if match {
			selectedTests = append(selectedTests, curTest)
		}
	}
	return selectedTests
}

func selectTestsByTag(testList []aTestParams, tag string) []aTestParams {
	var selectedTests []aTestParams
	for _, curTest := range testList {
		tagMatch := false
		for _, curTag := range curTest.paramsUnmarshaled.Tags {
			match, err := regexp.MatchString(tag, curTag)
			if err != nil {
				log.Fatal("Incorrect regexp for tag")
			}
			//tagMatch = true
			if match {
				tagMatch = true
				break
			}
		}
		if tagMatch {
			selectedTests = append(selectedTests, curTest)
		}
	}
	return selectedTests
}

func getCLIFlagValueByName(CLIFlags []CLIFlag, name string) string {
	for _, curFlag := range CLIFlags {
		if curFlag.Name == name {
			return curFlag.Value
		}
	}
	return "" // return empty string if flag not found
}

func isPresentedCLIFlagValueByName(CLIFlags []CLIFlag, name string) bool {
	for _, curFlag := range CLIFlags {
		if curFlag.Name == name {
			return true
		}
	}
	return false
}

func runTestsUsingParameters(testList []aTestParams, CLIFlags []CLIFlag) {

	var allParamsRequiredToTesting map[string]string
	allParamsRequiredToTesting = make(map[string]string)

	for _, curTest := range testList {
		for _, curParameter := range curTest.paramsUnmarshaled.Params {
			if true { //curParameter.Type == "StringParam" {
				allParamsRequiredToTesting[curParameter.Name] = ""
			}
		}
	}

	allRequiredFlagsPresented := true
	for curParameterKey, _ := range allParamsRequiredToTesting {
		if !isPresentedCLIFlagValueByName(CLIFlags, curParameterKey) {
			println("Flag \"" + curParameterKey + "\" is required.")
			allRequiredFlagsPresented = false
		}
	}

	if !allRequiredFlagsPresented {
		os.Exit(1)
	}

	for _, curTest := range testList {
		println(curTest.Name)
		cmd := exec.Command("go", "run", "./Tests/Tests/"+curTest.Name+`/`+"main.go", "-HostName=http://aas-uat.graph.uk", "-SessionTimestamp=123", "-Locale=EN", "-AdminName=AdminName")
		cmd.Env = os.Environ()
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()

		//var curTestParams aTestParams
		//curTestParams.Name = curTestFile
		//curTestParams.ParamsJSON = out.Bytes()
		println(string(out.Bytes()))
		if err != nil {
			println("Exit code: Error")
		} else {
			println("Exit code: Ok")
		}
		println("")
		//os.Exit(0)
	}
}

func parseAllCLIFlags() []CLIFlag {
	var allFlags []CLIFlag
	re := regexp.MustCompile("-.*=")
	for _, curArgument := range os.Args {
		flagName := re.FindString(curArgument)
		if len(flagName) > 2 {
			flagName = flagName[:len(flagName)-1] // trim last character
			flagName = flagName[1:]               // trim first character
			flagName = strings.TrimSpace(flagName)
			var curFlag CLIFlag
			curFlag.Name = flagName
			allFlags = append(allFlags, curFlag)
		}

	}
	for curFlagIndex, curFlag := range allFlags {
		flag.StringVar(&allFlags[curFlagIndex].Value, curFlag.Name, "", "variant")
	}

	flag.Parse()
	return allFlags
}

func main() {
	allCLIFlags := parseAllCLIFlags()

	command := "run"
	if len(os.Args) > 1 {
		command = os.Args[len(os.Args)-1]
	}

	switch command {
	case "help":
		println("Help")
		os.Exit(0)
	case "list":
		fullTestList := getFullTestList()
		testListFilteredByName := selectTestsByName(fullTestList, getCLIFlagValueByName(allCLIFlags, "name"))
		testsWithParamsFilteredByName := loadTestParams(testListFilteredByName)
		testsWithParamsFilteredByNameAndTag := selectTestsByTag(testsWithParamsFilteredByName, getCLIFlagValueByName(allCLIFlags, "tags"))
		listTestsOrderByName(testsWithParamsFilteredByNameAndTag)
		os.Exit(0)
	case "run":
		fullTestList := getFullTestList()
		testListFilteredByName := selectTestsByName(fullTestList, getCLIFlagValueByName(allCLIFlags, "name"))
		testsWithParamsFilteredByName := loadTestParams(testListFilteredByName)
		testsWithParamsFilteredByNameAndTag := selectTestsByTag(testsWithParamsFilteredByName, getCLIFlagValueByName(allCLIFlags, "tags"))
		runTestsUsingParameters(testsWithParamsFilteredByNameAndTag, allCLIFlags)
		os.Exit(0)
	default:
		println("Incorrect action. Please run combat help for find available actions.")
		os.Exit(1)
	}
	os.Exit(0)

	cmdFlag := flag.String("cmd", "run", "action (run/help/list...)")
	nameFlag := flag.String("name", "", "action (run/help/list...)")
	tagsFlag := flag.String("tags", "", "action (run/help/list...)")
	//paramsFlag := flag.String("name", "", "action (run/help/list...)")

	flag.Parse()
	defaultCmd := "run"
	if *cmdFlag == "" {
		cmdFlag = &defaultCmd
	}
	switch *cmdFlag {
	case "help":
		println("Help")
		os.Exit(0)
	case "list":
		fullTestList := getFullTestList()
		selectTestsByName(fullTestList, *nameFlag)
		testListFilteredByName := selectTestsByName(fullTestList, *nameFlag)
		testsWithParamsFilteredByName := loadTestParams(testListFilteredByName)

		testsWithParamsFilteredByNameAndTag := selectTestsByTag(testsWithParamsFilteredByName, *tagsFlag)
		listTestsOrderByName(testsWithParamsFilteredByNameAndTag)
		os.Exit(0)
	case "run":
		println("Run")
		os.Exit(0)
	default:
		println("Incorrect action. Please run combat help for find available actions.")
		os.Exit(1)
	}

	//Парсим все значения параметров из CLI.

	//Сначала фильтруем тесты по имени, потом по тегам.
	//combat run -name="lynx" -tags="xnd"
	//-name
	//-nameReg
	//-tags
	//-tagsReg
	//-locale="sdf"
	//-nyx="sdf"

	//Потом собираем инфу о параметрах, которые нужно предоставить.
	//Читаем параметры.
	//Проверяем, что все необходимые параметры предоставлены. Если параметров предоставлено больше - ошибка.
	//Проверяем, что перечисления вписываются в ограничения.
	//Запускаем/выводим статистику.

	//combat <...> run  запуск
	//combat <...> list вывод тестов с параметрами и тегами, группировка по имени теста.
	//combat <...> params вывод тестов с параметрами и тегами, группировка по параметрам.
	//combat <...> tags вывод тестов с параметрами и тегами, группировка по тегам.
	//combat help вывод справки.

	//Print all tests with params

	return
}
