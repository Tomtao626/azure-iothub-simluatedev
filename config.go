package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

const serviceName = "simluatedev"

var BuildVersion = "default"

var Config struct {
	DeviceProvisioning struct {
		PolicyName       string
		Method           string
		IdScope          string
		OrderId          string
		DevEndPoint      string
		RegBaseUrl       string
		RegStatusBaseUrl string
	}

	CommonKeys struct {
		PrimaryKey     string
		SecondaryKey   string
		RegistrationId string
	}

	DataUploadConf struct {
		IdentityId      string
		Pid             string
		MsgType         string
		FirmwareVersion string
		TimeVal         time.Duration
	}

	DirectMethodConf struct {
		MethodName        string
		SuccessStatusCode int
		MethodJsonPath    string
	}
}

func LoadCfg() {
	showVersion := flag.Bool("v", false, "show version")
	configFilePath := flag.String("c", "", "configure file path")
	configFileTest := flag.Bool("t", false, "configure file test")
	flag.Parse()
	if *showVersion {
		fmt.Println(BuildVersion)
		os.Exit(0)
	}
	Info("cfg path :", *configFilePath)
	if *configFilePath == "" {
		*configFilePath = "config.toml"
	}

	var err error
	if _, err = toml.DecodeFile(*configFilePath, &Config); err != nil {
		Fatal("toml fail to parse file :", err)
		os.Exit(-1)
	}

	Infof("%+v", Config)

	if *configFileTest {
		fmt.Printf("configuration file %s test is successful\n", *configFilePath)
		os.Exit(0)
	}
}
