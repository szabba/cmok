package cmok

import (
	"net/http"
)

type AuthService interface {
	Authenticate(w http.ResponseWriter, r *http.Request) (User, bool)
}

type AccessPolicy interface {
	Protect(storage Storage, user User) Storage
}

type User string

const AnonymousUser = User("")
