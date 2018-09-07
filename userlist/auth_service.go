package userlist

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/szabba/cmok"
)

const (
	_AuthorizationHeader         = "Authorization"
	_WWWAuthenticateHeader       = "WWW-Authenticate"
	_WWWAuthenticateHeaderFormat = `Basic realm=%q`
	_HeaderPrefix                = "Basic"
	_CredentialsSeparator        = ":"
)

var (
	_ErrFailedAuth = errors.New("failed authentication")
)

type AuthService struct {
	config Config
}

func NewAuthService(config Config) *AuthService {
	return &AuthService{config}
}

var _ cmok.AuthService = new(AuthService)

func (svc *AuthService) Authenticate(w http.ResponseWriter, r *http.Request) (cmok.User, bool) {
	user, ok := svc.authenticate(r)
	if !ok {
		svc.writeUnauthorized(w)
	}
	return user, ok
}

func (svc *AuthService) authenticate(r *http.Request) (cmok.User, bool) {
	user, password, ok := svc.getCredentials(r)
	if user == cmok.AnonymousUser {
		return cmok.AnonymousUser, ok
	}

	userDetails, present := svc.config.Users[user]
	if !present || password != userDetails.Password {
		return cmok.AnonymousUser, false
	}

	return user, true
}

func (svc *AuthService) writeUnauthorized(w http.ResponseWriter) {
	wwwAuthentiacateHeaderValue := fmt.Sprintf(_WWWAuthenticateHeaderFormat, svc.config.Realm)
	w.Header().Set(_WWWAuthenticateHeader, wwwAuthentiacateHeaderValue)
	http.Error(w, "authentication failed", http.StatusUnauthorized)
}

func (svc *AuthService) getCredentials(r *http.Request) (cmok.User, Password, bool) {
	authHeader := r.Header.Get(_AuthorizationHeader)
	if authHeader == "" {
		return cmok.AnonymousUser, NoPassword, true
	}

	encodedCredentials, ok := svc.extractCredentials(authHeader)
	if !ok {
		return cmok.AnonymousUser, NoPassword, false
	}

	credentials, ok := svc.decodeCredentials(encodedCredentials)
	if !ok {
		return cmok.AnonymousUser, NoPassword, false
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

func (svc *AuthService) splitCredentials(creds string) (cmok.User, Password, bool) {
	parts := strings.SplitN(creds, _CredentialsSeparator, 2)
	if len(parts) != 2 {
		return cmok.AnonymousUser, NoPassword, false
	}
	return cmok.User(parts[0]), Password(parts[1]), true
}
