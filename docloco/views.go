package docloco

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"html/template"
)

type displayResult struct {
	Url string
	Title string
	Highlight template.HTML
	Project string
	Version string
}

func searchView(c *gin.Context) {
	// search for some text
	queryString := c.Query("q")
	if queryString != "" {
		searchResults, err := searchController(queryString)
		if err != nil {
			c.JSON(400, gin.H{
				"status": "failed",
				"error":  "something went wrong",
			})}
		fmt.Println(searchResults)
		results := make([]displayResult, len(searchResults.Hits))
		for i := 0; i < len(searchResults.Hits); i++ {
			var highlight template.HTML
			if len(searchResults.Hits[i].Fragments["PlainContent"]) > 0 {
				highlight = template.HTML(searchResults.Hits[i].Fragments["PlainContent"][0])
			} else {
				highlight = template.HTML("")
			}

			results[i] = displayResult{
				Url: searchResults.Hits[i].ID,
				Title: string(searchResults.Hits[i].Fields["Title"].(string)),
				Highlight: highlight,
				Project: string(searchResults.Hits[i].Fields["Project"].(string)),
				Version: string(searchResults.Hits[i].Fields["Version"].(string)),
			}
		}
		c.HTML(http.StatusOK, "results.html", gin.H{
			"results": results,
			"query":   queryString,
			"Took": searchResults.Took,
			"Found": len(searchResults.Hits),
		})
	} else {
		// Render search form
		c.HTML(http.StatusOK, "index.html", gin.H{})
	}
}

func uploadView(c *gin.Context) {
	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		log.Fatal(err)
	}
	filename := header.Filename
	name := c.Request.FormValue("name")
	version := c.Request.FormValue("version")
	if err := uploadController(&file, filename, name, version); err != nil {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	c.JSON(200, gin.H{
		"status": "success",
		"error":  "Upload succeeded",
	})

}
