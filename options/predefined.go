package options

import (
	"encoding/json"
	"io"
	"log"

	"gitlab.rayark.com/backend/envconf"
)

func LogStatusOfEnvironmentVariables(w io.Writer) envconf.Option {
	return envconf.CustomHandleEnvironmentVariablesOption(func(status map[string]bool) {
		result := map[string]interface{}{}
		result["message"] = "envconf: show environment variables used by configuration and whether they are set"
		result["environment-variables"] = status
		b, _ := json.MarshalIndent(result, "", "    ")
		log.New(w, "", 0).Printf(string(b))
	})
}
