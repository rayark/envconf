package options

import (
	"encoding/json"
	"io"
	"log"

	"gitlab.rayark.com/backend/envconf"
)

// MakeJSONLogger will log the status of environment variables used by given structure.
// The message is in JSON format.
func MakeJSONLogger(w io.Writer) func(map[string]*envconf.EnvStatus) {
	return func(status map[string]*envconf.EnvStatus) {
		logger := log.New(w, "", 0)

		result := map[string]interface{}{}
		result["message"] = "envconf: show environment variables used by configuration and whether they are set"
		result["environment-variables"] = status
		b, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			logger.Printf("envconf encounters an error: %v", err.Error())
			return
		}
		logger.Printf(string(b))
	}
}
