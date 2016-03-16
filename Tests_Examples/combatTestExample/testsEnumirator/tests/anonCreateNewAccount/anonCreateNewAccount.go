package anonCreateNewAccount

import (
	"TestBot/aas"
	"TestBot/testsEnumirator/aTest"
	//	"log"
	//"time"
)

type AnonCreateNewAccount struct {
	aTest.ATest
}

func (a *AnonCreateNewAccount) Init(aas *aas.AAS) error {
	a.ATest.Init()
	a.AAS = aas
	a.Name = "anonCreateNewAccount"
	return nil
}

func (a *AnonCreateNewAccount) Run() {
	defer func() {
		if r := recover(); r != nil {
			a.Success = false
		}
	}()
	a.AAS.Browser.RestartIfDied()
	a.AAS.Pages.SignUpPage.FillFirstName("lynx")

	//a.AAS.Pages
	a.Success = true
}
