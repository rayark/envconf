package envconf

type Option func(*loader)

// CustomHandleEnvVarsOption creates an option that calls `cb` with a map
// indicating the environment variables checked by envconf.
//
// The map keys are the environment variable names that are checked by envconf,
// and the values present whether the corresponding environment variables are set.
func CustomHandleEnvVarsOption(cb func(map[string]*EnvStatus)) Option {
	return func(l *loader) {
		l.handleEnvironmentVariables = cb
	}
}
