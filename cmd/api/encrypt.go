package api

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func GenerateRSAKeys() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	privateKeyFile, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	privateKeyFile.Write(privateKeyPEM)

	publicKey := &privateKey.PublicKey
	publicKeyFile, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	defer publicKeyFile.Close()

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(publicKey)})
	publicKeyFile.Write(publicKeyPEM)

	return nil
}
