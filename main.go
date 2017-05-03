package main

import (
	"github.com/docloco/utils"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/analysis/char/html"
	"github.com/blevesearch/bleve/analysis/token/lowercase"
	"github.com/blevesearch/bleve/analysis/tokenizer/unicode"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-zglob"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"net/http"
)

func getIndex() (bleve.Index, error) {
	// open a new index
	var index bleve.Index
	var err error
	if index, err = bleve.Open("docs.bleve"); err != nil {
		mapping := bleve.NewIndexMapping()
		err = mapping.AddCustomAnalyzer("html", map[string]interface{}{
			"type": custom.Name,
			"char_filters": []string{
				html.Name,
			},
			"tokenizer": unicode.Name,
			"token_filters": []string{
				lowercase.Name,
			},
		})
		if err != nil {
			panic(err)
		}
		mapping.DefaultAnalyzer = "html"
		index, err = bleve.New("docs.bleve", mapping)
	}
	return index, err
}

func indexFile(path string, index bleve.Index) error {
	dat, err := ioutil.ReadFile(path)
	data := struct {
		Content string
		Path    string
	}{
		Content: string(dat),
		Path:    path,
	}
	fmt.Println(path)

	index.Index(path, data)
	return err
}

func saveFile(file multipart.File, filename string) string {
	dest := "./tmp/" + filename
	out, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		panic(err)
	}
	return dest
}

func main() {
	var idx bleve.Index
	var err error
	if idx, err = getIndex(); err != nil {
		panic(err)
	}

	glob := "**/*.html"

	r := gin.Default()
	r.GET("/search", func(c *gin.Context) {
		// search for some text
		query := bleve.NewMatchQuery("Exception")
		search := bleve.NewSearchRequest(query)
		searchResults, err := idx.Search(search)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(searchResults)
		c.JSON(200, gin.H{
			"message": searchResults,
		})
	})

	r.POST("/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("upload")
		if err != nil {
			log.Fatal(err)
		}
		filename := header.Filename
		fmt.Println(header.Filename)
		zipFile := saveFile(file, filename)
		name := c.Request.FormValue("name")
		version := c.Request.FormValue("version")
		dest := filepath.Join("docs", name, version)
		os.MkdirAll(dest, os.ModePerm)
		utils.Unzip(zipFile, dest)
		globPath := filepath.Join(dest, glob)
		fmt.Println(globPath)
		matches, _ := zglob.Glob(globPath)
		fmt.Printf("%v", matches)
		for _, htmlPath := range matches {
			indexFile(htmlPath, idx)
		}

	})

	r.StaticFS("/docs", http.Dir("docs"))

	r.Run() // listen and serve on 0.0.0.0:8080

	//matches, _ := zglob.Glob(path)

	/*	for _, htmlPath := range matches {
			indexFile(htmlPath, idx)
		}

		// search for some text
		query := bleve.NewMatchQuery("Exception")
		search := bleve.NewSearchRequest(query)
		searchResults, err := idx.Search(search)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(searchResults)*/
}
