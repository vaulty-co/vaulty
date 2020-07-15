package config

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"github.com/vaulty/vaulty/ca"
)

type (
	Config struct {
		// Network address that Vaulty should listen on
		Address string

		// File containing route definitions
		RoutesFile string

		// Path to CA files
		CAPath string

		// Password for the forward proxy
		ProxyPassword string

		// Debug mode, exposes bodies of request and response
		Debug bool

		// Salt for hash action
		Salt string

		Encryption *Encryption

		Storage *Storage
	}

	Encryption struct {
		// Encryption Type
		Type string `envconfig:"VAULTY_ENCRYPTION_TYPE" default:"none"`

		// Encryption key that should be used for AES GCM encryption
		Key string `envconfig:"VAULTY_ENCRYPTION_KEY"`

		// Region of AWS KMS Master Key
		AWSKMSRegion string `envconfig:"VAULTY_ENCRYPTION_AWS_KMS_REGION"`

		// ID of AWS KMS Master Key (optional if KMSKeyAlias is set)
		AWSKMSKeyID string `envconfig:"VAULTY_ENCRYPTION_AWS_KMS_KEY_ID"`

		// Alias of AWS KMS Master Key (optional if KMSKeyID is set)
		AWSKMSKeyAlias string `envconfig:"VAULTY_ENCRYPTION_AWS_KMS_KEY_ALIAS"`
	}

	Storage struct {
		// Storage Type
		Type string `envconfig:"VAULTY_STORAGE_TYPE" default:"memory"`
		// ...
	}
)

func NewConfig() *Config {
	return &Config{Encryption: &Encryption{}}
}

func (c *Config) FromEnvironment() error {
	return envconfig.Process("vaulty", c)
}

// GenerateMissedValues generates proxy password if it's not provided
// and CA certificate and key
func (c *Config) GenerateMissedValues() error {
	var err error

	if c.ProxyPassword == "" {
		pass := make([]byte, 16)
		_, err = io.ReadFull(rand.Reader, pass)
		if err != nil {
			return err
		}
		c.ProxyPassword = fmt.Sprintf("%x", pass)
		fmt.Printf("No password for forward proxy provided (PROXY_PASS)!\nRandom password is used: %s\n", c.ProxyPassword)
	}

	caCertFile := filepath.Join(c.CAPath, "ca.cert")
	caKeyFile := filepath.Join(c.CAPath, "ca.key")
	if isFileMissed(caCertFile) || isFileMissed(caKeyFile) {
		fmt.Printf("No CA certificate / key found (in CA_PATH).\nGenerate CA cert: %s\nCA private key: %s\n",
			caCertFile, caKeyFile)

		rootCertPEM, rootKeyPEM := ca.GenCA()
		ioutil.WriteFile(caCertFile, rootCertPEM, 0644)
		ioutil.WriteFile(caKeyFile, rootKeyPEM, 0644)
	}

	return nil
}

func isFileMissed(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return true
	}

	return false
}
