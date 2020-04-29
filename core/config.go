package core

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Environment   string `yaml:"environment" envconfig:"PROXY_ENV"`
	BaseHost      string `yaml:"base_host" envconfig:"BASE_HOST"`
	ProxyPassword string `yaml:"proxy_pass" envconfig:"PROXY_PASS"`
	Redis         struct {
		URL string `yaml:"url" envconfig:"REDIS_URL"`
	}
}

var Config *Configuration

func LoadConfig(file string) *Configuration {
	Config = &Configuration{}

	readFile(file, Config)
	readEnv(Config)

	return Config
}

func readFile(file string, cfg *Configuration) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		panic(err)
	}
}

func readEnv(cfg *Configuration) {
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err)
	}
}
