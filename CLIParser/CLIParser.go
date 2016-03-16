package CLIParser

import (
	"flag"
	"os"
	"regexp"
	"strings"
)

type CLIFlag struct {
	Name  string
	Value string
}

func IsPresentedCLIFlagValueByName(CLIFlag []CLIFlag, name string) bool {
	for _, curFlag := range CLIFlag {
		if curFlag.Name == name {
			return true
		}
	}
	return false
}

func GetCLIFlagValueByName(CLIFlags []CLIFlag, name string) string {
	for _, curFlag := range CLIFlags {
		if curFlag.Name == name {
			return curFlag.Value
		}
	}
	return "" // return empty string if flag not found
}

func ParseAllCLIFlags() []CLIFlag {
	var allFlags []CLIFlag
	re := regexp.MustCompile("-.*=")
	for _, curArgument := range os.Args {
		flagName := re.FindString(curArgument)
		if len(flagName) > 2 {
			flagName = flagName[:len(flagName)-1] // trim last character
			flagName = flagName[1:]               // trim first character
			flagName = strings.TrimSpace(flagName)
			var curFlag CLIFlag
			curFlag.Name = flagName
			allFlags = append(allFlags, curFlag)
		}

	}
	for curFlagIndex, curFlag := range allFlags {
		flag.StringVar(&allFlags[curFlagIndex].Value, curFlag.Name, "", "variant")
	}

	flag.Parse()
	return allFlags
}

//-------------------------------------------------------------------------------
func GetAction() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	} else {
		return "" // return empty string if command not found
	}
}

func GetParams() map[string]string {
	var result map[string]string
	result = make(map[string]string)
	return result
}
