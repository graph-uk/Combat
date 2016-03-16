package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
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
	//println("go", "run", t.directory+"/"+t.name+`/`+"main.go", "paramsJSON")
	cmd.Env = os.Environ()
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	var TestParams UnmarshaledTestParams
	if err := json.Unmarshal(out.Bytes(), &TestParams); err != nil {
		log.Println("Cannot parse json for test: " + t.name)
		log.Println("JSON data: " + out.String())
		panic(err)
	}
	//t.params = TestParams.Params
	t.tags = TestParams.Tags
	for _, curParameter := range TestParams.Params {
		t.params[curParameter.Name] = curParameter
	}

	return nil
}
