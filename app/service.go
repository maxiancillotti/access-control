package app

import (
	"github.com/maxiancillotti/access-control/app/config"
	"github.com/maxiancillotti/access-control/internal/services/authtoken"
)

func GetServiceConfig(config *config.ServiceConfig) *authtoken.AuthTokenConfig {

	return authtoken.NewAuthTokenConfigBuilder().
		SetJWTSigningSecretKey(config.JWTSigningSecretKey).
		SetJWTEncryptionSecretKey(config.JWTEncryptionSecretKey).
		SetJWTExpirationDuration(config.JWTExpirationDuration.GetDuration()).
		Build()
}
