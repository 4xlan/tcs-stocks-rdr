package tcsrdrconfig

type TCSRDRConfiguration interface {
	GetConfig(config *TCSRDRCfgFile)
}

type TCSRDRCfgFile struct {
	Token   string `yaml:"token"`
	Url     string `yaml:"url"`
	Depth   int    `yaml:"depth"`
	Ip      string `yaml:"ip"`
	Port    string `yaml:"port"`
	Timeout int    `yaml:"timeout"`
}
