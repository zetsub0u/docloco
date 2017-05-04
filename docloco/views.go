package docloco

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"html/template"
	"github.com/jaytaylor/html2text"
)

type displayResult struct {
	Url string
	Title string
	Highlight template.HTML
}

func cleanHighlight(input string) template.HTML {
	txt, _ := html2text.FromString(input)
	return template.HTML(txt)
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
			results[i] = displayResult{
				Url: searchResults.Hits[i].ID,
				Title: string(searchResults.Hits[i].Fields["Title"].(string)),
				Highlight: cleanHighlight(searchResults.Hits[i].Fragments["Content"][0]),
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
