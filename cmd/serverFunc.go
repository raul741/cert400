/*
Copyright © 2024 raul <raul@bulgariu.xyz>
*/

package cmd

import (
	"embed"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	listenPort             = "1302"
	default_organization   = "Empresa SL"
	default_country        = "ES"
	default_province       = "Valencia"
	default_locality       = "Valencia"
	default_street_address = "C/Fake Address 123"
	default_postal_code    = "80085"
	default_not_after      = 10
	default_rsa_bitsize    = 3072
	default_dns_name       = "www.example.com"
)

//go:embed templates/**
var templateFolder embed.FS

func defaultHTML() {
	vipDefOrg := viper.GetString("Web.default_organization")
	if vipDefOrg != "" {
		default_organization = vipDefOrg
	}
	vipDefCount := viper.GetString("Web.default_country")
	if vipDefCount != "" {
		default_country = vipDefCount
	}
	vipDefProv := viper.GetString("Web.default_province")
	if vipDefProv != "" {
		default_province = vipDefProv
	}
	vipDefLocal := viper.GetString("Web.default_locality")
	if vipDefLocal != "" {
		default_locality = vipDefLocal
	}
	vipDefStreet := viper.GetString("Web.default_street_address")
	if vipDefStreet != "" {
		default_street_address = vipDefStreet
	}
	vipDefPostal := viper.GetString("Web.default_postal_code")
	if vipDefPostal != "" {
		default_postal_code = vipDefPostal
	}
	vipNotAfter := viper.GetInt("Web.default_not_after")
	if vipNotAfter != 0 {
		default_not_after = vipNotAfter
	}
	vipRSA := viper.GetInt("Web.default_rsa_bitsize")
	if vipRSA != 0 {
		default_rsa_bitsize = vipRSA
	}
	vipDNS := viper.GetString("Web.default_dns_name")
	if vipDNS != "" {
		default_dns_name = vipDNS
	}
}

func server() {
	defaultHTML()
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
	r.StaticFileFS("/favicon.ico", "./templates/ico/favicon.ico", http.FS(templateFolder))
	r.StaticFileFS("/favicon-32x32.png", "./templates/ico/favicon-32x32.png", http.FS(templateFolder))
	r.StaticFileFS("/favicon-16x16.png", "./templates/ico/favicon-16x16.png", http.FS(templateFolder))
	r.StaticFileFS("/apple-touch-icon.png", "./templates/ico/apple-touch-icon.png", http.FS(templateFolder))
	r.StaticFileFS("/android-chrome-512x512.png", "./templates/ico/android-chrome-512x512.png", http.FS(templateFolder))
	r.StaticFileFS("/android-chrome-192x192.png", "./templates/ico/android-chrome-192x192.png", http.FS(templateFolder))
	r.StaticFileFS("/site.webmanifest", "./templates/ico/site.webmanifest", http.FS(templateFolder))

	r.GET("/", returnIndex)
	r.GET("/cacert", returnCacert)
	r.GET("/download/:domain", returnCerts)
	r.POST("/api/upload", getDomainRequest)

	fmt.Printf("Listening on port %v...\n", listenPort)
	if err := r.Run(":" + listenPort); err != nil {
		log.Fatalf("Error occurred starting server: %v\n", err)
	}
}

func returnCacert(c *gin.Context) {
	home, err := os.UserHomeDir()
	if err != nil {
		e := fmt.Sprintf("Error happened fetching: %v\n", err)
		c.String(http.StatusInternalServerError, e)
		return
	}
	c.FileAttachment(home+"/.config/cert400/ca.crt", "cacert.cer")
}

func returnIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "templates/index.html", gin.H{
		"default_serial_number":  rand.IntN(20000-1) + 1,
		"default_organization":   default_organization,
		"default_country":        default_country,
		"default_province":       default_province,
		"default_locality":       default_locality,
		"default_street_address": default_street_address,
		"default_postal_code":    default_postal_code,
		"default_not_after":      default_not_after,
		"default_rsa_bitsize":    default_rsa_bitsize,
		"default_dns_name":       default_dns_name,
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
