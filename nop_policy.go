package cmok

import (
	"github.com/szabba/cmok/auth"
)

type NopPolicy struct{}

var _ AccessPolicy = NopPolicy{}

func (_ NopPolicy) Protect(storage Storage, _ auth.User) Storage {
	return storage
}
