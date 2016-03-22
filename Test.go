package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
//	"fmt"
)

type Test struct {
	directory string
	name      string
	params    map[string]TestParameter
	tags      []string
}

type TestParameter struct {
	Name     string
	Type     string
	Variants []string
}

func (t *Test) LoadTagsAndParams() error {
	type UnmarshaledTestParams struct {
		Params []TestParameter
		Tags   []string
	}

	// get test's params in JSON

	cmd := exec.Command("go", "run", t.directory+"/"+t.name+`/`+"main.go", "paramsJSON")
	cmd.Env = os.Environ()
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//println("LoadTagsAndParams: " + out.String())
	var TestParams UnmarshaledTestParams
	if err := json.Unmarshal(out.Bytes(), &TestParams); err != nil {
		log.Println("Cannot parse json for test: " + t.name)
		log.Println("JSON data: " + out.String())
		panic(err)
	}
	t.tags = TestParams.Tags

	for _, curParameter := range TestParams.Params {
		t.params[curParameter.Name] = curParameter
	}

	return nil
}


func (t *Test) IsCasesEqual(case1 []string, case2 []string) bool{
	if len(case1)!=len(case2){
		return false
	}

	result := true
	for _, curParameter := range case1{
		parameterFound := false
		for _, curParameter2 := range case2{
			if curParameter == curParameter2{
				parameterFound = true
				break
			}
		}
		if !parameterFound{
			return false
		}
	}
	return result
}

func (t *Test) IsCasePresented(allCases [][]string, aCase []string) bool{
	for _, curCase := range allCases{
		if t.IsCasesEqual(curCase,aCase){
			return true
		}
	}
	return false
}

func (t *Test) GetCasesByParameterCombinations(paramCombinations []*map[string]string) [][]string {
	var result [][]string

	for _, curCombination := range paramCombinations{
		curCombinationAccepted := true
		curCombinationCase := []string{t.name}
		for nameOfcurParamOfTest, curParamOfTest := range t.params{
			if curParamOfTest.Type == "EnumParam" {
				if !stringInSlice((*curCombination)[nameOfcurParamOfTest], curParamOfTest.Variants) {
					curCombinationAccepted = false
					break
				}
			}
		}
		if curCombinationAccepted{
			for nameOfcurParamOfTest, _ := range t.params{
				curCombinationCase = append(curCombinationCase, "-"+nameOfcurParamOfTest+"="+(*curCombination)[nameOfcurParamOfTest])
			}
			if !t.IsCasePresented(result,curCombinationCase){
				result = append(result,curCombinationCase)
			}
		}
	}
	return result
}