package userlist_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/szabba/assert"

	"github.com/szabba/cmok"
	"github.com/szabba/cmok/auth"
	"github.com/szabba/cmok/userlist"
)

func TestAccessPolicyDeniesReadAccessToUserWithoutReadPermissions(t *testing.T) {
	// given
	policy := userlist.NewAccessPolicy(userlist.AccessConfig{})

	storage := policy.Protect(cmok.NewNopStorage(), auth.AnonymousUser)

	// when
	_, _, err := storage.Get("/file")

	// then
	assert.That(err == cmok.ErrAccessDenied, t.Errorf, "got error %q, want %q", err, cmok.ErrAccessDenied)
}

func TestAccessPolicyDeniesWriteAccessToUserWithoutWritePermissions(t *testing.T) {
	// given
	policy := userlist.NewAccessPolicy(userlist.AccessConfig{})

	storage := policy.Protect(cmok.NewNopStorage(), auth.AnonymousUser)

	// when
	err := storage.Set("/file", ioutil.NopCloser(strings.NewReader("")))

	// then
	assert.That(err == cmok.ErrAccessDenied, t.Errorf, "got error %q, want %q", err, cmok.ErrAccessDenied)
}

func TestAccessPolicyAllowsReadAccessToUserWithReadPermission(t *testing.T) {
	// given
	config := userlist.AccessConfig{
		Permissions: map[auth.User]userlist.Permissions{
			"downloader": userlist.Read(),
		},
	}

	policy := userlist.NewAccessPolicy(config)

	storage := policy.Protect(cmok.NewNopStorage(), "downloader")

	// when
	_, _, err := storage.Get("/file")

	// then
	assert.That(err != cmok.ErrAccessDenied, t.Errorf, "error should not be %q", err)
}

func TestAccessPolicyAllowsWriteAccessToUserWithWritePermissions(t *testing.T) {
	// given
	config := userlist.AccessConfig{
		Permissions: map[auth.User]userlist.Permissions{
			"uploader": userlist.Write(),
		},
	}
	policy := userlist.NewAccessPolicy(config)

	storage := policy.Protect(cmok.NewNopStorage(), "uploader")

	content := ioutil.NopCloser(strings.NewReader(""))

	// when
	err := storage.Set("/file", content)

	// then
	assert.That(err != cmok.ErrAccessDenied, t.Errorf, "errors should not be %q", cmok.ErrAccessDenied)
}
