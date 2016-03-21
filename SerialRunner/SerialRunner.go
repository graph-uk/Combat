package SerialRunner

import (
	"fmt"
	"os/exec"
	//"os"
	"bytes"
	"os"
)

func RunCasesSerial(cases [][]string, directory string)  int{
	fmt.Println("Running cases.")
	for _, curCase := range cases{
		curCase[0] = directory+"/"+curCase[0]+"/main.go"
		fmt.Println(curCase)
		os.Exit(0)

		cmd := exec.Command("go",curCase[1:]...)
		cmd.Env = os.Environ()
		var out bytes.Buffer
		var outErr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &outErr
		//err :=
		cmd.Run()
		os.Exit(0)
	}

	return 0
}