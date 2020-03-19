package proxy

import (
	"fmt"
	"os"
	"testing"

	"github.com/vaulty/proxy/core"
)

func TestMain(m *testing.M) {
	core.Config().Environment = "test"
	core.Config().Host = "proxy.test"

	exitCode := m.Run()

	fmt.Println("after all tests")
	os.Exit(exitCode)

}
