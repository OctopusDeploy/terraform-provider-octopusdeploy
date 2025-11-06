package internal

import (
	"flag"
	"os"
	"testing"
)

var createSharedContainer = flag.Bool("createSharedContainer", false, "Set to true to run integration tests in containers")

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
