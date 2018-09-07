package userlist

import (
	"github.com/szabba/cmok"
)

type Config struct {
	Realm string                   `json:"realm"`
	Users map[cmok.User]UserConfig `json:"users"`
}

type UserConfig struct {
	Password Password `json:"password"`
}

type Password string

const NoPassword = Password("")
