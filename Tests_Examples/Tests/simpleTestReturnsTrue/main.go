package main

import (
	"Combat/Tests_Examples/Tests_shared/aTest"
	"log"
)

type theTest struct {
	aTest  aTest.ATest
	params struct {
		HostName         aTest.StringParam
		SessionTimestamp aTest.StringParam
		Locale           aTest.EnumParam
		AdminName        aTest.StringParam
	}
}

func createNewTest() (*theTest, error) {
	var result theTest

	result.params.Locale.AcceptedValues = append(result.params.Locale.AcceptedValues, "EN")
	result.params.Locale.AcceptedValues = append(result.params.Locale.AcceptedValues, "RU")
	result.params.AdminName.Value = "TestDefaultValue"
	result.aTest.Tags = append(result.aTest.Tags, "NotForLive")

	result.aTest.FillParamsFromCLI(&result.params)
	return &result, nil
}

func main() {
	_, err := createNewTest()
	if err != nil {
		panic(err)
	}

	log.Println("ok")
	return
}
