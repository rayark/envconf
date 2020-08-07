package envconf

type Option func(*loader)

func CustomHandleEnvVarsOption(cb func(map[string]bool)) Option {
	return func(l *loader) {
		l.handleEnvironmentVariables = cb
	}
}
