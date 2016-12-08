package config

import (
//Standard library packages
	"path/filepath"
	"io/ioutil"
	"log"
	"flag"
//Third party packages
	"gopkg.in/yaml.v2"
)

var configWrapper *ConfigWrapper

func InitConfig(env string, configPath string) {
	if configWrapper == nil {
		filename, _ := filepath.Abs(configPath)
		yamlFile, err := ioutil.ReadFile(filename)
		var config Config
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		configWrapper = &ConfigWrapper{}
		switch env {
		case "development":
			configWrapper.currentConfig = config.Development
		case "test":
			configWrapper.currentConfig = config.Test
		case "production":
			configWrapper.currentConfig = config.Production
		}
	}
}

func GetConfigWrapper() *ConfigWrapper {
	if configWrapper == nil {
		configPath := flag.String("config", "empty", "configuration file (.yml) full path")
		env := flag.String("env", "wrong", "environment: development | test | production")
		flag.Parse()
		filename, _ := filepath.Abs(*configPath)
		log.Println("trying to read file ", filename)
		yamlFile, err := ioutil.ReadFile(filename)
		var config Config

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		configWrapper = &ConfigWrapper{}
		switch *env {
		case "development":
			configWrapper.currentConfig = config.Development
		case "test":
			configWrapper.currentConfig = config.Test
		case "production":
			configWrapper.currentConfig = config.Production
		}


	}
	return configWrapper
}

type ConfigWrapper struct {
	Config        Config
	currentConfig EnvConfig
}

type Config struct {
	Development EnvConfig `yaml:"development"`
	Test        EnvConfig `yaml:"test"`
	Production  EnvConfig `yaml:"production"`

}

type EnvConfig struct {
	MongoHost                    string `yaml:"mongo_host"`
	Env                          string `yaml:"env"`
	EmailServerAddress           string `yaml:"email_server_address"`
	EmailServerPort              int    `yaml:"email_server_port"`
	EmailServerUsername          string `yaml:"email_server_username"`
	EmailServerPassword          string `yaml:"email_server_password"`
	EmailServerFrom              string `yaml:"email_server_from"`
	EmailServerBcc               string `yaml:"email_server_bcc"`
	LogPath                      string `yaml:"log_path"`
	AdminAuth                    string `yaml:"admin_auth"`
}

func (configWrapper *ConfigWrapper) GetCurrent() *EnvConfig {
	return &configWrapper.currentConfig
}
