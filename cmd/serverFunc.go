package cmd

import (
	"embed"
	"fmt"
	"net/http"

	//"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	listenPort = "1302"
)

//go:embed templates/*
var templatesFolder embed.FS

func server() {
	r := gin.Default()
	//r.Use(static.Serve("/", static.EmbedFolder(templatesFolder, "templates")))
	r.GET("/", returnIndex)
	fmt.Printf("Listening on port %v...\n", listenPort)
	r.Run(":" + listenPort)
}

func returnIndex(c *gin.Context) {
	c.HTML(http.StatusOK, templatesFolder, gin.H{
		"UserAgent": c.Request.UserAgent(),
	})
}
