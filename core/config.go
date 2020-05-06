package core

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Environment   string `yaml:"environment" envconfig:"PROXY_ENV"`
	BaseHost      string `yaml:"base_host" envconfig:"BASE_HOST"`
	ProxyPassword string `yaml:"proxy_pass" envconfig:"PROXY_PASS"`
	CaPath        string `yaml:"ca_path" envconfig:"CA_PATH"`
	Redis         struct {
		URL string `yaml:"url" envconfig:"REDIS_URL"`
	}
}

var Config *Configuration

func LoadConfig(file string) *Configuration {
	Config = &Configuration{}

	readFile(file, Config)
	readEnv(Config)
	setDefaults(Config)

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

func setDefaults(cfg *Configuration) {
	if cfg.CaPath == "" {
		cfg.CaPath = "~/.vaulty"
	}

	var err error
	cfg.CaPath, err = homedir.Expand(cfg.CaPath)

	if err != nil {
		panic(err)
	}
}
