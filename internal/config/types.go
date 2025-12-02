package config

type Config struct {
	Version      int                    `toml:"version"`
	Sandboxes    map[string]Sandbox     `toml:"sandboxes"`
	Interceptors map[string]Interceptor `toml:"interceptors"`
	Services     map[string]Service     `toml:"services"`
}

type Sandbox struct {
	Source  string         `toml:"source"`
	Version string         `toml:"version"`
	Config  map[string]any `toml:"config"`
}

type Interceptor struct {
	Source   string         `toml:"source"`
	Version  string         `toml:"version"`
	Priority int            `toml:"priority"`
	Config   map[string]any `toml:"config"`
}

type Service struct {
	Version string         `toml:"version"`
	Config  map[string]any `toml:"config"`
}
