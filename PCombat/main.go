package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/graph-uk/Combat/PCombat/combatClient"
)

func main() {
	defaultSessionTimeout := 60 //minutes

	client, err := combatClient.NewCombatClient()
	if err != nil {
		fmt.Println("Cannot init combat client")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	sessionID, err := client.CreateNewSession(defaultSessionTimeout)
	if err != nil {
		fmt.Println("Cannot create combat session")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	failCount := client.GetSessionResult(sessionID)

	if failCount == 0 {
		fmt.Println("All tests are passed")
		os.Exit(0)
	} else {
		fmt.Println("Total failed tests: " + strconv.Itoa(failCount))
		os.Exit(failCount)
	}
}
