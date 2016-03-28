package aTest

import (
	"TestBot/aas"
	//"TestBot/browser"
	"io/ioutil"
)

type ATest struct {
	//Browser *browser.Browser
	AAS        *aas.AAS
	Name       string
	Enabled    bool
	Success    bool
	TriesCount int
	Locales    []string
	Tags       []string
}

func (a *ATest) GetLog() string {
	//a.AAS.Browser.ClearCookiesAndStorage()
	//result := a.AAS.Browser.Log
	//return result
	return a.AAS.Browser.Log
}

func (a *ATest) GetName() *string {
	return &a.Name
}

func (a *ATest) GetLocales() *[]string {
	return &a.Locales

}
func (a *ATest) GetTags() *[]string {
	return &a.Tags
}

func (a *ATest) IsEnabled() bool {
	return a.Enabled
}

func (a *ATest) IsSuccess() bool {
	return a.Success
}

func (a *ATest) Enable() {
	a.Enabled = true
}

func (a *ATest) Disable() {
	a.Enabled = false
}

func (a *ATest) Init() {
	a.Disable()
	a.Success = false
	a.TriesCount = 0
}

func (a *ATest) WriteLog(filename string) {
	ioutil.WriteFile(filename, []byte(a.GetLog()), 0644)
}
