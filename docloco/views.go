package docloco

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func searchView(c *gin.Context) {
	// search for some text
	queryString := c.Query("q")
	if queryString != "" {
		searchResults, err := searchController(queryString)
		if err != nil {
			c.JSON(400, gin.H{
				"status": "failed",
				"error":  "something went wrong",
			})
		}
		fmt.Println(searchResults)
		c.HTML(http.StatusOK, "results.html", gin.H{
			"results": searchResults,
			"query":   queryString})
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
