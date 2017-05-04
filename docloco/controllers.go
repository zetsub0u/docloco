package docloco

import (
	"github.com/blevesearch/bleve"
	"github.com/mattn/go-zglob"
	"github.com/zetsub0u/docloco/config"
	"mime/multipart"
	"os"
	"path/filepath"
)

func uploadController(file *multipart.File, filename, name, version string) error {
	zipFile, err := saveFile(file, filename)
	defer os.Remove(zipFile)
	if err != nil {
		return err
	}
	dest := filepath.Join(config.Store.StorageDir, name, version)
	os.MkdirAll(dest, os.ModePerm)
	if err := Unzip(zipFile, dest); err != nil {
		return err
	}
	globPath := filepath.Join(dest, globPattern)
	matches, _ := zglob.Glob(globPath)
	for _, htmlPath := range matches {
		if err := indexFile(htmlPath, idx); err != nil {
			return err
		}
	}
	return nil
}

func searchController(queryString string) (*bleve.SearchResult, error) {
	return doSearch(queryString)
}
