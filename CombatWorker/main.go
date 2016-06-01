package main

import (
	"fmt"
	"os"

	"github.com/graph-uk/Combat/CombatWorker/worker"
)

func main() {
	//	for _, CurVal := range os.Environ() {
	//		fmt.Println(CurVal)
	//	}
	//	//fmt.Println(os.Environ())
	//	fmt.Println(string(os.PathListSeparator))
	//	os.Exit(0)
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
