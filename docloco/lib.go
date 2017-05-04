package docloco

import (
	"fmt"
	"github.com/blevesearch/bleve"
	custom_analyzer "github.com/blevesearch/bleve/analysis/analyzer/custom"
	html_filter "github.com/blevesearch/bleve/analysis/char/html"
	lowercase_filter "github.com/blevesearch/bleve/analysis/token/lowercase"
	unicode_tokenizer "github.com/blevesearch/bleve/analysis/tokenizer/unicode"
	"github.com/zetsub0u/docloco/config"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strings"
)


// Open or create a new Bleve Index repository, it also initializes the mapper analyzers.
func getIndex() (bleve.Index, error) {
	// open a new index
	var index bleve.Index
	var err error
	if index, err = bleve.Open(config.Store.IndexDir); err != nil {

		mapping := bleve.NewIndexMapping()

		err = mapping.AddCustomAnalyzer("html", map[string]interface{}{
			"type": custom_analyzer.Name,
			"char_filters": []string{
				html_filter.Name,
			},
			"tokenizer": unicode_tokenizer.Name,
			"token_filters": []string{
				lowercase_filter.Name,
			},
		})
		if err != nil {
			panic(err)
		}
		mapping.DefaultAnalyzer = "html"
		index, err = bleve.New(config.Store.IndexDir, mapping)
	}
	return index, err
}

// Index the contents of an html file with some basic parsing
func indexFile(path string, index bleve.Index) error {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	data := struct {
		_type   string
		Content string
		Path    string
		Title   string
	}{
		_type:   "doc",
		Content: string(dat),
		Path:    path,
		Title:   getTitle(string(dat)),
	}
	fmt.Println(path)

	err = index.Index(path, data)
	return err
}

// Save the received zipfile in the filesystem to be later processed
func saveFile(file *multipart.File, filename string) (string, error) {
	tmp := "/tmp/" + filename
	out, err := os.Create(tmp)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err := io.Copy(out, *file); err != nil {
		return "", err
	}
	return tmp, nil
}


// Extract the title
func getTitle(data string) string {
	doc, err := html.Parse(strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node) string
	f = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "title" {
			fmt.Printf("==> %s\n", n.FirstChild.Data)
			return n.FirstChild.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			xx := f(c)
			if len(xx) > 0 {
				return xx
			}
		}
		return ""
	}
	xx := f(doc)
	fmt.Printf("--> %s\n", xx)
	return xx
}


// Search the index for a given keyword
func doSearch(queryString string) (*bleve.SearchResult, error) {
	query := bleve.NewMatchQuery(queryString)
	search := bleve.NewSearchRequest(query)
	search.Highlight = bleve.NewHighlightWithStyle(html_filter.Name)
	search.Fields = []string{"*"}
	return idx.Search(search)
}
