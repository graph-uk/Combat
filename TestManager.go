package main

import (
	"Combat/CLIParser"
	"io/ioutil"
	"log"
	"regexp"
)

// This is the base struct contain all required in all test fields
type TestManager struct {
	tests                map[string]Test
	CLIParameters        map[string]string
	testMergedParameters TestParameter
}

// Parse all parameters from CLI. Fill default values if needed.
func (t *TestManager) parseCLIParameters() {
	t.CLIParameters = CLIParser.ParseAllCLIFlags()
	if _, ok := t.CLIParameters["name"]; !ok {
		t.CLIParameters["name"] = ""
	}

	if _, ok := t.CLIParameters["tag"]; !ok {
		t.CLIParameters["tag"] = ""
	}
}

func (t *TestManager) Init(directory string, params map[string]string) error {
	t.parseCLIParameters()
	t.SelectAllTests(directory)
	t.FilterTestsByName(t.CLIParameters["name"])
	t.FilterTestsByTag(t.CLIParameters["tag"])

	for curTestName, _ := range t.tests {
		println(curTestName)
	}

	return nil //a.AAS.Browser.Log
}

func (t *TestManager) PrintListOrderedByNames() {
}

func (t *TestManager) PrintListOrderedByTags() {
}

func (t *TestManager) PrintListOrderedByParams() {
}

func (t *TestManager) RunTests() {
}

//Select all tests in the directory, load that's parameters, and collect it to t.tests
func (t *TestManager) SelectAllTests(directory string) error {
	// clear test list
	t.tests = make(map[string]Test)

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
		t.tests[curTestFile.Name()] = Test{
			directory: directory,
			name:      curTestFile.Name(),
			params:    map[string]TestParameter{},
			tags:      []string{},
		}
	}

	for _, curTest := range t.tests {
		curTest.LoadTagsAndParams()
		//GO BUG?
		//This line prints correct values, but outside of cycle - zero len
		//println(len(curTest.tags))
	}
	// load params of each test in t.test

	//GO BUG?
	//This line prints correct values, but outside of cycle - zero len
	//println(len(curTest.tags))

	for _, curTest := range t.tests {
		println(len(curTest.tags))
	}

	return nil
}

// Saves in t.tests the tests with a suitable name only
func (t *TestManager) FilterTestsByName(name string) error {
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
func (t *TestManager) FilterTestsByTag(tag string) error {
	for curTestName, curTest := range t.tests {
		tagFound := false
		for _, curTag := range curTest.tags {
			println("xcvxcv")
			println(curTag)
			println(tag)
			match, err := regexp.MatchString(tag, curTag)
			println(match)
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
