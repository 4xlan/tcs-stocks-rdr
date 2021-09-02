package tcsrdrconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type TCSRDRCfgFile struct {
	Token   string `yaml:"token"`
	Url     string `yaml:"url"`
	Depth   int    `yaml:"depth"`
	Ip      string `yaml:"ip"`
	Port    string `yaml:"port"`
	Timeout int    `yaml:"timeout"`
}

func (cfg *TCSRDRCfgFile)GetConfig(defaultConfigPath *string) error {

	// Calculating the ABS path for config file
	absPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	configFile := fmt.Sprintf("%s/%s", absPath, *defaultConfigPath)

	// Open file and read it
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	// Reading cfg file
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	return nil
}
