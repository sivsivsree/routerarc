package utils

import (
	"net/http"
	"strings"
)

type FileSystem struct {
	FS http.FileSystem
}

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

type StaticFileServer struct {
}

func (sfs *StaticFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
