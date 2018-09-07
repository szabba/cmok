package userlist_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/szabba/assert"

	"github.com/szabba/cmok"
	"github.com/szabba/cmok/userlist"
)

func TestAuthServiceAuthenticatesAnAnonymousUserWhenNoCredentialsArePresent(t *testing.T) {
	// given
	authSvc := userlist.NewAuthService(userlist.Config{})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	// when
	user, ok := authSvc.Authenticate(w, r)

	// then
	assertAuthOK(t, user, ok, cmok.AnonymousUser)
	assertNoOutput(t, w)
}

func TestAuthServiceFailsToAuthenticateAUserWithAnInvalidHeaderFormat(t *testing.T) {
	// given
	authSvc := userlist.NewAuthService(userlist.Config{})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Sezame, open!")

	// when
	user, ok := authSvc.Authenticate(w, r)

	// then
	assertAuthFailed(t, user, ok, w)
}

func TestAuthServiceFailsToAuthenticateUserWithWrongCredentials(t *testing.T) {
	// given
	config := userlist.Config{
		Users: map[cmok.User]userlist.UserConfig{
			"uploader": {Password: "download"},
		},
	}
	authSvc := userlist.NewAuthService(config)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Basic dXBsb2FkZXI6dXBsb2Fk")

	// when
	user, ok := authSvc.Authenticate(w, r)

	// then
	assertAuthFailed(t, user, ok, w)
}

func TestAuthServiceAuthenticatesUserWithRightCredentials(t *testing.T) {
	// given
	config := userlist.Config{
		Users: map[cmok.User]userlist.UserConfig{
			"uploader": {Password: "upload"},
		},
	}
	authSvc := userlist.NewAuthService(config)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Basic dXBsb2FkZXI6dXBsb2Fk")

	// when
	user, ok := authSvc.Authenticate(w, r)

	// then
	assertAuthOK(t, user, ok, cmok.User("uploader"))
	assertNoOutput(t, w)
}

func assertAuthOK(t *testing.T, user cmok.User, ok bool, wantUser cmok.User) {
	assert.That(ok, t.Errorf, "authentication failed, should succeed")
	assert.That(user == wantUser, t.Errorf, "got user %q, want %q", user, wantUser)
}

func assertNoOutput(t *testing.T, w *httptest.ResponseRecorder) {
	body := w.Body.String()
	assert.That(body == "", t.Errorf, "the service has written %q, want no output", body)
}

func assertAuthFailed(t *testing.T, user cmok.User, ok bool, w *httptest.ResponseRecorder) {
	assert.That(!ok, t.Errorf, "authentication succeeded, should fail")
	assert.That(user == cmok.AnonymousUser, t.Errorf, "got user %q, want %q", user, cmok.AnonymousUser)

	status := w.Code
	assert.That(
		status == http.StatusUnauthorized, t.Errorf,
		"got status %d %s, want %d %s",
		status, http.StatusText(status),
		http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
}
