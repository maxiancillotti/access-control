package config

type DatabaseConfig struct {
	Hostname string `toml:"host" env:"DB_HOSTNAME" env-required`
	Port     int    `toml:"port" env:"DB_PORT" env-required`
	Instance string `toml:"instance" env:"DB_INSTANCE"`
	DBName   string `toml:"dbname" env:"DB_NAME" env-required`
	User     string `toml:"user" env:"DB_USER" env-required`
	Password string `toml:"password" env:"DB_PASSWORD" env-required`
}
