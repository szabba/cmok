package userlist

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/szabba/cmok/auth"
)

const (
	_AuthorizationHeader  = "Authorization"
	_HeaderPrefix         = "Basic"
	_CredentialsSeparator = ":"
)

var (
	_ErrFailedAuth = errors.New("failed authentication")
)

type AuthService struct {
	config AuthConfig
}

func NewAuthService(config AuthConfig) *AuthService {
	return &AuthService{config}
}

var _ auth.Service = new(AuthService)

func (svc *AuthService) Authenticate(r *http.Request) (auth.User, bool) {
	user, password, ok := svc.getCredentials(r)
	if user == auth.AnonymousUser {
		return auth.AnonymousUser, ok
	}

	userDetails, present := svc.config.Users[user]
	if !present || password != userDetails.Password {
		return auth.AnonymousUser, false
	}

	return user, true
}

func (svc *AuthService) getCredentials(r *http.Request) (auth.User, Password, bool) {
	authHeader := r.Header.Get(_AuthorizationHeader)
	if authHeader == "" {
		return auth.AnonymousUser, NoPassword, true
	}

	encodedCredentials, ok := svc.extractCredentials(authHeader)
	if !ok {
		return auth.AnonymousUser, NoPassword, false
	}

	credentials, ok := svc.decodeCredentials(encodedCredentials)
	if !ok {
		return auth.AnonymousUser, NoPassword, false
	}

	return svc.splitCredentials(credentials)
}

func (svc *AuthService) extractCredentials(header string) (string, bool) {
	out := header
	out = strings.TrimSpace(out)
	if !strings.HasPrefix(out, _HeaderPrefix) {
		return "", false
	}
	out = strings.TrimPrefix(out, _HeaderPrefix)
	out = strings.TrimSpace(out)
	return out, true
}

func (svc *AuthService) decodeCredentials(encoded string) (string, bool) {
	dec := base64.NewDecoder(base64.URLEncoding, strings.NewReader(encoded))
	decoded, err := ioutil.ReadAll(dec)
	return string(decoded), err == nil
}

func (svc *AuthService) splitCredentials(creds string) (auth.User, Password, bool) {
	parts := strings.SplitN(creds, _CredentialsSeparator, 2)
	if len(parts) != 2 {
		return auth.AnonymousUser, NoPassword, false
	}
	return auth.User(parts[0]), Password(parts[1]), true
}
