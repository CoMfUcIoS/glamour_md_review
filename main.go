package main

import (
	"embed"
	"net/http"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/gin-gonic/gin"
)

// Files contains the embedded files for the web server
//
//go:embed xterm
var Files embed.FS

// index.html
//
//go:embed index.html
var Index embed.FS

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Static files from embed.FS
	r.StaticFS("/static", http.FS(Files))

	// POST /api/md
	r.POST("/api/md", func(c *gin.Context) {
		// Get the input from the request
		var input struct {
			// Lines []string `json:"lines"`
			Markdown string `json:"markdown"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		r, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(100),
			glamour.WithStylesFromJSONBytes([]byte(`{
			"document": {
				"margin": 0
			}
		}`)),
		)

		// Render the markdown
		rendered, _ := r.Render(input.Markdown)

		// Remove all \n
		rendered = strings.ReplaceAll(rendered, "\n", "")

		// Return the rendered html
		c.JSON(http.StatusOK, gin.H{"markdown": rendered})

		// // iterate over the lines and render them
		// var lines []string
		// for _, line := range input.Lines {
		// 	rendered, _ := r.Render(line)
		// 	//remove all \n
		// 	rendered = strings.ReplaceAll(rendered, "\n", "")
		// 	lines = append(lines, rendered)
		// }

		// // Return the rendered html
		// c.JSON(http.StatusOK, gin.H{"lines": lines})
	})

	// GET /
	r.GET("/", func(c *gin.Context) {
		// Get the index.html file from the embedded files
		index, err := Index.ReadFile("index.html")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the index.html file
		c.Data(http.StatusOK, "text/html", index)
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
