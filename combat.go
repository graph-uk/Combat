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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func listTestsOrderByName(allTestsWithParams []aTestParams) {
	//println(len(allTestsWithParams))
	for _, curTestParams := range allTestsWithParams {
		println(curTestParams.Name)
		//fmt.Printf("%-20s %-20s", "Name:", curTestParams.Name)
		//println()
		println("-------------------------------------------------")
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
		println("-------------------------------------------------")
		fmt.Printf("%-20s ", "Tags:")
		for curTagKey, curTag := range params.Tags {
			print(curTag)
			if curTagKey < len(params.Tags)-1 {
				print(",")
			}
		}
		println()
		println()
		//curTestParams.paramsUnmarshaled = params
	}
}

func listTestsOrderByTag(allTestsWithParams []aTestParams) {
	var allTags map[string][]string
	allTags = make(map[string][]string)

	for _, curTestParams := range allTestsWithParams {
		for _, curTag := range curTestParams.paramsUnmarshaled.Tags {
			allTags[curTag] = append(allTags[curTag], curTestParams.Name)
		}
	}

	for curTagKey, curTagTests := range allTags {
		fmt.Printf("%s(%d)\r\n", curTagKey, len(curTagTests))
		for _, curTagTest := range curTagTests {
			println(curTagTest)
		}
		println()
	}
}

func listTestsOrderByParameter(allTestsWithParams []aTestParams) {
	var allParametersTests map[string][]string
	allParametersTests = make(map[string][]string)

	var allParametersVariants map[string][]string
	allParametersVariants = make(map[string][]string)

	for _, curTestParams := range allTestsWithParams {
		for _, curParameter := range curTestParams.paramsUnmarshaled.Params {
			allParametersTests[curParameter.Name] = append(allParametersTests[curParameter.Name], curTestParams.Name)
			if curParameter.Type == "EnumParam" {
				for _, curVariant := range curParameter.Variants {
					if !stringInSlice(curVariant, allParametersVariants[curParameter.Name]) {
						allParametersVariants[curParameter.Name] = append(allParametersVariants[curParameter.Name], curVariant)
					}
				}
			}
		}
	}

	for curParameterKey, curParameter := range allParametersTests {
		print(curParameterKey)
		if len(allParametersVariants[curParameterKey]) > 1 {
			print("(")
			for curVariantKey, curVariant := range allParametersVariants[curParameterKey] {
				print(curVariant)
				if curVariantKey < len(allParametersVariants[curParameterKey])-1 {
					print(",")
				}
			}
			print(")")
		}
		println()
		println("-------------------------------------------------")
		for _, curParameterTest := range curParameter {
			print(curParameterTest)
			if len(allParametersVariants[curParameterKey]) > 1 {

			}
			println()
		}
		println()
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

// Run all tests using parameters
// Return count of failed tests
func runTestsUsingParameters(testList []aTestParams, CLIFlags []CLIFlag) int {

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

	errorCount := 0
	for _, curTest := range testList {
		println(curTest.Name)

		params := []string{"run"}
		params = append(params, "./Tests/Tests/"+curTest.Name+`/`+"main.go")
		for _, curParameterKey := range curTest.paramsUnmarshaled.Params { // collect parameters needed for the test.
			params = append(params, "-"+curParameterKey.Name+"="+getCLIFlagValueByName(CLIFlags, curParameterKey.Name))
		}
		cmd := exec.Command("go", params...)
		cmd.Env = os.Environ()
		var out bytes.Buffer
		var outerr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &outerr
		err := cmd.Run()

		println(string(out.Bytes()))
		if err != nil {
			errorCount++
			println(string(outerr.Bytes()))
			println("Exit code: Error")
		} else {
			println("Exit code: Ok")
		}
		println("")
		//os.Exit(0)
	}
	return errorCount
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
	case "tags":
		fullTestList := getFullTestList()
		testListFilteredByName := selectTestsByName(fullTestList, getCLIFlagValueByName(allCLIFlags, "name"))
		testsWithParamsFilteredByName := loadTestParams(testListFilteredByName)
		testsWithParamsFilteredByNameAndTag := selectTestsByTag(testsWithParamsFilteredByName, getCLIFlagValueByName(allCLIFlags, "tags"))
		listTestsOrderByTag(testsWithParamsFilteredByNameAndTag)
		os.Exit(0)
	case "params":
		fullTestList := getFullTestList()
		testListFilteredByName := selectTestsByName(fullTestList, getCLIFlagValueByName(allCLIFlags, "name"))
		testsWithParamsFilteredByName := loadTestParams(testListFilteredByName)
		testsWithParamsFilteredByNameAndTag := selectTestsByTag(testsWithParamsFilteredByName, getCLIFlagValueByName(allCLIFlags, "tags"))
		listTestsOrderByParameter(testsWithParamsFilteredByNameAndTag)
		os.Exit(0)
	case "run":
		fullTestList := getFullTestList()
		testListFilteredByName := selectTestsByName(fullTestList, getCLIFlagValueByName(allCLIFlags, "name"))
		testsWithParamsFilteredByName := loadTestParams(testListFilteredByName)
		testsWithParamsFilteredByNameAndTag := selectTestsByTag(testsWithParamsFilteredByName, getCLIFlagValueByName(allCLIFlags, "tags"))
		errorCount := runTestsUsingParameters(testsWithParamsFilteredByNameAndTag, allCLIFlags)
		println("Total failed:", errorCount)
		os.Exit(errorCount)
	default:
		println("Incorrect action. Please run combat help for find available actions.")
		os.Exit(1)
	}
	os.Exit(0)

	//Парсим все значения параметров из CLI.

	//Сначала фильтруем тесты по имени, потом по тегам.
	//combat run -name="lynx" -tags="xnd"
	//-name
	//-tags
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
