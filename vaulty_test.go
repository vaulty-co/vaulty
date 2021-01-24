package vaulty

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vaulty/vaulty/config"
)

func TestCreateProxyWithConfig(t *testing.T) {
	// create temporary directory to generate CA files
	tmpCAdir, err := ioutil.TempDir("", "vaulty-ca")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(tmpCAdir) // clean up

	conf := &config.Config{
		CAPath:     tmpCAdir,
		RoutesFile: "./routing/testdata/routes.json",
	}

	err = conf.FromEnvironment()
	require.NoError(t, err)

	err = conf.GenerateMissedValues()
	require.NoError(t, err)

	done := make(chan struct{}, 1)

	go func() {
		err := Run(conf)
		require.NoError(t, err)
	}()

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM)
		done <- struct{}{}
		<-sigs
	}()

	<-done

	time.Sleep(1 * time.Second / 2)

	pid := syscall.Getpid()
	p, err := os.FindProcess(pid)
	require.NoErrorf(t, err, "Failed to find current process: %v", err)

	err = p.Signal(syscall.SIGTERM)
	require.NoErrorf(t, err, "Failed to signal process: %v", err)
}
