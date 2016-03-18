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

func (t *Test) GetCasesByParameterCombinations(paramCombinations []*map[string]string) [][]string {
	var result [][]string

	//fmt.Println(t.name)

	for _, curCombination := range paramCombinations{
		//fmt.Print(*curCombination)
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
			result = append(result,curCombinationCase)
		}
		//fmt.Println()
	}
	return result
}