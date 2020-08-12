package options_test

import (
	"bytes"
	"os"
	"testing"

	"gitlab.rayark.com/backend/envconf"
	"gitlab.rayark.com/backend/envconf/options"
)

func TestLogger(t *testing.T) {
	os.Setenv("TEST_INTEGER", "-3")
	os.Setenv("TEST_UNSIGNED_INTEGER", "3")

	config := struct {
		String  string `env:"string"`
		Integer int    `env:"integer"`
	}{}
	buf := bytes.NewBuffer(nil)
	envconf.Load("TEST", &config, envconf.CustomHandleEnvVarsOption(options.LogStatusOfEnvVars(buf)))

	if buf.String() == "" {
		t.Errorf("failed to log environment variable status")
	}
}
