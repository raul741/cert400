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
)

func generateCA() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error happened looking up user home directory: %v\n", err)
	}
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2024),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco sucks"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 3072)
	if err != nil {
		log.Fatalf("Error happened generating private key: %v\n", err)
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		log.Fatalf("Error happened creating certificate: %v\n", err)
	}

	//////
	//fmt.Printf("Checking %v/.config/\n", home)
	_, err = os.Stat(home + "/.config")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("$HOME/.config doesn't exist, creating...\n")
			os.Mkdir(home+"/.config", 0700)
		} else {
			log.Fatalf("Error happened accessing .config: %v", err)
		}
	}

	//fmt.Printf("Checking %v/.config/cert400\n", home)
	_, err = os.Stat(home + "/.config/cert400")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("$HOME/.config/cert400 doesn't exist, creating...\n")
			os.Mkdir(home+"/.config/cert400", 0700)
		} else {
			log.Fatalf("Error happened accessing cert400: %v\n", err)
		}
	}

	fmt.Printf("Generating %v/.config/cert400/server.key...\n", home)
	// TODO: Check if keys already exist and warn the user about it
	key, err := os.OpenFile(home+"/.config/cert400/server.key", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Error happened opening file: %v\n", err)
	}
	defer key.Close()

	fmt.Printf("Generating %v/.config/cert400/server.cert...\n", home)
	cert, err := os.OpenFile(home+"/.config/cert400/server.crt", os.O_WRONLY|os.O_CREATE, 0600)
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
	// readCert, err := io.ReadAll(caPEM)
	// if err != nil {
	// 	log.Fatalf("Error happened preparing to write cert: %v\n", err)
	// }
	cert.WriteString(string(caPEM.Bytes()))

	// readKey, err := io.ReadAll(caPrivKeyPEM)
	// if err != nil {
	// 	log.Fatalf("Error happened preparing to write key: %v\n", err)
	// }
	key.WriteString(string(caPrivKeyPEM.Bytes()))

}

func generateCert() {
	// cert := &x509.Certificate{
	// 	SerialNumber: big.NewInt(1658),
	// 	Subject: pkix.Name{
	// 		Organization:  []string{"LOL Company"},
	// 		Country:       []string{"US"},
	// 		Province:      []string{""},
	// 		Locality:      []string{"San Francisco"},
	// 		StreetAddress: []string{"Golden Gate Bridge"},
	// 		PostalCode:    []string{"94016"},
	// 	},
	// 	PermittedDNSDomains: []string{"test.bulgariu.xyz"},
	// 	NotBefore:           time.Now(),
	// 	NotAfter:            time.Now().AddDate(10, 0, 0),
	// 	SubjectKeyId:        []byte{1, 2, 3, 4, 6},
	// 	ExtKeyUsage:         []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	// 	KeyUsage:            x509.KeyUsageDigitalSignature,
	// }

}
