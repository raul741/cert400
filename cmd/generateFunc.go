/*
Copyright © 2024 raul <raul@bulgariu.xyz>
*/

package cmd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/spf13/viper"
)

var (
	serialnumber  int
	organization  string
	country       string
	province      string
	locality      string
	streetaddress string
	postalcode    string
	notafter      int
	RSA_bitsize   int
)

func generateCA() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error happened looking up user home directory: %v\n", err)
	}
	checkFolders(home)
	RSA_bitsize = viper.GetInt("CA.rsa_bitsize")
	serialnumber = viper.GetInt("CA.serial_number")
	organization = viper.GetString("CA.organization")
	country = viper.GetString("CA.country")
	province = viper.GetString("CA.province")
	locality = viper.GetString("CA.locality")
	streetaddress = viper.GetString("CA.street_address")
	postalcode = viper.GetString("CA.postal_code")
	notafter = viper.GetInt("CA.not_after")
	// fmt.Println(serialnumber, organization, country, province, locality, streetaddress, postalcode, notafter)
	// os.Exit(0)

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(int64(serialnumber)),
		Subject: pkix.Name{
			Organization:  []string{organization},
			Country:       []string{country},
			Province:      []string{province},
			Locality:      []string{locality},
			StreetAddress: []string{streetaddress},
			PostalCode:    []string{postalcode},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(notafter, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caPrivKey, err := rsa.GenerateKey(rand.Reader, RSA_bitsize)
	if err != nil {
		log.Fatalf("Error happened generating private key: %v\n", err)
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		log.Fatalf("Error happened creating certificate: %v\n", err)
	}

	//////
	//fmt.Printf("Checking %v/.config/\n", home)
	// _, err = os.Stat(home + "/.config")
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		fmt.Printf("$HOME/.config doesn't exist, creating...\n")
	// 		os.Mkdir(home+"/.config", 0700)
	// 	} else {
	// 		log.Fatalf("Error happened accessing .config: %v", err)
	// 	}
	// }

	//fmt.Printf("Checking %v/.config/cert400\n", home)
	// _, err = os.Stat(home + "/.config/cert400")
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		fmt.Printf("$HOME/.config/cert400 doesn't exist, creating...\n")
	// 		os.Mkdir(home+"/.config/cert400", 0700)
	// 	} else {
	// 		log.Fatalf("Error happened accessing cert400: %v\n", err)
	// 	}
	// }

	fmt.Printf("Generating %v/.config/cert400/ca.key...\n", home)
	// TODO: Check if keys already exist and warn the user about it
	key, err := os.OpenFile(home+"/.config/cert400/ca.key", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Error happened opening file: %v\n", err)
	}
	defer key.Close()

	fmt.Printf("Generating %v/.config/cert400/ca.cert...\n", home)
	cert, err := os.OpenFile(home+"/.config/cert400/ca.crt", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Error happened opening file: %v\n", err)
	}
	defer cert.Close()
	//////

	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	cert.WriteString(string(caPEM.Bytes()))
	key.WriteString(string(caPrivKeyPEM.Bytes()))
}

func generateCert(serNumber int, clientOrg string, clientCountry string,
	clientProvince string, clientLocality string, clientStreetAddr string,
	clientPostCode string, clientDomain string, clientNotAfter int, clientBitSize int) (string, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		e := fmt.Errorf("Error happened looking up user home directory: %v\n", err)
		return "", "", e
	}
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(int64(serNumber)),
		Subject: pkix.Name{
			Organization:  []string{clientOrg},
			Country:       []string{clientCountry},
			Province:      []string{clientProvince},
			Locality:      []string{clientLocality},
			StreetAddress: []string{clientStreetAddr},
			PostalCode:    []string{clientPostCode},
		},
		PermittedDNSDomains: []string{clientDomain},
		NotBefore:           time.Now(),
		NotAfter:            time.Now().AddDate(clientNotAfter, 0, 0),
		SubjectKeyId:        []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:         []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:            x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, clientBitSize)
	if err != nil {
		e := fmt.Errorf("Error happened generating client privkey: %v\n", err)
		return "", "", e
	}

	// Parse private CA certificate /////////////////////////////
	caFile, err := os.ReadFile(home + "/.config/cert400/ca.crt")
	if err != nil {
		e := fmt.Errorf("Error happened reading from CA certificate: %v\n", err)
		return "", "", e
	}
	caPubBlock, _ := pem.Decode(caFile)
	caCert, err := x509.ParseCertificate(caPubBlock.Bytes)
	if err != nil {
		e := fmt.Errorf("Error happened parsing CA certificate: %v\n", err)
		return "", "", e
	}
	/////////////////////////////////////////////////////////////

	// Parse public CA certificate /////////////////////////////
	caPrivFile, err := os.ReadFile(home + "/.config/cert400/ca.key")
	if err != nil {
		e := fmt.Errorf("Error happened reading from CA privkey: %v\n", err)
		return "", "", e
	}
	caPrivBlock, _ := pem.Decode(caPrivFile)
	caPrivKey, err := x509.ParsePKCS1PrivateKey(caPrivBlock.Bytes)
	if err != nil {
		e := fmt.Errorf("Error happened parsing CA privkey: %v\n", err)
		return "", "", e
	}

	// Generate signed client certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		e := fmt.Errorf("Error happened signing certificate: %v\n", err)
		return "", "", e
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	log.Printf("Generating %v/.config/")
	fmt.Printf("Generating %v/.config/cert400/clientCertificates/"+clientDomain+".key...\n", home)
	clientkey, err := os.OpenFile(home+"/.config/cert400/clientCertificates/"+clientDomain+".key", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		e := fmt.Errorf("Error happened opening file: %v\n", err)
		return "", "", e
	}
	defer clientkey.Close()

	fmt.Printf("Generating %v/.config/cert400/clientCertificates/"+clientDomain+".crt...\n", home)
	clientcert, err := os.OpenFile(home+"/.config/cert400/clientCertificates/"+clientDomain+".crt", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		e := fmt.Errorf("Error happened opening file: %v\n", err)
		return "", "", e
	}
	defer clientcert.Close()

	certPath := fmt.Sprintf("%v/.config/cert400/clientCertificates/%v.crt", home, clientDomain)
	keyPath := fmt.Sprintf("%v/.config/cert400/clientCertificates/%v.key", home, clientDomain)

	clientkey.WriteString(string(certPrivKeyPEM.Bytes()))
	clientcert.WriteString(string(certPEM.Bytes()))

	return certPath, keyPath, nil
}

func checkFolders(home string) {
	_, err := os.Stat(home + "/.config")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("$HOME/.config doesn't exist, creating...\n")
			os.Mkdir(home+"/.config", 0700)
		} else {
			log.Fatalf("Error happened accessing .config: %v", err)
		}
	}
	_, err = os.Stat(home + "/.config/cert400")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("$HOME/.config/cert400 doesn't exist, creating...\n")
			os.Mkdir(home+"/.config/cert400", 0700)
		} else {
			log.Fatalf("Error happened accessing cert400: %v\n", err)
		}
	}

	_, err = os.Stat(home + "/.config/cert400/clientCertificates")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("$HOME/.config/cert400/clientCertificates doesn't exist, creating...\n")
			os.Mkdir(home+"/.config/cert400/clientCertificates", 0700)
		} else {
			log.Fatalf("Error happened accessing cert400: %v\n", err)
		}
	}

}
