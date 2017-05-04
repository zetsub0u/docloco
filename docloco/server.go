package docloco

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/blevesearch/bleve"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/zetsub0u/docloco/config"
	"html/template"
	"net/http"
)

var idx bleve.Index

// The glob pattern used to find files to be indexed inside the supplied documentation package
const globPattern string = "**/*.html"

// Load templates from rice.go, this will either read the templates from the local filesystem
// or from the embeded generated go source code if rice embed-go is run before building.
func loadTemplates(list ...string) multitemplate.Render {
	templateBox, err := rice.FindBox("../templates")
	if err != nil {
		fmt.Println(err)
	}

	r := multitemplate.New()

	for _, x := range list {
		templateString, err := templateBox.String(x)
		if err != nil {
			fmt.Println(err)
		}

		tmplMessage, err := template.New(x).Parse(templateString)
		if err != nil {
			fmt.Println(err)
		}

		r.Add(x, tmplMessage)
	}

	return r
}

// Starts the Web Server
func RunServer() {
	if !config.Store.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	var err error
	if idx, err = getIndex(); err != nil {
		panic(err)
	}

	r := gin.Default()

	// Static Views
	//r.LoadHTMLGlob("templates/*")
	r.HTMLRender = loadTemplates("index.html", "results.html")
	r.StaticFS("/docs", http.Dir(config.Store.StorageDir))

	// Dynamic Views
	r.GET("/", searchView)
	r.POST("/upload", uploadView)

	r.Run(fmt.Sprintf("%s:%d", config.Store.Server.Host, config.Store.Server.Port))
}
