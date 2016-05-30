package combatClient

import (
	"errors"
	"fmt"
	"os"
	//"strconv"
	"io/ioutil"
	"time"
)

type CombatClient struct {
	serverURL string
	sessionID string
}

func (t *CombatClient) getServerUrlFromCLI() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("Server URL is required")
	}
	return os.Args[1], nil
}

func NewCombatClient() (*CombatClient, error) {
	var result CombatClient
	var err error

	result.serverURL, err = result.getServerUrlFromCLI()
	if err != nil {
		return &result, err
	}

	return &result, nil
}

func (t *CombatClient) packTests() (string, error) {
	fmt.Println("Packing tests")
	tmpFile, err := ioutil.TempFile("", "combatSession")
	if err != nil {
		panic(err)
	}
	tmpFile.Close()
	//	fmt.Println(tmpFile.Name())
	//	tmpFile.Close()
	zipit("./..", tmpFile.Name())
	return tmpFile.Name(), nil
}

func (t *CombatClient) cleanupTests() error {
	fmt.Println("Cleanup tests")
	return nil
}

func (t *CombatClient) getParams() string {
	params := ""
	for curArgIndex, curArg := range os.Args {
		if curArgIndex > 1 {
			params += curArg + " "
		}
	}
	return params
}

func (t *CombatClient) createSessionOnServer(archiveFileName string) string {
	fmt.Println("Uploading session...")
	sessionName := ""

	var err error
	for i := 1; i <= 10; i++ {
		sessionName, err = postSession(archiveFileName, t.getParams(), t.serverURL+"/createSession")
		if err != nil {
			time.Sleep(5 * time.Second)
			fmt.Println(err.Error())
		} else {
			break
		}
	}
	if err != nil {
		fmt.Println("Cannot upload file. Check is server available.")
	}
	return sessionName
}

func (t *CombatClient) CreateNewSession(timeoutMinutes int) (string, error) {
	err := t.cleanupTests()
	if err != nil {
		fmt.Println("Cannot cleanup tests")
		return "", err
	}

	testsArchiveFileName, err := t.packTests()
	if err != nil {
		fmt.Println("Cannot pack tests to zip archive")
		return "", err
	}

	sessionName := t.createSessionOnServer(testsArchiveFileName)
	fmt.Println("Session: " + sessionName)

	//	cleanupTests()
	//	testsArchiveFileName := packTests()
	//	sessionName, _ := createSessionOnServer(testsArchiveFileName)
	//	fmt.Println("Session: " + sessionName)
	return sessionName, nil
}

func (t *CombatClient) GetSessionResult(sessionID string) int {
	for {
		finished, _, err := getSessionStatusJSON(sessionID)
		if err != nil {
			fmt.Println(err.Error())
		}
		if finished {
			break
		}
		time.Sleep(5 * time.Second)
	}
	return 0
}
