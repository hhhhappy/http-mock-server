package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

const (
	configFilePath  = "/config/config.yml"
	requestFilePath = "/config/request.yml"
)

type Config struct {
	Port              string            `yaml:"port"`
	PrintBufferLength int               `yaml:"printBufferLength"`
	LogPath           string            `yaml:"logPath"`
	LogAccessSummary  bool              `yaml:"logAccessSummary"`
	DefaultHeaders    map[string]string `yaml:"defaultHeaders"`
	Requests          []Request         `yaml:"urlList"`
}

type Request struct {
	Url            string            `yaml:"url"`
	Type           string            `yaml:"type"`
	ReturnBodyFile string            `yaml:"returnBodyFile"`
	Code           int               `yaml:"code"`
	Header         map[string]string `yaml:"header"`
}

var conf *Config

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	//load config file
	configFileByte, err := ioutil.ReadFile(path.Join(pwd, configFilePath))
	if err != nil {
		panic(err)
	}

	conf = &Config{}
	err = yaml.Unmarshal(configFileByte, conf)
	if err != nil {
		panic(err)
	}

	//load request file
	requestFileByte, err := ioutil.ReadFile(path.Join(pwd, requestFilePath))
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(requestFileByte, &conf.Requests)
	if err != nil {
		panic(err)
	}
}

func GetConf() Config {
	return *conf
}

func (conf Config) GetRequestDefinition(url string) (obj *Request) {
	for key, unit := range conf.Requests {
		if unit.Url == url {
			return &conf.Requests[key]
		}
	}

	return nil
}
