package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func getJob(host string) (command string, params string) {
	response, err := http.Get(host + "/getJob")
	if err != nil {
		fmt.Println()
		fmt.Printf("%s", err)
	} else {
		fmt.Println("getJob - " + response.Header.Get("command"))
		defer response.Body.Close()
		command = response.Header.Get("command")
		if command == "idle" {
			return command, ""
		}
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
		}
		//fmt.Printf("%s\n", string(contents))
		ioutil.WriteFile("./job/archived.zip", contents, 0777)
	}
	return command, params
}

func doCasesExplore(params string) (status int, cases string) {
	fmt.Println("CasesExplore")
	err := unzip("./job/archived.zip", "./job/unarch")
	if err != nil {
		fmt.Print(err.Error())
	}
	os.Chdir("job/unarch/Tests")
	return 1, ""
}

func main() {
	os.RemoveAll("job")
	time.Sleep(1 * time.Second)
	os.MkdirAll("job", 0777)

	for {
		fmt.Println(".")
		command, params := getJob("http://localhost:9090")
		if command == "CasesExplore" {
			doCasesExplore(params)
		}
		time.Sleep(5 * time.Second)
	}

}
