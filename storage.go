package cmok

import (
	"io"
)

type Entry struct {
	Path string
	Name string
}

type Storage interface {
	Get(path string) ([]Entry, io.ReadCloser, error)
	Set(path string, content io.ReadCloser) error
}
