package main

import (
	"fmt"
	"os"

	"github.com/graph-uk/Combat/CombatWorkerLauncher/LauncherServer"
)

func main() {
	launcherServerServer, err := LauncherServer.NewLauncherServer()
	if err != nil {
		fmt.Println("Cannot init launcher server")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = launcherServerServer.CheckAmazonCredentials()
	if err != nil {
		fmt.Println("Amazon credentials are not valid. Check AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY in config.json")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = launcherServerServer.Serve()
	if err != nil {
		fmt.Println("Cannot serve")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
