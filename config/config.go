package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
)

// Global config object singleton
var gloablConfigSingleton *GlobalConfig

func InitGlobalConfig(conf *GlobalConfig) {
	gloablConfigSingleton = conf
}

// GetProvisionerConfig returns a copy of global config object
func GetGlobalConfig() GlobalConfig {
	return *gloablConfigSingleton
}

func PrintConfigTemplate(format string) {
	if len(format) < 1 {
		format = "yaml"
	}

	printFile := func(content []byte, format string) {
		f, err := os.OpenFile("provision-conf.template."+format, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
		f.Write(content)
		f.Sync()
	}

	confTemplate := gloablConfigSingleton

	if confTemplate == nil {
		confTemplate = &GlobalConfig{
			LogConfig: LoggerConfig{
				Level:     "debug",
				LogFormat: "text",
			},
		}
	}

	var conf []byte
	switch format {
	case "json":
		conf, _ = json.MarshalIndent(&confTemplate, "", "  ")
	case "yaml":
		conf, _ = yaml.Marshal(&confTemplate)
	case "toml":
		conf, _ = toml.Marshal(&confTemplate)
	}
	printFile(conf, format)
}
