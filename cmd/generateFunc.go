package cmd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

func generateCert() {
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
	key, err := os.OpenFile("./server.key", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Error happened opening file: %v\n", err)
	}
	defer key.Close()

	cert, err := os.OpenFile("./server.crt", os.O_WRONLY|os.O_CREATE, 0600)
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
