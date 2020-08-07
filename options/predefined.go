package options

import (
	"encoding/json"
	"io"
	"log"
)

// LogStatusOfEnvVars will log the status of environment variables used by given structure.
// The message is in JSON format.
func LogStatusOfEnvVars(w io.Writer) func(map[string]bool) {
	return func(status map[string]bool) {
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
