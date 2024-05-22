package cmd

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	listenPort = "1302"
)

//go:embed templates/**
var templateFolder embed.FS

func server() {
	r := gin.Default()
	r.Static("/css", "./cmd/templates/css")
	LoadHTMLFromEmbedFS(r, templateFolder, "templates/*.html")

	r.GET("/", returnIndex)
	r.POST("/api/upload", getCert)

	fmt.Printf("Listening on port %v...\n", listenPort)
	r.Run(":" + listenPort)
}

func returnIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "templates/index.html", gin.H{
		"UserAgent": c.Request.UserAgent(),
	})
}

func getCert(c *gin.Context) {
	text, _ := c.GetPostForm("domain")
	fmt.Println(text)
}
