package anonGoToMainpage

import (
	"TestBot/aas"
	"TestBot/testsEnumirator/aTest"
	//	"log"
	//"time"
)

type AnonGoToMainpage struct {
	aTest.ATest
	//aas *aas.AAS
}

func (a *AnonGoToMainpage) Init(aas *aas.AAS) error {
	a.ATest.Init()
	a.AAS = aas
	a.Name = "AnonGoToMainpage"
	return nil
}

func (a *AnonGoToMainpage) Run() {
	defer func() {
		if r := recover(); r != nil {
			a.Success = false
		}
	}()
	a.AAS.Browser.RestartIfDied()

	a.AAS.Pages.MainPage.Open()
	a.AAS.Browser.Log += "Click all header's menu\r\n"
	a.AAS.Pages.MainPage.Header_ClickLogIn()
	a.AAS.Pages.MainPage.Header_ClickResetPassword(true)
	a.AAS.Pages.MainPage.Header_ClickSignUp()
	a.Success = true
}
