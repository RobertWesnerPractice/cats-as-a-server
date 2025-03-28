package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"net/http"
)

//go:embed cats/*
var fCats embed.FS

//go:embed assets/* index.html
var fAssets embed.FS

func main() {
	router := gin.Default()
	router.StaticFS("/public", http.FS(fAssets))

	{
		api := router.Group("/api")

		// TODO: remove .jpeg, etc. from file
		api.GET("/list", func(c *gin.Context) {
			entries, err := fCats.ReadDir("cats")
			if err != nil {
				c.JSON(http.StatusInternalServerError, struct {
					Message string `json:"message"`
				}{err.Error()})

				return
			}
			results := make([]string, len(entries))

			for i, entry := range entries {
				results[i] = entry.Name()
			}

			c.JSON(http.StatusOK, results)
		})

		// TODO: find by name without extension
		api.GET("/get/:file", func(c *gin.Context) {
			file, err := fCats.ReadFile("cats/" + c.Param("file"))
			if err != nil {
				c.Data(http.StatusNotFound, gin.MIMEHTML, []byte(err.Error()))
			}

			c.Data(http.StatusOK, http.DetectContentType(file), file)
		})
	}

	router.GET("/", func(c *gin.Context) {
		file, _ := fAssets.ReadFile("index.html")
		c.Data(http.StatusOK, gin.MIMEHTML, file)
	})

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
