package main

import (
	"fmt"
	"os"

	"github.com/graph-uk/Combat/CombatServer/server"
)

func main() {
	combatServer, err := server.NewCombatServer()
	if err != nil {
		fmt.Println("Cannot init server server.")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = combatServer.Serve()
	if err != nil {
		fmt.Println("Cannot serve.")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
