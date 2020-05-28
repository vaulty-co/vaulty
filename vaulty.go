package vaulty

import (
	"fmt"
	"net/http"

	"github.com/vaulty/vaulty/encrypt"
	"github.com/vaulty/vaulty/proxy"
	"github.com/vaulty/vaulty/secrets"
	"github.com/vaulty/vaulty/storage"
	"github.com/vaulty/vaulty/storage/inmem"
	"github.com/vaulty/vaulty/transform/action"
)

func Run(config *Config) error {
	// functionality that goes to router
	st := inmem.NewStorage()

	encrypter, err := encrypt.NewEncrypter(config.EncryptionKey)
	if err != nil {
		return err
	}

	secretStorage := secrets.NewEphemeralStorage(encrypter)

	loaderOptions := &storage.LoaderOptions{
		ActionOptions: &action.Options{
			Encrypter:     encrypter,
			SecretStorage: secretStorage,
		},
		Storage: st,
	}

	err = storage.LoadFromFile(config.RoutesFile, loaderOptions)
	if err != nil {
		return err
	}

	proxy, err := proxy.NewProxy(&proxy.Options{
		ProxyPassword: config.ProxyPassword,
		CAPath:        config.CAPath,
		Storage:       st,
	})
	if err != nil {
		return err
	}

	fmt.Printf("==> Vaulty proxy server started on %v!\n", config.Address)
	return http.ListenAndServe(config.Address, proxy)
}
