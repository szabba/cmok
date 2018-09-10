package userlist_test

import (
	"net/http/httptest"
	"testing"

	"github.com/szabba/assert"

	"github.com/szabba/cmok/auth"
	"github.com/szabba/cmok/userlist"
)

func TestAuthServiceAuthenticatesAnAnonymousUserWhenNoCredentialsArePresent(t *testing.T) {
	// given
	authSvc := userlist.NewAuthService(userlist.AuthConfig{})

	r := httptest.NewRequest("GET", "/", nil)

	// when
	user, ok := authSvc.Authenticate(r)

	// then
	assertAuthOK(t, user, ok, auth.AnonymousUser)
}

func TestAuthServiceFailsToAuthenticateAUserWithAnInvalidHeaderFormat(t *testing.T) {
	// given
	authSvc := userlist.NewAuthService(userlist.AuthConfig{})

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Sezame, open!")

	// when
	user, ok := authSvc.Authenticate(r)

	// then
	assertAuthFailed(t, user, ok)
}

func TestAuthServiceFailsToAuthenticateUserWithWrongCredentials(t *testing.T) {
	// given
	config := userlist.AuthConfig{
		Users: map[auth.User]userlist.UserConfig{
			"uploader": {Password: "download"},
		},
	}
	authSvc := userlist.NewAuthService(config)

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Basic dXBsb2FkZXI6dXBsb2Fk")

	// when
	user, ok := authSvc.Authenticate(r)

	// then
	assertAuthFailed(t, user, ok)
}

func TestAuthServiceAuthenticatesUserWithRightCredentials(t *testing.T) {
	// given
	config := userlist.AuthConfig{
		Users: map[auth.User]userlist.UserConfig{
			"uploader": {Password: "upload"},
		},
	}
	authSvc := userlist.NewAuthService(config)

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Basic dXBsb2FkZXI6dXBsb2Fk")

	// when
	user, ok := authSvc.Authenticate(r)

	// then
	assertAuthOK(t, user, ok, "uploader")
}

func assertAuthOK(t *testing.T, user auth.User, ok bool, wantUser auth.User) {
	assert.That(ok, t.Errorf, "authentication failed, should succeed")
	assert.That(user == wantUser, t.Errorf, "got user %q, want %q", user, wantUser)
}

func assertAuthFailed(t *testing.T, user auth.User, ok bool) {
	assert.That(!ok, t.Errorf, "authentication succeeded, should fail")
	assert.That(user == auth.AnonymousUser, t.Errorf, "got user %q, want %q", user, auth.AnonymousUser)

}
