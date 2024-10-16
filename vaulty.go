package vaulty

import (
	"context"
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
	"github.com/vaulty/vaulty/secrets/redisstorage"
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

var storages = map[string]secrets.Factory{
	"memory": memorystorage.Factory,
	"redis":  redisstorage.Factory,
}

func Run(ctx context.Context, conf *config.Config) error {
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

	secretsStorage, err := storages[conf.Storage.Type](&secrets.Config{
		Encrypter:     encrypter,
		StorageConfig: conf.Storage,
	})
	if err != nil {
		return err
	}
	defer secretsStorage.Close()

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

	done := make(chan error, 1)

	server := &http.Server{Addr: conf.Address, Handler: proxy}
	go func() {
		log.Infof("Vaulty proxy server started on %v!\n", conf.Address)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Errorf("Failed to listen and serve: %v", err)
			done <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("Shutting down Vaulty...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Errorf("Failed to clearly shutdown Vaulty: %v", err)
		}
		return nil
	case err := <-done:
		return err
	}
}
