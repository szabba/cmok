package cmok

import (
	"net/http"
)

type AuthService interface {
	Authenticate(w http.ResponseWriter, r *http.Request) (User, bool)
}

type User string

const AnonymousUser = User("")
