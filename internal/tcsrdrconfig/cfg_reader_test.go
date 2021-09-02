package tcsrdrconfig

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfig(t *testing.T) {
	defaultConfigPath := "config.yaml.example"
	req := require.New(t)

	testList := map[string]struct {
		configFile string
		configPath string
		state int
	} {
		"Correct": {
			configFile: "IyMjIENvbm5lY3QgdG8gbWFya2V0IHBhcmFtcwojIFVybCBmb3IgcG9ydGZvbGlvCnVybDogImh0dHBzOi8vYXBpLWludmVzdC50aW5rb2ZmLnJ1L29wZW5hcGkiCiMgVG9rZW4gZm9yIEFQSQp0b2tlbjogIjFUb2tlbjEiCiMgVGltZW91dCBpbiBzZWMuCnRpbWVvdXQ6IDIwCiMgRGVwdGggb2Ygb3JkZXJib29rCmRlcHRoOiAzMAoKIyMjIFNlcnZlciBwYXJhbXMKIyBMaXN0ZW4gb24gdGhlIG5leHQgYWRkcmVzcwppcDogIjE2OS4yNTQuMC4xIgojIExpc3RlbiBvbiB0aGUgbmV4dCBwb3J0CnBvcnQ6ICI5MTAwIgo=",
			configPath: defaultConfigPath,
			state: 0, // All ok
		},
		"LostField": {
			configFile: "IyMjIENvbm5lY3QgdG8gbWFya2V0IHBhcmFtcwojIFVybCBmb3IgcG9ydGZvbGlvCnVybDogImh0dHBzOi8vYXBpLWludmVzdC50aW5rb2ZmLnJ1L29wZW5hcGkiCiMgVG9rZW4gZm9yIEFQSQp0b2tlbjogIiIKIyBUaW1lb3V0IGluIHNlYy4KdGltZW91dDogNQojIERlcHRoIG9mIG9yZGVyYm9vawoKIyMjIFNlcnZlciBwYXJhbXMKIyBMaXN0ZW4gb24gdGhlIG5leHQgYWRkcmVzcwppcDogIiIgCiMgTGlzdGVuIG9uIHRoZSBuZXh0IHBvcnQKcG9ydDogIjkxMDAiCg==",
			configPath: defaultConfigPath,
			state: 1, // One field lost
		},
		"IncorrectPath": {
			configFile: "ZW1wdHlTdHJpbmcK",
			configPath: "config.yaml.example10",
			state: 2, // Incorrect path
		},
	}

	exampleConfig := TCSRDRCfgFile{
		Token: "1Token1",
		Url: "https://api-invest.tinkoff.ru/openapi",
		Timeout: 20,
		Depth: 30,
		Ip: "169.254.0.1",
		Port: "9100",
	}

	for name, testCase := range testList {
		t.Run(name, func(t *testing.T) {
			cfg, err := configInitAndGet(&testCase.configPath, &defaultConfigPath, &testCase.configFile)
			switch testCase.state {
			case 0:
				req.NoError(err)
				req.Equal(exampleConfig, *cfg)
			case 1:
				req.NoError(err)
				req.NotEqual(exampleConfig, *cfg)
			case 2:
				req.Error(err)
				req.NotEqual(exampleConfig, *cfg)
			}

		})
	}
}

func configInitAndGet(configPath *string, defaultConfigPath *string, configFileContent *string) (*TCSRDRCfgFile, error) {

	tFile := &TCSRDRCfgFile{}

	err := writeConfigExample(configPath, configFileContent)
	if err != nil {
		return tFile, err
	}

	err = tFile.GetConfig(defaultConfigPath)
	log.Println(tFile)
	if err != nil {
		return tFile, err
	}

	return tFile, nil
}

func writeConfigExample(configPath *string, b64File *string) error{

	stringDecoded, err := base64.StdEncoding.DecodeString(*b64File)
	if err != nil {
		return err
	}

	curExec, err := os.Executable()
	if err != nil {
		return err
	}

	curFilePath := filepath.Dir(curExec)
	err = os.WriteFile(fmt.Sprintf("%s/%s", curFilePath, *configPath), stringDecoded, 0644)

	if err != nil {
		return err
	}

	return nil
}