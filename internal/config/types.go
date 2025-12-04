package config

type Config struct {
	Version      int                    `toml:"version"`
	Sandboxes    map[string]Sandbox     `toml:"sandboxes"`
	Interceptors map[string]Interceptor `toml:"interceptors"`
	Services     map[string]Service     `toml:"services"`
}

type Sandbox struct {
	Enable         *bool           `toml:"enable"`
	AutoReconnect  *bool           `toml:"auto_reconnect"`
	ConnectTimeout *string         `toml:"connect_timeout"`
	Source         string          `toml:"source"`
	Hash           string          `toml:"hash"`
	Config         *map[string]any `toml:"config"`
}

type Interceptor struct {
	Enable         *bool           `toml:"enable"`
	AutoReconnect  *bool           `toml:"auto_reconnect"`
	ConnectTimeout *string         `toml:"connect_timeout"`
	Source         string          `toml:"source"`
	Priority       int             `toml:"priority"`
	Hash           string          `toml:"hash"`
	Config         *map[string]any `toml:"config"`
}

type Service struct {
	Enable         *bool           `toml:"enable"`
	AutoReconnect  *bool           `toml:"auto_reconnect"`
	ConnectTimeout *string         `toml:"connect_timeout"`
	Source         string          `toml:"source"`
	Hash           string          `toml:"hash"`
	Config         *map[string]any `toml:"config"`
}
