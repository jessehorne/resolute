package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

func GenerateKeyRSA2048() (*rsa.PrivateKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func EncryptMessage(key *rsa.PublicKey, content string) (string, error) {
	b, err := rsa.EncryptPKCS1v15(rand.Reader, key, []byte(content))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func DecryptMessage(key *rsa.PrivateKey, content string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}
	s, err := rsa.DecryptPKCS1v15(rand.Reader, key, b)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func ParsePublicKey(k string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(k))
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

func PublicKeyToString(pk *rsa.PublicKey) string {
	b := x509.MarshalPKCS1PublicKey(pk)
	key := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: b,
	})

	return string(key)
}
