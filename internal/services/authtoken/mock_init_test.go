package authtoken

import "time"

var (
	testSigningSecretKey     = "J@NcRfUjWnZr4u7x!A%D*G-KaPdSgVkY"
	testEncryptionSecretyKey = "%C*F-JaNdRgUkXp2s5v8x/A?D(G+KbPe"
	testExpirationDuration   = time.Hour * 87600 // Ten years

	testAuthTokenServicesConfig = NewAuthTokenConfigBuilder().
					SetJWTSigningSecretKey(testSigningSecretKey).
					SetJWTEncryptionSecretKey(testEncryptionSecretyKey).
					SetJWTExpirationDuration(testExpirationDuration).
					Build()

	//testAuthTokenServices = NewJWTServices(testAuthTokenServicesConfig)
	testAuthTokenServices = &jwtInteractor{
		config: testAuthTokenServicesConfig,
		jws:    &jwsInteractor{},
		jwe:    &jweInteractor{},
	}
)
