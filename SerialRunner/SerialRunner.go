package SerialRunner

import (
	"fmt"
	"os/exec"
	"bytes"
	"os"
	"strings"
)

func addLeftTab(str string) string  {
	result := ""
	strArray := strings.Split(str, "\n")
	for _, curStr := range strArray{
		result += "    "+strings.TrimSpace(curStr)+"\r\n"
	}
	result = "    "+strings.TrimSpace(result)
	return result
}

func RunCasesSerial(cases [][]string, directory string)  int{
	fmt.Println("Run cases.")
	FailedCasesCount := 0
	for _, curCase := range cases{
		curCaseTestName := curCase[0]
		curCase[0] = directory+"/"+curCase[0]+"/main.go"
		curCase = append([]string{"run"},curCase...)
		//fmt.Println(curCase)
		//os.Exit(0)
		//fmt.Println("sdf")
		cmd := exec.Command("go",curCase...)
		cmd.Env = os.Environ()
		var out bytes.Buffer
		var outErr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &outErr
		//fmt.Println(curCase[1:])
		fmt.Println(curCaseTestName)
		exitStatus := cmd.Run()

		if exitStatus !=nil{
			FailedCasesCount++
			fmt.Println(addLeftTab(exitStatus.Error()))
			fmt.Println(addLeftTab(out.String()))
			fmt.Println(addLeftTab(outErr.String()))
		}else{
			fmt.Println(addLeftTab("          OK"))
		}
		fmt.Println()
	}
	fmt.Println("Total failed cases: ", FailedCasesCount)
	return 0
}