/*
Copyright © 2024 raul
*/

package cmd

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	listenPort = "1302"
)

//go:embed templates/**
var templateFolder embed.FS

func server() {
	lPort := viper.GetString("Server.port")
	if lPort != "" {
		listenPort = lPort
	}
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
