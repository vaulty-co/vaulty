package vaulty

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/config"
)

func TestCreateProxyWithConfig(t *testing.T) {
	// create temporary directory to generate CA files
	cadir := "testdata"

	conf := &config.Config{
		CAPath:     cadir,
		RoutesFile: "./routing/testdata/routes.json",
	}

	err := conf.FromEnvironment()
	require.NoError(t, err)

	err = conf.GenerateMissedValues()
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		err := Run(ctx, conf)
		require.NoError(t, err)
		close(done)
	}()

	// wait for the server to start
	time.Sleep(1 * time.Second)

	// send signal to stop the server
	cancel()

	// wait for the server to stop
	<-done
}
