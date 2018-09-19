package userlist

import (
	"github.com/szabba/cmok/auth"
)

type AuthConfig struct {
	Realm string                   `json:"realm"`
	Users map[auth.User]UserConfig `json:"users"`
}

type AccessConfig struct {
	Permissions map[auth.User]Permissions `json:"permissions"`
}

type UserConfig struct {
	Password Password `json:"password"`
}

type Password string

const NoPassword = Password("")

type Permissions struct {
	CanRead  bool `json:"can_read"`
	CanWrite bool `json:"can_write"`
}

func None() Permissions {
	return Permissions{}
}

func Read() Permissions {
	return Permissions{CanRead: true}
}

func Write() Permissions {
	return Permissions{CanWrite: true}
}

func All() Permissions {
	return Permissions{CanRead: true, CanWrite: true}
}
