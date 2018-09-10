package cmok

import (
	"errors"
	"github.com/szabba/cmok/auth"
)

var ErrAccessDenied = errors.New("access denied")

type AccessPolicy interface {
	Protect(storage Storage, user auth.User) Storage
}
