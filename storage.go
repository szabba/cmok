package cmok

import (
	"io"
)

type Entry struct {
	Path string
	Name string
}

type Storage interface {
	List(path string) ([]Entry, error)
	Get(path string) (io.ReadCloser, error)
	Put(path string, content io.ReadCloser) error
}
