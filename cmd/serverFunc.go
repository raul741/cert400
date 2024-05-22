/*
Copyright © 2024 raul
*/

package cmd

import (
	"embed"
	"fmt"
	"net/http"
	"os"

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
	LoadHTMLFromEmbedFS(r, templateFolder, "templates/*.html")
	r.Static("/css", "./cmd/templates/css")

	r.GET("/", returnIndex)
	r.GET("/cacert", returnCacert)
	r.POST("/api/upload", getCert)

	fmt.Printf("Listening on port %v...\n", listenPort)
	r.Run(":" + listenPort)
}

func returnCacert(c *gin.Context) {
	home, err := os.UserHomeDir()
	if err != nil {
		e := fmt.Sprintf("Error happened fetching: %v\n", err)
		c.String(http.StatusInternalServerError, e)
		return
	}
	c.File(home + "/.config/cert400/ca.crt")
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
