package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

func main() {
	// crreate a self-signed certificate

	// create a private keyx
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// create a certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		panic(err)
	}

	// save the private key
	privBytes := x509.MarshalPKCS1PrivateKey(priv)

	// save the private key to a file
	// write to private.pem
	privFile, _ := os.Create("private.pem")
	defer privFile.Close()
	privBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}
	pem.Encode(privFile, privBlock)
	fmt.Println("created private.pem")

	// save the certificate to a file
	// write to cert.pem
	certFile, _ := os.Create("cert.pem")
	defer certFile.Close()
	certBlock := &pem.Block{Type: "CERTIFICATE", Bytes: cert}
	pem.Encode(certFile, certBlock)
	fmt.Println("created cert.pem")
}
