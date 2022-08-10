package authtoken

import (
	"strconv"
	"time"

	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/lestrrat-go/jwx/jwt"

	"github.com/pkg/errors"
)

func NewJWTServices(config *AuthTokenConfig) internal.AuthTokenServices {
	return &jwtInteractor{
		config: config,
		jws:    &jwsInteractor{},
		jwe:    &jweInteractor{},
	}
}

type jwtInteractor struct {
	config *AuthTokenConfig
	jws    jwsServices
	jwe    jweServices
}

func (s *jwtInteractor) GenerateToken(userID uint, userPermissions domain.UserPermissions) (respToken []byte, svcErr *svcerr.ServiceError) {

	token, err := s.getPayload(userID, userPermissions)
	if err != nil {
		svcErr = svcerr.New(
			errors.Wrap(err, internal.ErrMsgGeneratingPayload.Error()),
			internal.ErrorCategoryInternal,
		)
		return
	}

	signedToken, err := s.jws.signToken(token, []byte(s.config.SigningSecretKey()))
	if err != nil {
		svcErr = svcerr.New(
			errors.Wrap(err, internal.ErrMsgSignFailed.Error()),
			internal.ErrorCategoryInternal,
		)
		return
	}

	encryptedToken, err := s.jwe.encryptToken(signedToken, []byte(s.config.EncryptionSecretKey()))
	if err != nil {
		svcErr = svcerr.New(
			errors.Wrap(err, internal.ErrMsgFailedToEncryptPayload.Error()),
			internal.ErrorCategoryInternal,
		)
		return
	}
	respToken = encryptedToken
	return
}

func (s *jwtInteractor) getPayload(userID uint, userPermissions domain.UserPermissions) (jwt.Token, error) {

	token := jwt.New()
	err := token.Set(jwt.IssuerKey, issuerKeyValue) // security token service (STS)
	if err != nil {
		return nil, errors.Wrap(err, "token issuer value set failed")
	}
	err = token.Set(jwt.SubjectKey, strconv.Itoa(int(userID)))
	if err != nil {
		return nil, errors.Wrap(err, "token subject value set failed")
	}
	err = token.Set(jwt.ExpirationKey, time.Now().Add(s.config.ExpirationDuration()))
	if err != nil {
		return nil, errors.Wrap(err, "token expiration value set failed")
	}
	err = token.Set(permissionsKeyName, userPermissions)
	if err != nil {
		return nil, errors.Wrap(err, "token permissions value set failed")
	}
	return token, nil
}

func (s *jwtInteractor) ValidateToken(token string) (permissions interface{}, svcErr *svcerr.ServiceError) {

	decryptedToken, err := s.jwe.decryptToken([]byte(token), []byte(s.config.EncryptionSecretKey()))
	if err != nil {
		svcErr = svcerr.New(
			errors.Wrap(err, internal.ErrMsgFailedToDecryptToken.Error()),
			internal.ErrorCategoryInvalidToken,
		)
		return
	}

	tokenObj, err := s.jws.verifyToken(decryptedToken, []byte(s.config.SigningSecretKey()))
	if err != nil {
		svcErr = svcerr.New(
			errors.Wrap(err, internal.ErrMsgInvalidToken.Error()),
			internal.ErrorCategoryInvalidToken,
		)
		return
	}
	permissions, ok := tokenObj.Get(permissionsKeyName)
	if !ok {
		permissions, svcErr = nil, svcerr.New(
			internal.ErrMsgTokenDoesntClaimAnyPermissions,
			internal.ErrorCategorySemanticallyUnprocesable,
		)
	}
	return
}
