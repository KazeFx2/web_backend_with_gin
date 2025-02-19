package Config

import (
	"gopkg.in/yaml.v3"
	"main/Logger"
	"os"
)

func LoadYaml(yamlPath string, container interface{}) error {
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		Logger.LogE("can not read yaml config file '%s': %v", yamlPath, err)
		return err
	}

	if err := yaml.Unmarshal(data, container); err != nil {
		Logger.LogE("parse YAML failed: %v", err)
		return err
	}
	return nil
}

func WriteYaml(yamlPath string, info interface{}) error {
	yamlData, err := yaml.Marshal(info)
	if err != nil {
		Logger.LogE("transfer config to yaml failed: %v", err)
		return err
	}
	err = os.WriteFile(yamlPath, yamlData, 0644)
	if err != nil {
		Logger.LogE("write to file '%s' failed: %v", yamlPath, err)
		return err
	}
	return nil
}

/*
Config Example

type Item struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Items []Item `yaml:"items"`
}

Yaml Example

server:
  port: 8080
  host: localhost
database:
  user: MyUser
  password: MyPassword
  name: MyDatabase
items:
  - key: item1
    value: value1
  - key: item2
    value: value2
  - key: item3
    value: value3
*/
