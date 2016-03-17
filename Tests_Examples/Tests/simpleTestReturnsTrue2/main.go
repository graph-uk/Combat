package main

import (
	"TB2/Tests/Tests_shared/aTest"
	"log"
	"os"
)

type theTest struct {
	aTest  aTest.ATest
	params struct {
		HostName         aTest.StringParam
		SessionTimestamp aTest.StringParam
		Locale           aTest.EnumParam
	}
}

func createNewTest() (*theTest, error) {
	var result theTest

	result.params.Locale.AcceptedValues = append(result.params.Locale.AcceptedValues, "EN")
	result.params.Locale.AcceptedValues = append(result.params.Locale.AcceptedValues, "RU")
	result.params.Locale.AcceptedValues = append(result.params.Locale.AcceptedValues, "US")
	//result.params.AdminName.Value = "TestDefaultValue"
	result.aTest.Tags = append(result.aTest.Tags, "NotForLive")
	result.aTest.Tags = append(result.aTest.Tags, "Lynx")
	result.aTest.Tags = append(result.aTest.Tags, "AlwaysFailedTest")

	result.aTest.FillParamsFromCLI(&result.params)
	return &result, nil
}

func main() {
	_, err := createNewTest()
	if err != nil {
		panic(err)
	}
	os.Exit(12)
	log.Println("ok")
	return
}