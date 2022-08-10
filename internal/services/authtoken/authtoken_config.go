package authtoken

import "time"

const (
	tokenExpirationDefault = 30 * time.Minute
	issuerKeyValue         = "access-control-api"
	permissionsKeyName     = "permissions"
)

type AuthTokenConfigBuilder interface {
	SetJWTSigningSecretKey(string) AuthTokenConfigBuilder
	SetJWTEncryptionSecretKey(string) AuthTokenConfigBuilder
	SetJWTExpirationDuration(time.Duration) AuthTokenConfigBuilder
	Build() *AuthTokenConfig
}

type authTokenConfigBuilder struct {
	jwtSigningSecretKey     string
	jwtEncryptionSecretKey  string
	tokenExpirationDuration time.Duration
}

func NewAuthTokenConfigBuilder() AuthTokenConfigBuilder {
	return &authTokenConfigBuilder{
		tokenExpirationDuration: tokenExpirationDefault,
	}
}

func (b *authTokenConfigBuilder) Build() *AuthTokenConfig {
	return &AuthTokenConfig{
		builder: b,
	}
}

func (b *authTokenConfigBuilder) SetJWTSigningSecretKey(key string) AuthTokenConfigBuilder {
	b.jwtSigningSecretKey = key
	return b
}

func (b *authTokenConfigBuilder) SetJWTEncryptionSecretKey(key string) AuthTokenConfigBuilder {
	b.jwtEncryptionSecretKey = key
	return b
}

func (b *authTokenConfigBuilder) SetJWTExpirationDuration(duration time.Duration) AuthTokenConfigBuilder {
	if duration > 0 {
		b.tokenExpirationDuration = duration
	}
	return b
}

type AuthTokenConfig struct {
	builder *authTokenConfigBuilder
}

func (c *AuthTokenConfig) SigningSecretKey() string {
	return c.builder.jwtSigningSecretKey
}

func (c *AuthTokenConfig) EncryptionSecretKey() string {
	return c.builder.jwtEncryptionSecretKey
}

func (c *AuthTokenConfig) ExpirationDuration() time.Duration {
	return c.builder.tokenExpirationDuration
}
