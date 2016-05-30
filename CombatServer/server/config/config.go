package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Port                 int
	CountOfSavedSessions int
}

func PrintConfigExample() {
	var conf Config
	conf.Port = 9090
	conf.CountOfSavedSessions = 10
	confBytes, err := json.Marshal(conf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(confBytes))
}

func defaultConfig() string {
	return `{
	"port":9090, 
	"countOfSavedSessionsFiles":10
}`
}

//Try to load config - if not found - create and load again.
//If cannot create or load - print error, help text and exit(1)
func LoadConfig() (*Config, error) {
	var conf Config

	// create default config.json if not exist
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		fmt.Println("config.json is not found. Default config will be created")
		if makeDefaultConfig() != nil {
			fmt.Println("Cannot create default config.json")
			fmt.Println(err.Error())
			return &conf, err
		}
	}

	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Cannot read config.json")
		fmt.Println(err.Error())
		return &conf, err
	}

	if json.Unmarshal(bytes, &conf) != nil {
		fmt.Println("Cannot unmarshal config.json. Check format, or delete config.json. Default config will be created at next run")
		fmt.Println(err.Error())
		return &conf, err
	}

	return &conf, nil
}

func makeDefaultConfig() error {
	return ioutil.WriteFile("config.json", []byte(defaultConfig()), 0777)
}
