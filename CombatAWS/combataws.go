package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func packTests() string {
	fmt.Println("Packing tests")
	tmpFile, err := ioutil.TempFile("", "combatSession")
	if err != nil {
		panic(err)
	}
	tmpFile.Close()
	//	fmt.Println(tmpFile.Name())
	//	tmpFile.Close()
	zipit("./..", tmpFile.Name())
	return tmpFile.Name()
}

func cleanupTests() {
	fmt.Println("Cleanup tests")
}

func getParams() string {
	params := ""
	for curArgIndex, curArg := range os.Args {
		if curArgIndex > 0 {
			params += curArg + " "
		}
	}
	return params
}

func createSessionOnServer(archiveFileName string) (string, error) {
	fmt.Println("Uploading session...")
	sessionName := ""

	var err error
	for i := 1; i <= 10; i++ {
		sessionName, err = postSession(archiveFileName, getParams(), "http://localhost:9090/createSession")
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
	return sessionName, err
}

func main() {
	cleanupTests()
	testsArchiveFileName := packTests()
	sessionName, _ := createSessionOnServer(testsArchiveFileName)
	fmt.Println("Session: " + sessionName)
	for {
		finished, _, err := getSessionStatusJSON(sessionName)
		if err != nil {
			fmt.Println(err.Error())
		}
		if finished {
			break
		}
		time.Sleep(5 * time.Second)
	}
}
