package authtoken

import (
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwe"
)

type jweServices interface {
	encryptToken(token, secretKey []byte) ([]byte, error)
	decryptToken(token, secretKey []byte) ([]byte, error)
}

type jweInteractor struct{}

func (s *jweInteractor) encryptToken(token, secretKey []byte) ([]byte, error) {
	encryptedToken, err := jwe.Encrypt(token, jwa.DIRECT, secretKey, jwa.A128CBC_HS256, jwa.NoCompress)
	if err != nil {
		return nil, err
	}
	return encryptedToken, nil
}

func (s *jweInteractor) decryptToken(token, secretKey []byte) ([]byte, error) {
	decryptedToken, err := jwe.Decrypt(token, jwa.DIRECT, secretKey)
	if err != nil {
		return nil, err
	}
	return decryptedToken, nil
}
