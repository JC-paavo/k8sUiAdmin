package main

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var distFS embed.FS

func serveFrontend(r *gin.Engine) {
	staticFiles, err := fs.Sub(distFS, "dist")
	if err != nil {
		return
	}

	fileServer := http.FileServer(http.FS(staticFiles))

	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.Next()
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)

		if c.Writer.Status() == 404 {
			indexBytes, err := distFS.ReadFile("dist/index.html")
			if err != nil {
				c.String(http.StatusNotFound, "Not Found")
				c.Abort()
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", indexBytes)
			c.Abort()
		}
	})
}
