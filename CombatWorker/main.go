package main

import (
	"fmt"
	"os"

	"github.com/graph-uk/Combat/CombatWorker/worker"
)

func main() {

	worker, err := combatWorker.NewCombatWorker()
	if err != nil {
		fmt.Println("Cannot init combat worker")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for {
		worker.Work()
	}

}
