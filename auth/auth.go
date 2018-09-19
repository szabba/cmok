package auth

import (
	"net/http"
)

type Service interface {
	Authenticate(r *http.Request) (User, bool)
}

type User string

const AnonymousUser = User("")
