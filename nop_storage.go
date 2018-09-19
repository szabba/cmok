package cmok

import (
	"errors"
	"io"
)

var _ErrNop = errors.New("nop implementation")

type NopStorage struct{}

var _ Storage = new(NopStorage)

func NewNopStorage() *NopStorage {
	return &NopStorage{}
}

func (_ *NopStorage) Get(_ string) ([]Entry, io.ReadCloser, error) {
	return nil, nil, _ErrNop
}

func (_ *NopStorage) Set(_ string, _ io.ReadCloser) error {
	return _ErrNop
}
