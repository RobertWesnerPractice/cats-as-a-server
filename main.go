package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed cats/*
var fCats embed.FS

//go:embed assets/* index.html
var fAssets embed.FS

func main() {
	router := gin.Default()
	router.StaticFS("/public", http.FS(fAssets))

	{
		entries, err := fCats.ReadDir("cats")
		if err != nil {
			panic(err)
		}

		// Buffer everything here since it cant change at runtime
		filenameMap := make(map[string]string, len(entries))
		list := make([]string, len(entries))
		for i, entry := range entries {
			filename := entry.Name()
			name := strings.TrimSuffix(filename, filepath.Ext(filename))
			filenameMap[name] = filename
			list[i] = name
		}

		api := router.Group("/api")

		api.GET("/list", func(c *gin.Context) {
			c.JSON(http.StatusOK, list)
		})

		api.GET("/get/:name", func(c *gin.Context) {
			filename, ok := filenameMap[c.Param("name")]
			if !ok {
				c.Data(http.StatusNotFound, "", nil)

				return
			}

			c.Param("file")
			file, err := fCats.ReadFile("cats/" + filename)
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
