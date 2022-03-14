package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"
)

// Generate TLS using go: https://eli.thegreenplace.net/2021/go-https-servers-with-tls/
// generating ssl -> https://devcenter.heroku.com/articles/ssl-certificate-self

func main() {

	// default location which will be served
	cwd, err := os.Getwd()

	// generating cert for server.
	genCert()

	if err != nil {
		log.Fatal("Error wille riding cwd, terminating app.")
	}

	// default port
	port := flag.String("port", "5000", "Port to use for file server.")
	dir := flag.String("dir", cwd, "directory to serve. Default current working dir.")
	certFile := flag.String("certfile", "cert.pem", "certificate PEM file")
	keyFile := flag.String("keyfile", "key.pem", "key PEM file")

	// info for user
	fmt.Printf("Starting file server on PORT: %s with DIR: %s", *port, *dir)

	// meaty part of this program
	handler := http.FileServer(http.Dir(string(*dir)))
	log.Fatal(http.ListenAndServeTLS(":"+*port, *certFile, *keyFile, handler))
}

func genCert() {
	// new privete key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}
	// new cert template
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"My Corp"},
		},
		DNSNames:  []string{"localhost"},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	// saving new cert to file cert.pem
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if pemCert == nil {
		log.Fatal("Failed to encode certificate to PEM")
	}
	if err := os.WriteFile("cert.pem", pemCert, 0644); err != nil {
		log.Fatal(err)
	}
	log.Print("wrote cert.pem\n")

	// saving new private key to file key.pem
	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if pemKey == nil {
		log.Fatal("Failed to encode key to PEM")
	}
	if err := os.WriteFile("key.pem", pemKey, 0600); err != nil {
		log.Fatal(err)
	}
	log.Print("wrote key.pem\n")
}
