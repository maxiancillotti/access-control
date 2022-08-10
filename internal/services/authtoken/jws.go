package authtoken

import (
	"encoding/base64"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

type jwsServices interface {
	signToken(token jwt.Token, secretKey []byte) ([]byte, error)
	verifyToken(token, secretKey []byte) (jwt.Token, error)
}

type jwsInteractor struct{}

func (s *jwsInteractor) signToken(token jwt.Token, secretKey []byte) ([]byte, error) {
	b64SignSecretKeyStr := base64.URLEncoding.EncodeToString(secretKey)
	signedToken, err := jwt.Sign(token, jwa.HS256, []byte(b64SignSecretKeyStr))
	if err != nil {
		return nil, err
	}
	return signedToken, nil
}

func (s *jwsInteractor) verifyToken(token, secretKey []byte) (jwt.Token, error) {
	b64SignSecretKeyStr := base64.URLEncoding.EncodeToString(secretKey)

	jwtToken, err := jwt.Parse(token,
		jwt.WithVerify(jwa.HS256, []byte(b64SignSecretKeyStr)),
		jwt.WithValidate(true),
		jwt.WithIssuer(issuerKeyValue))
	if err != nil {
		return nil, err
	}
	return jwtToken, nil
}
