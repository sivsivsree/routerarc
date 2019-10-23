package utils

import (
	"net/http"
	"strings"
)

// FileSystem is used for static file serving
type FileSystem struct {
	FS http.FileSystem
}

// Open is used to implement the Filesystem
func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.FS.Open(path)
	if err != nil {
		return nil, err
	}

	s, errStat := f.Stat()

	if errStat != nil {
		return nil, err
	}

	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.FS.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}
