package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func doEncrypt(key string, message []byte) ([]byte, error) {
	if len(key) == 0 {
		return message, nil
	}
	encKey := []byte(key)
	return encrypt(encKey, string(message))
}

func encrypt(key []byte, message string) (encmess []byte, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	encmessS := base64.URLEncoding.EncodeToString(cipherText)
	encmess = []byte(encmessS)
	return
}

func doDecrypt(key string, message []byte) ([]byte, error) {
	if len(key) == 0 {
		return message, nil
	}
	encKey := []byte(key)
	return decrypt(encKey, string(message))
}

func decrypt(key []byte, securemess string) (decodedmess []byte, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess = cipherText
	return
}
