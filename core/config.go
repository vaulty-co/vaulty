package core

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/go-homedir"
	"github.com/vaulty/proxy/ca"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Environment       string `yaml:"environment" envconfig:"PROXY_ENV"`
	BaseHost          string `yaml:"base_host" envconfig:"BASE_HOST"`
	ProxyPassword     string `yaml:"proxy_pass" envconfig:"PROXY_PASS"`
	RoutesFile        string `default:"~/.vaulty/routes.json" yaml:"routes_file" envconfig:"ROUTES_FILE"`
	CaPath            string `default:"~/.vaulty" yaml:"ca_path" envconfig:"CA_PATH"`
	EncryptionKey     string `yaml:"encryption_key" envconfig:"ENCRYPTION_KEY"`
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
	setDefaults(Config)

	return Config
}

func readFile(file string, cfg *Configuration) {
	if _, err := os.Stat(file); err != nil {
		fmt.Println("No configuration file found. Read config from ENV")
		return
	}

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

func setDefaults(cfg *Configuration) {
	var err error

	if cfg.ProxyPassword == "" {
		pass := make([]byte, 16)
		_, err = io.ReadFull(rand.Reader, pass)
		if err != nil {
			panic(err)
		}
		cfg.ProxyPassword = fmt.Sprintf("%x", pass)
		fmt.Printf("No password for forward proxy provided (PROXY_PASS)!\nRandom password is used: %s\n", cfg.ProxyPassword)
	}

	if isFileMissed(filepath.Join(cfg.CaPath, "ca.pem")) || isFileMissed(filepath.Join(cfg.CaPath, "ca.key")) {
		fmt.Printf("No CA certificate / key found (in CA_PATH).\nGenerate CA cert: %s\nCA private key: %s\n",
			cfg.CaPath+"/ca.pem", cfg.CaPath+"/ca.key")

		rootCertPEM, rootKeyPEM := ca.GenCA()
		ioutil.WriteFile(cfg.CaPath+"/ca.pem", rootCertPEM, 0644)
		ioutil.WriteFile(cfg.CaPath+"/ca.key", rootKeyPEM, 0644)
	}
}

func isFileMissed(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return true
	}

	return false
}
