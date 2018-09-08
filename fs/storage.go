package fs

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/szabba/cmok"
)

type Storage struct {
	rootPath string
}

func NewStorage(rootPath string) *Storage {
	return &Storage{rootPath}
}

func (st *Storage) Get(path string) ([]cmok.Entry, io.ReadCloser, error) {
	children, err := st.list(path)
	if err == nil {
		return children, nil, nil
	}
	content, err := st.getOne(path)
	return nil, content, err
}

func (st *Storage) list(path string) ([]cmok.Entry, error) {
	localPath, err := st.sanitize(path)
	if err != nil {
		return nil, fmt.Errorf("cannot sanitize path %q", path)
	}

	f, err := os.Open(localPath)
	if err != nil {
		return nil, fmt.Errorf("cannot access %q", path)
	}
	defer f.Close()

	fis, err := f.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("cannot list children of %q", path)
	}

	entries := make([]cmok.Entry, len(fis))
	for i, fi := range fis {
		entries[i] = st.fileInfoToEntry(path, fi)
	}

	return entries, nil
}

func (st *Storage) fileInfoToEntry(at string, fi os.FileInfo) cmok.Entry {
	name := fi.Name()
	if fi.IsDir() {
		name += "/"
	}
	return cmok.Entry{
		Name: name,
		Path: path.Join(at, fi.Name()),
	}
}

func (st *Storage) getOne(path string) (io.ReadCloser, error) {
	localPath, err := st.sanitize(path)
	if err != nil {
		return nil, fmt.Errorf("cannot sanitize path %q", err)
	}

	f, err := os.Open(localPath)
	if err != nil {
		return nil, fmt.Errorf("cannot access %q", path)
	}

	return f, nil
}

func (st *Storage) Set(path string, content io.ReadCloser) error {
	localPath, err := st.sanitize(path)
	if err != nil {
		return fmt.Errorf("cannot sanitize path %q", err)
	}

	f, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("cannot access %q", path)
	}
	defer st.close(f, &err, path)

	_, err = io.Copy(f, content)
	if err != nil {
		err = fmt.Errorf("cannot write %q", path)
	}
	return err
}

func (st *Storage) sanitize(path string) (string, error) {
	fullRoot, err := filepath.Abs(st.rootPath)
	if err != nil {
		return "", err
	}

	joined := filepath.Clean(filepath.Join(fullRoot, path))
	_, err = filepath.Rel(fullRoot, joined)
	if err != nil {
		return "", err
	}

	return joined, nil
}

func (st *Storage) close(f *os.File, err *error, path string) {
	closeErr := f.Close()
	if closeErr != nil && *err == nil {
		*err = fmt.Errorf("cannot write %q", path)
	}
}
