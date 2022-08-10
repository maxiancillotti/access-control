package config

type HttpServerConfig struct {
	Hostname          string          `toml:"hostname" env:"HTTPSERVER_HOSTNAME" env-required`
	HostPort          int             `toml:"host_port" env:"HTTPSERVER_HOST_PORT" env-required`
	ReadHeaderTimeout StrTimeDuration `toml:"read_header_timeout_duration" env:"HTTPSERVER_READ_HEADER_TIMEOUT_DURATION" env-required`
	WriteTimeout      StrTimeDuration `toml:"write_timeout_duration" env:"HTTPSERVER_WRITE_TIMEOUT_DURATION" env-required`
}

/*
type HttpClientConfig struct {
	Timeout StrTimeDuration `toml:"timeout_duration" env:"HTTP_CLIENT_TIMEOUT_DURATION" env`
}
*/
