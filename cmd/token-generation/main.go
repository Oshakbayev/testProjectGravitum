package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateEd25519KeyPairBase64() (string, string, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}
	// Кодируем ключи в Base64
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKey)
	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKey)
	return publicKeyBase64, privateKeyBase64, nil
}

func main() {
	publicKey, privateKey, err := GenerateEd25519KeyPairBase64()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("publicKey: ", publicKey)
	fmt.Println("privateKey: ", privateKey)
}
