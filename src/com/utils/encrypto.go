package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"
)

var root_key = randStringBytes(16)

func Encrypto(plainText []byte, refer string) []byte {
	// AES-128,key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
	key := randStringBytes(16)
	work_iv := randStringBytes(16)
	password_iv := randStringBytes(16)

	password_key, password_err := aesEncrypt(key, root_key, work_iv)
	if password_err != nil {
		panic(password_err)
	}
	password_key1 := base64.StdEncoding.EncodeToString(password_key)
	work_iv1 := base64.StdEncoding.EncodeToString(work_iv)
	password_iv1 := base64.StdEncoding.EncodeToString(password_iv)
	ioutil.WriteFile("com/utils/"+refer+"_decrypt_dependency.txt", []byte(password_key1+"\r\n"+work_iv1+"\r\n"+password_iv1), 0777)

	result, err := aesEncrypt(plainText, key, password_iv)
	if err != nil {
		panic(err)
	}
	return result
}
func Decrypto(result []byte, refer string) string {
	key, password_iv, err := getWorkKey(root_key, refer)
	origData, err := aesDecrypt(result, key, password_iv)
	if err != nil {
		panic(err)
	}
	return string(origData)
}
func getWorkKey(root_key []byte, refer string) ([]byte, []byte, error) {
	password_key, err := ioutil.ReadFile("com/utils/" + refer + "_decrypt_dependency.txt")
	if err != nil {
		return nil, nil, err
	}
	keywrods := strings.Split(string(password_key), "\r\n")
	password_key1, err := base64.StdEncoding.DecodeString(keywrods[0])
	if err != nil {
		panic(err)
	}
	work_iv1, err := base64.StdEncoding.DecodeString(keywrods[1])
	if err != nil {
		panic(err)
	}
	password_iv, err := base64.StdEncoding.DecodeString(keywrods[2])
	if err != nil {
		panic(err)
	}
	work_key, err := aesDecrypt(password_key1, root_key, work_iv1)
	if err != nil {
		panic(err)
	}

	return work_key, password_iv, nil
}

func aesEncrypt(origData, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pkcs5Padding(origData, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))

	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func aesDecrypt(crypted, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pkcs5UnPadding(origData)
	return origData, nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)

	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func randStringBytes(n int) []byte {
	k := make([]byte, n)
	if _, err := rand.Read(k); err != nil {
		fmt.Printf("rand.Read() error:%v \n", err)
	}
	return k
}
