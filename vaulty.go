package vaulty

import (
	"fmt"
	"net/http"

	"github.com/vaulty/vaulty/encrypt"
	"github.com/vaulty/vaulty/proxy"
	"github.com/vaulty/vaulty/routing"
	"github.com/vaulty/vaulty/secrets"
	"github.com/vaulty/vaulty/transformer"
	"github.com/vaulty/vaulty/transformer/json"
	"github.com/vaulty/vaulty/transformer/regexp"
)

func Run(config *Config) error {
	encrypter, err := encrypt.NewEncrypter(config.EncryptionKey)
	if err != nil {
		return err
	}

	secretsStorage := secrets.NewEphemeralStorage(encrypter)

	// Create router and load routes from file into router
	loader := routing.NewFileLoader(&routing.FileLoaderOptions{
		Enc:            encrypter,
		SecretsStorage: secretsStorage,
		TransformerFactory: map[string]transformer.Factory{
			"json":   json.Factory,
			"regexp": regexp.Factory,
		},
	})
	routes, err := loader.Load(config.RoutesFile)
	if err != nil {
		return err
	}
	if len(routes) == 0 {
		return fmt.Errorf("No routes were loaded from file: %s", config.RoutesFile)
	}

	router := routing.NewRouter()
	router.SetRoutes(routes)

	proxy, err := proxy.NewProxy(&proxy.Options{
		ProxyPassword: config.ProxyPassword,
		CAPath:        config.CAPath,
		Router:        router,
	})
	if err != nil {
		return err
	}

	fmt.Printf("==> Vaulty proxy server started on %v!\n", config.Address)
	return http.ListenAndServe(config.Address, proxy)
}
