package main

import (
	"combat/Tests_Examples/Tests/Tests_shared/aTest"
	"combat/Tests_Examples/Tests/Tests_shared/aas"
	"log"
	"os"
)

type theTest struct {
	aTest  aTest.ATest
	params struct {
		HostName         aTest.StringParam
		SessionTimestamp aTest.StringParam
		Locale           aTest.EnumParam
		AdminName        aTest.StringParam
	}
	aas aas.AAS
}

func createNewTest() (*theTest, error) {
	var result theTest

	result.params.Locale.AcceptedValues = append(result.params.Locale.AcceptedValues, "EN")
	result.params.Locale.AcceptedValues = append(result.params.Locale.AcceptedValues, "RU")
	result.params.AdminName.Value = "TestDefaultValue"
	result.aTest.Tags = append(result.aTest.Tags, "NotForLive")

	result.aTest.FillParamsFromCLI(&result.params)
	result.aas.Init(result.params.HostName.Value)
	//	result.aas.Browser.RestartIfDied()
	//	result.aas.Pages.MainPage.Open()
	//	os.Exit(0)
	return &result, nil
}

func main() {
	fillLineTest, err := createNewTest()
	if err != nil {
		panic(err)
	}

	fillLineTest.aas.Browser.RestartIfDied()
	fillLineTest.aas.Pages.MainPage.Open()
	fillLineTest.aas.Pages.MainPage.Header_ClickSignUp()
	fillLineTest.aas.Pages.SignUpPage.FillFirstName(fillLineTest.params.AdminName.Value)
	log.Println("ok")
	os.Exit(0)
}
