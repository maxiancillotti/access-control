package config

type ServiceConfig struct {
	JWTSigningSecretKey string `toml:"jwt_signing_secret_key" env:"AUTH_JWT_SIGNING_SECRET_KEY" env-required`

	JWTEncryptionSecretKey string `toml:"jwt_encryption_secret_key" env:"AUTH_JWT_ENCRYPTION_SECRET_KEY" env-required`

	JWTExpirationDuration StrTimeDuration `toml:"jwt_expiration_time_duration" env:"AUTH_JWT_EXPIRATION_TIME_DURATION"`
}
