package testsEnumerator

import (
	"TestBot/aas"
	//"TestBot/browser"
	//"TestBot/testsEnumirator/tests/anonCreateNewAccount"
	"TestBot/testsEnumirator/tests/anonGoToMainpage"
	"log"
)

type aTestInterface interface {
	Init(*aas.AAS) error
	GetName() *string
	GetLocales() *[]string
	GetTags() *[]string
	GetLog() string
	IsEnabled() bool
	Enable()
	Disable()
	Run()
	WriteLog(string)
	IsSuccess() bool
}

type TestsEnumerator struct {
	TestsArray []aTestInterface
	Log        string
}

func (b *TestsEnumerator) MockEnableTestIfShouldRun(t aTestInterface) {
	t.Enable()
}

func (b *TestsEnumerator) Init(aas *aas.AAS) error {
	b.TestsArray = append(b.TestsArray, new(anonGoToMainpage.AnonGoToMainpage))
	//b.TestsArray = append(b.TestsArray, new(anonCreateNewAccount.AnonCreateNewAccount))

	//init
	for _, curTest := range b.TestsArray {
		curTest.Init(aas)
	}

	//enable
	for _, curTest := range b.TestsArray {
		b.MockEnableTestIfShouldRun(curTest)
	}

	//print list of test that should run
	for _, curTest := range b.TestsArray {
		if curTest.IsEnabled() {
			log.Println(*curTest.GetName())
		}
	}
	return nil
}

func (b *TestsEnumerator) RunEnabledTests() error {
	for _, curTest := range b.TestsArray {
		if curTest.IsEnabled() {
			curTest.Run()
			curTest.WriteLog("D:\\" + *curTest.GetName() + ".txt")
			if curTest.IsSuccess() {
				log.Println("Success")
			} else {
				log.Println("Fail")
			}
		}
	}
	return nil
}
