package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func packTests() string {
	fmt.Println("Packing tests")
	tmpFile, err := ioutil.TempFile("", "combatSession")
	if err != nil {
		panic(err)
	}
	tmpFile.Close()
	fmt.Println(tmpFile.Name())
	tmpFile.Close()
	zipit("./", tmpFile.Name())
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

func createSessionOnServer(sessionName string, archiveFileName string) {
	fmt.Println("Uploading session...")

	var err error
	for i := 1; i <= 10; i++ {
		err = postFile(archiveFileName, sessionName, "http://localhost:9090/uploadSession")
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
}

func main() {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	fmt.Println("Session: " + timestamp)
	fmt.Println(getParams())
	cleanupTests()
	testsArchiveFileName := packTests()
	createSessionOnServer(timestamp, testsArchiveFileName)
}
