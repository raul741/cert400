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
	gin.SetMode(gin.ReleaseMode)
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

	//r.Static("/css", "./cmd/templates/css")
	r.StaticFileFS("/css/style.css", "./templates/css/style.css", http.FS(templateFolder))

	r.GET("/", returnIndex)
	r.GET("/cacert", returnCacert)
	r.GET("/download/:domain", returnCerts)
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
	c.HTML(http.StatusOK, "templates/index.html", gin.H{})
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
	result := fmt.Sprintf(`
	<b>
	<p align="center">
		<a target="_blank" href="download/%v.key">Download privkey</a>
	</p>
	<p align="center">
		<a target="_blank" href="download/%v.crt">Download cert</a>
	</p>
	</b>
		`, DNSName, DNSName)
	c.String(http.StatusOK, result)
	//c.File(keyDownload)

	fmt.Println(certDownload, keyDownload)
	return
}

func returnCerts(c *gin.Context) {
	domain := c.Param("domain")
	base, err := os.UserHomeDir()
	if err != nil {
		e := fmt.Sprintf("Error happened finding user home directory: %v\n", err)
		c.String(http.StatusInternalServerError, e)
		return
	}
	c.FileAttachment(base+"/.config/cert400/clientCertificates/"+domain, domain)
	//c.File(base + "/.config/cert400/clientCertificates/" + domain + ".key")
}
