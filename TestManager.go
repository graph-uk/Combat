package main

import (
	"Combat/CLIParser"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

// This is the base struct contain all required in all test fields
type TestManager struct {
	tests                map[string]*Test
	parametersFromCLI    map[string]string
	testMergedParameters TestParameter
}
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Parse all parameters from CLI. Fill default values if needed.
func (t *TestManager) parseCLIParameters() {
	t.parametersFromCLI = CLIParser.ParseAllCLIFlags()
	if _, ok := t.parametersFromCLI["name"]; !ok {
		t.parametersFromCLI["name"] = ""
	}

	if _, ok := t.parametersFromCLI["tag"]; !ok {
		t.parametersFromCLI["tag"] = ""
	}
}

func (t *TestManager) Init(directory string, params map[string]string) error {
	t.parseCLIParameters()
	t.selectAllTests(directory)
	t.filterTestsByName()
	t.filterTestsByTag()
	return nil //a.AAS.Browser.Log
}


//Select all tests in the directory, load that's parameters, and collect it to t.tests
func (t *TestManager) selectAllTests(directory string) error {
	// clear test list
	t.tests = make(map[string]*Test)

	// read test's directory
	testsFileList, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Println("Error: cannot list directory: " + directory)
		log.Fatal(err)
	}

	// check that no files in the test's directory
	for _, curTestFile := range testsFileList {
		if !curTestFile.IsDir() {
			log.Fatal("File " + curTestFile.Name() + " in tests directory: " + directory + ". There is should exist folders only.")
		}
	}

	// create new items in t.tests,
	for _, curTestFile := range testsFileList {
		t.tests[curTestFile.Name()] = &Test{
			directory: directory,
			name:      curTestFile.Name(),
			params:    map[string]TestParameter{},
			tags:      []string{},
		}
	}

	for _, curTest := range t.tests {
		curTest.LoadTagsAndParams()
	}
	return nil
}

// Saves in t.tests the tests with a suitable name only
func (t *TestManager) filterTestsByName() error {
	name := t.parametersFromCLI["Name"]
	for curTestName, _ := range t.tests {
		match, err := regexp.MatchString(name, curTestName)
		if err != nil {
			log.Fatal("Incorrect regexp in name parameter")
		}
		if !match {
			delete(t.tests, curTestName)
		}
	}
	return nil
}

// Saves in t.tests the tests with a suitable tag only
func (t *TestManager) filterTestsByTag() error {
	tag := t.parametersFromCLI["tag"]
	for curTestName, curTest := range t.tests {
		tagFound := false
		for _, curTag := range curTest.tags {
			match, err := regexp.MatchString(tag, curTag)
			if err != nil {
				log.Fatal("Incorrect regexp in name parameter")
			}
			if match {
				tagFound = true
				break
			}
		}
		if !tagFound {
			delete(t.tests, curTestName)
		}
	}
	return nil
}

// Print to STDOUT list of tests ordered by name
func (t *TestManager) PrintListOrderedByNames() error {
	for _, curTest := range t.tests {
		fmt.Println(curTest.name)
		fmt.Println("-------------------------------------------------")

		for _, curParam := range curTest.params {
			fmt.Printf("%-20s %-20s", curParam.Name, curParam.Type)
			if curParam.Type == "EnumParam" {
				for _, curEnumVariant := range curParam.Variants {
					fmt.Print(curEnumVariant + " ")
				}
			}
			fmt.Println()
		}
	}
	return nil
}

// Print to STDOUT list of tests ordered by tag
func (t *TestManager) PrintListOrderedByTag() error {
	var allTags map[string][]string
	allTags = make(map[string][]string)

	for _, curTest := range t.tests {
		for _, curTag := range curTest.tags {
			allTags[curTag] = append(allTags[curTag], curTest.name)
		}
	}
	for curTagKey, curTagTests := range allTags {
		fmt.Printf("%s(%d)\r\n", curTagKey, len(curTagTests))
		for _, curTagTest := range curTagTests {
			fmt.Println(curTagTest)
		}

		fmt.Println()
	}
	return nil
}

// Print to STDOUT list of tests ordered by parameter
func (t *TestManager) PrintListOrderedByParameter() error {
		var allParametersTests map[string][]string
		allParametersTests = make(map[string][]string)

		var allParametersVariants map[string][]string
		allParametersVariants = make(map[string][]string)

		for _, curTest := range t.tests {
			for _, curParameter := range curTest.params {
				allParametersTests[curParameter.Name] = append(allParametersTests[curParameter.Name], curTest.name)
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
			fmt.Print(curParameterKey)
			if len(allParametersVariants[curParameterKey]) > 1 {
				fmt.Print("(")
				for curVariantKey, curVariant := range allParametersVariants[curParameterKey] {
					fmt.Print(curVariant)
					if curVariantKey < len(allParametersVariants[curParameterKey])-1 {
						fmt.Print(",")
					}
				}
				fmt.Print(")")
			}
			fmt.Println()
			fmt.Println("-------------------------------------------------")
			for _, curParameterTest := range curParameter {
				if t.tests[curParameterTest].params[curParameterKey].Type == "EnumParam"{
					fmt.Println(curParameterTest, t.tests[curParameterTest].params[curParameterKey].Variants)
				}else{
					fmt.Println(curParameterTest)
				}

			}
			fmt.Println()
		}
	return nil
}


// Print to STDOUT all cases are allowed for this parameters combination
func (t *TestManager) PrintCases() error {
	var allEnumParameters map[string][]string
	allEnumParameters = make(map[string][]string)

	for _, curTest := range t.tests {
		for _, curParameter := range curTest.params {
			if curParameter.Type == "EnumParam" {
				allEnumParameters[curParameter.Name] = []string{}
				for _, curVariant := range curParameter.Variants {
					if !stringInSlice(curVariant, allEnumParameters[curParameter.Name]) {
						allEnumParameters[curParameter.Name] = append(allEnumParameters[curParameter.Name], curVariant)
					}
				}
			}
		}
	}

	for curParameterName, curParameterVariants := range allEnumParameters{
		fmt.Println(curParameterName, curParameterVariants)
	}
	return nil
}