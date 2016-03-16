package main

import (
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

func (t *TestManager) Init(directory string, params map[string]string) error {
	t.SelectAllTests(directory)

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

	// load params of each test in t.test
	for _, curTest := range t.tests {
		curTest.LoadTagsAndParams()
	}

	return nil
}

// Saves in t.tests the tests with a suitable name only
func (t *TestManager) FilterTestsByName(name string) error {
	for curTestName, _ := range t.tests {
		match, err := regexp.MatchString(curTestName, name)
		if err != nil {
			log.Fatal("Incorrect regexp for name")
		}
		if !match {
			delete(t.tests, curTestName)
		}
	}
	return nil
}

//func selectTestsByName(fullTestList []string, name string) []string {
//	var selectedTests []string
//	for _, curTest := range fullTestList {
//		match, err := regexp.MatchString(name, curTest)
//		if err != nil {
//			log.Fatal("Incorrect regexp for name")
//		}
//		if match {
//			selectedTests = append(selectedTests, curTest)
//		}
//	}
//	return selectedTests
//}
