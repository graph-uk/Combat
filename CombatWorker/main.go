package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func packOutputToTemp() string {
	fmt.Println("Packing output")
	tmpFile, err := ioutil.TempFile("", "combatOutput")
	if err != nil {
		panic(err)
	}
	tmpFile.Close()
	zipit("out", tmpFile.Name())
	return tmpFile.Name()
}

func getJob(host string) (command string, params string, sessionID string) {
	response, err := http.Get(host + "/getJob")
	if err != nil {
		fmt.Println()
		fmt.Printf("%s", err)
	} else {
		fmt.Println("getJob - " + response.Header.Get("command"))
		defer response.Body.Close()
		command = response.Header.Get("Command")
		if command == "idle" {
			return command, "", ""
		}
		params = response.Header.Get("Params")
		sessionID = response.Header.Get("SessionID")
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
		}
		//fmt.Printf("%s\n", string(contents))
		ioutil.WriteFile("./job/archived.zip", contents, 0777)
	}
	return command, params, sessionID
}

func postCases(cases string, sessionID string) error {
	fmt.Print("post cases... ")
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("sessionID", sessionID)
	bodyWriter.WriteField("cases", cases)
	//fileWriter, err := bodyWriter.("uploadfile", filename)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post("http://deltaview.ru:9090/setSessionCases", contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Fail: incorrect request status - " + strconv.Itoa(resp.StatusCode))
		return errors.New("Incorrect request status: " + strconv.Itoa(resp.StatusCode))
	} else {
		fmt.Println("Ok")
	}
	cleanupJob()
	return nil
}

func doCasesExplore(params, sessionID string) (status int, cases string) {
	fmt.Println("CasesExplore")
	err := unzip("./job/archived.zip", "./job/unarch")
	if err != nil {
		fmt.Print(err.Error())
	}
	os.Chdir("job/unarch/Tests")

	command := []string{"cases"}
	fmt.Println(params)
	command = append(command, strings.Split(params, " ")...)

	cmd := exec.Command("combat", command...)
	cmd.Env = os.Environ()
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outErr
	fmt.Println(command)
	fmt.Print("Run combat cases... ")
	exitStatus := cmd.Run()

	if exitStatus == nil {
		fmt.Println("Ok")
		postCases(out.String(), sessionID)
		//fmt.Println(out.String())
	} else {
		fmt.Println("Fail")
		fmt.Println(out.String())
		fmt.Println(outErr.String())
	}
	return 1, ""
}

func doRunCase(params, caseID string) {
	fmt.Println("CaseRunning " + params)
	err := unzip("./job/archived.zip", "./job/unarch")
	if err != nil {
		fmt.Print(err.Error())
	}
	os.Chdir("job/unarch/Tests")

	command := []string{"run"}
	//fmt.Println(params)
	command = append(command, strings.Split(params, " ")...)
	os.Chdir(command[1])
	command[1] = "main.go"

	cmd := exec.Command("go", command...)
	cmd.Env = os.Environ()
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outErr
	//	fmt.Println(command)
	fmt.Print("Run case... ")
	exitStatus := cmd.Run()

	exitStatusString := ""
	if exitStatus == nil {
		exitStatusString = "0"
		fmt.Println("Ok")
		//postCases(out.String(), sessionID)
		//fmt.Println(out.String())
	} else {
		exitStatusString = exitStatus.Error()
		fmt.Println("Fail")
		fmt.Println(out.String())
		fmt.Println(outErr.String())
	}

	postCaseResult(caseID, exitStatusString, out.String()+outErr.String())
	return
}

func postCaseResult(caseID, exitStatus, stdout string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	if _, err := os.Stat("out"); err != nil { // if we have has not "out" directory - create it
		os.MkdirAll("out", 0777)
	}
	outFileName := packOutputToTemp()
	//zipit("out", "out.zip")

	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", outFileName)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// open file handle
	fh, err := os.Open(outFileName)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.WriteField("caseID", caseID)
	bodyWriter.WriteField("exitStatus", exitStatus)
	bodyWriter.WriteField("stdOut", stdout)
	bodyWriter.Close()

	resp, err := http.Post("http://deltaview.ru:9090"+"/setCaseResult", contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
	if resp.StatusCode != 200 {
		return errors.New("Incorrect request status: " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}

func cleanupJob() {
	err := os.RemoveAll("job")
	if err != nil {
		//panic(err)
		//fmt.Println(err.Error())
	}
	time.Sleep(1 * time.Second)
	os.MkdirAll("job", 0777)
}

var myPath string

func main() {
	myPath, _ = os.Getwd()
	for {
		os.Chdir(myPath)
		cleanupJob()
		fmt.Println(".")
		command, params, sessionID := getJob("http://deltaview.ru:9090")
		if command == "CasesExplore" {
			doCasesExplore(params, sessionID)
		}
		if command == "RunCase" {
			doRunCase(params, sessionID)
		}
		if command == "idle" {
			time.Sleep(5 * time.Second)
		}
	}
}
