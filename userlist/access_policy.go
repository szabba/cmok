package userlist

import (
	"io"

	"github.com/szabba/cmok"

	"github.com/szabba/cmok/auth"
)

type AccessPolicy struct {
	config AccessConfig
}

func NewAccessPolicy(config AccessConfig) *AccessPolicy {
	return &AccessPolicy{config}
}

var _ cmok.AccessPolicy = new(AccessPolicy)

func (p *AccessPolicy) Protect(storage cmok.Storage, user auth.User) cmok.Storage {
	return &protectedStorage{
		p.config.Permissions[user],
		storage,
	}
}

type protectedStorage struct {
	perm    Permissions
	storage cmok.Storage
}

func (p *protectedStorage) Get(path string) ([]cmok.Entry, io.ReadCloser, error) {
	if !p.perm.CanRead {
		return nil, nil, cmok.ErrAccessDenied
	}
	return p.storage.Get(path)
}

func (p *protectedStorage) Set(path string, r io.ReadCloser) error {
	if !p.perm.CanWrite {
		return cmok.ErrAccessDenied
	}
	return p.storage.Set(path, r)
}
