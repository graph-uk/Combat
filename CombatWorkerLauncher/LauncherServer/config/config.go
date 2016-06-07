package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Port                  int
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	Region                string
	AMIId                 string
	ShutdownBehavior      string
	InstanceType          string
	SecurityGroupId       string
	KeyName               string
	TagName               string
	LauncherSecretKey     string
}

func defaultConfig() string {
	return `{
	"Port":9091, 
	"AWS_ACCESS_KEY_ID":"", 
	"AWS_SECRET_ACCESS_KEY":"",
	"Region":"eu-central-1", 
	"AMIId":"", 
	"ShutdownBehavior":"terminate", 
	"InstanceType":"t2.nano", 
	"SecurityGroupId":"", 
	"KeyName":"", 
	"TagName":"CombatWorker",
	"LauncherSecretKey":""
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

	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		fmt.Println("Cannot unmarshal config.json. Check format, or delete config.json. Default config will be created at next run")
		fmt.Println(err.Error())
		return &conf, err
	}

	return &conf, nil
}

func makeDefaultConfig() error {
	return ioutil.WriteFile("config.json", []byte(defaultConfig()), 0777)
}
