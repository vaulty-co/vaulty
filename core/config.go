package core

import (
	"io"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Environment       string `yaml:"environment" envconfig:"PROXY_ENV"`
	BaseHost          string `yaml:"base_host" envconfig:"BASE_HOST"`
	ProxyPassword     string `yaml:"proxy_pass" envconfig:"PROXY_PASS"`
	RoutesFile        string `default:"~/.vaulty/routes.json" yaml:"routes_file" envconfig:"ROUTES_FILE"`
	CaPath            string `default:"~/.vaulty" yaml:"ca_path" envconfig:"CA_PATH"`
	IsSingleVaultMode bool
	Redis             struct {
		URL string `yaml:"url" envconfig:"REDIS_URL"`
	}
}

var Config *Configuration

func LoadConfig(file string) *Configuration {
	Config = &Configuration{
		// for the MVP this will be the only working option even
		// if we have code for multu vault mode. In the future
		// we will see if there is a need to have multiple vaults
		IsSingleVaultMode: true,
	}

	readFile(file, Config)
	readEnv(Config)
	expandPaths(Config)

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
	if err != nil && err != io.EOF {
		panic(err)
	}
}

func readEnv(cfg *Configuration) {
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err)
	}
}

func expandPaths(cfg *Configuration) {
	var err error

	cfg.CaPath, err = homedir.Expand(cfg.CaPath)
	if err != nil {
		panic(err)
	}

	cfg.RoutesFile, err = homedir.Expand(cfg.RoutesFile)
	if err != nil {
		panic(err)
	}
}
