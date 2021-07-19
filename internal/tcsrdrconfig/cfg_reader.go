package tcsrdrconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func GetConfig(config *TCSRDRCfgFile, defaultConfigPath *string) error {

	// Calculating the ABS path for config file
	absp, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	filepath := fmt.Sprintf("%s/%s", absp, *defaultConfigPath)

	// Open file and read it
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Reading cfg file
	yaml.Unmarshal([]byte(data), config)
	return nil
}
