package vaulty

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/vaulty/vaulty/config"
	"github.com/vaulty/vaulty/encryption"
	"github.com/vaulty/vaulty/encryption/aesgcm"
	"github.com/vaulty/vaulty/encryption/awskms"
	"github.com/vaulty/vaulty/encryption/noneenc"
	"github.com/vaulty/vaulty/proxy"
	"github.com/vaulty/vaulty/routing"
	"github.com/vaulty/vaulty/secrets"
	"github.com/vaulty/vaulty/secrets/memorystorage"
	"github.com/vaulty/vaulty/transformer"
	"github.com/vaulty/vaulty/transformer/form"
	"github.com/vaulty/vaulty/transformer/json"
	"github.com/vaulty/vaulty/transformer/regexp"
)

var encrypters = map[string]encryption.Factory{
	"awskms": awskms.Factory,
	"aesgcm": aesgcm.Factory,
	"none":   noneenc.Factory,
}

var transformers = map[string]transformer.Factory{
	"json":   json.Factory,
	"regexp": regexp.Factory,
	"form":   form.Factory,
}

func Run(conf *config.Config) error {
	if conf.Debug {
		log.SetFormatter(&log.TextFormatter{
			ForceColors: true,
		})
		log.SetLevel(log.DebugLevel)
		fmt.Println("Warning! Body of requests and responses will be exposed in logs!")
	}

	encrypter, err := encrypters[conf.Encryption.Type](conf)
	if err != nil {
		return err
	}

	secretsStorage, err := memorystorage.Factory(&secrets.Config{
		Encrypter:      encrypter,
		StorageConfing: conf.Storage,
	})
	if err != nil {
		return err
	}

	// Create router and load routes from file into router
	loader := routing.NewFileLoader(&routing.FileLoaderOptions{
		Enc:                encrypter,
		SecretsStorage:     secretsStorage,
		Salt:               conf.Salt,
		TransformerFactory: transformers,
	})
	routes, err := loader.Load(conf.RoutesFile)
	if err != nil {
		return err
	}
	if len(routes) == 0 {
		return fmt.Errorf("No routes were loaded from file: %s", conf.RoutesFile)
	}

	router := routing.NewRouter()
	router.SetRoutes(routes)

	proxy, err := proxy.NewProxy(&proxy.Options{
		ProxyPassword: conf.ProxyPassword,
		CAPath:        conf.CAPath,
		Router:        router,
	})
	if err != nil {
		return err
	}

	fmt.Printf("==> Vaulty proxy server started on %v!\n", conf.Address)
	return http.ListenAndServe(conf.Address, proxy)
}
