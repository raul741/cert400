/*
Copyright © 2024 raul <raul@bulgariu.xyz>
*/

package cmd

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	listenPort = "1302"
)

//go:embed templates/**
var templateFolder embed.FS

func server() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error happened looking up user home directory: %v\n", err)
	}
	checkFolders(home)
	lPort := viper.GetString("Server.port")
	if lPort != "" {
		listenPort = lPort
	}
	r := gin.Default()
	LoadHTMLFromEmbedFS(r, templateFolder, "templates/*.html")
	r.Static("/css", "./cmd/templates/css")

	r.GET("/", returnIndex)
	r.GET("/cacert", returnCacert)
	r.POST("/api/upload", getDomainRequest)

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

func getDomainRequest(c *gin.Context) {
	serNum, _ := c.GetPostForm("formSerNum")
	serNumInt, err := strconv.Atoi(serNum)
	if err != nil {
		e := fmt.Sprint(err)
		c.String(http.StatusInternalServerError, e)
		return
	}
	Org, _ := c.GetPostForm("formOrg")
	Country, _ := c.GetPostForm("formCountry")
	Province, _ := c.GetPostForm("formProv")
	Locality, _ := c.GetPostForm("formLocal")
	StreetAddr, _ := c.GetPostForm("formStreet")
	PostCode, _ := c.GetPostForm("formPostal")
	ExpiryTime, _ := c.GetPostForm("formAge")
	ExpiryTimeInt, err := strconv.Atoi(ExpiryTime)
	if err != nil {
		e := fmt.Sprint(err)
		c.String(http.StatusInternalServerError, e)
	}
	BitSize, _ := c.GetPostForm("formBit")
	BitSizeInt, err := strconv.Atoi(BitSize)
	if err != nil {
		e := fmt.Sprint(err)
		c.String(http.StatusInternalServerError, e)
		return
	}
	DNSName, _ := c.GetPostForm("formDNS")
	certDownload, keyDownload, err := generateCert(serNumInt, Org,
		Country, Province, Locality, StreetAddr,
		PostCode, DNSName, ExpiryTimeInt, BitSizeInt)
	if err != nil {
		e := fmt.Sprint(err)
		c.String(http.StatusInternalServerError, e)
		return
	}
	c.File(certDownload)
	c.File(keyDownload)
	return
}
