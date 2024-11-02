package api

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/brownei/chifunds-api/utils"
	"golang.org/x/crypto/bcrypt"
)

type Key struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func RsaEncrypt(origData []byte) (string, string, error) {
	publicKey, err := os.ReadFile("public.pem")
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", "", errors.New("failed to parse public key PEM")
	}

	pubInterface, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return "", "", fmt.Errorf("public key parsing error: %v", err)
	}

	aesKey := make([]byte, 32) // Use AES-256 (32 bytes key)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return "", "", fmt.Errorf("AES key generation error: %v", err)
	}

	encryptedAesKey, err := rsa.EncryptPKCS1v15(rand.Reader, pubInterface, aesKey)
	if err != nil {
		return "", "", fmt.Errorf("encryption error: %v", err)
	}

	blockCipher, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", "", fmt.Errorf("AES cipher creation error: %v", err)
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return "", "", fmt.Errorf("GCM mode creation error: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", fmt.Errorf("nonce generation error: %v", err)
	}

	encryptedData := gcm.Seal(nonce, nonce, origData, nil)

	return base64.StdEncoding.EncodeToString(encryptedData), base64.StdEncoding.EncodeToString(encryptedAesKey), nil
}

func RsaDecrypt(body string) ([]byte, error) {
	privateKey, err := os.ReadFile("private.pem")
	if err != nil {
		return nil, fmt.Errorf("could not read private key file: %v", err)
	}

	cipherText, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return nil, fmt.Errorf("base64 decode error: %v", err)
	}

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("failed to parse private key PEM")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("private key parsing error: %v", err)
	}

	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, priv, cipherText)
	if err != nil {
		return nil, fmt.Errorf("decryption error: %v", err)
	}

	return decryptedData, nil
}

func (a *application) GetAllKeysHandler(w http.ResponseWriter, r *http.Request) {
	// Load public key from file
	publicKey, err := os.ReadFile("public.pem")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	privateKey, err := os.ReadFile("private.pem")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	keys := Key{
		PublicKey:  string(publicKey),
		PrivateKey: string(privateKey),
	}

	// Write public key as JSON response
	utils.WriteJSON(w, http.StatusAccepted, keys)
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
