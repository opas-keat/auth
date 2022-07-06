package util

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"omsoft.com/auth/cmd/config"
)

func SuccessResponse[T any](c *fiber.Ctx, msg string, status int, v T) error {
	c.Status(status).JSON(fiber.Map{
		"msg":   msg,
		"error": nil,
		"data":  v,
	})
	return nil
}

func FailOnError[T any](c *fiber.Ctx, err error, msg string, status int, v T) error {
	if err != nil {
		c.Status(status).JSON(fiber.Map{
			"msg":   msg,
			"error": err.Error(),
			"data":  v,
		})
	}
	return nil
}

//ECC mode decryption
func ECBDecrypt(crypted []byte) ([]byte, error) {
	keyHash := md5.Sum([]byte([]byte(config.KEY)))
	key := keyHash[:]
	if !validKey(key) {
		return nil, fmt.Errorf("the length of the secret key is wrong, the current incoming length is %d", len(key))
	}
	if len(crypted) < 1 {
		return nil, fmt.Errorf("source data length cannot be 0")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(crypted)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("the source data length must be an integer multiple of %d, the current length is %d", block.BlockSize(), len(crypted))
	}
	var dst []byte
	tmpData := make([]byte, block.BlockSize())

	for index := 0; index < len(crypted); index += block.BlockSize() {
		block.Decrypt(tmpData, crypted[index:index+block.BlockSize()])
		dst = append(dst, tmpData...)
	}
	dst, err = PKCS5UnPadding(dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

//ECC mode encryption
func ECBEncrypt(src []byte) ([]byte, error) {
	keyHash := md5.Sum([]byte([]byte(config.KEY)))
	key := keyHash[:]
	if !validKey(key) {
		return nil, fmt.Errorf("the length of the secret key is wrong, the current incoming length is %d", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(src) < 1 {
		return nil, fmt.Errorf("source data length cannot be 0")
	}
	src = PKCS5Padding(src, block.BlockSize())
	if len(src)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("the source data length must be an integer multiple of %d, the current length is %d", block.BlockSize(), len(src))
	}
	var dst []byte
	tmpData := make([]byte, block.BlockSize())
	for index := 0; index < len(src); index += block.BlockSize() {
		block.Encrypt(tmpData, src[index:index+block.BlockSize()])
		dst = append(dst, tmpData...)
	}
	return dst, nil
}

//Pkcs5 filling
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//Remove pkcs5 filling
func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])

	if length < unpadding {
		return nil, fmt.Errorf("invalid unpadding length")
	}
	return origData[:(length - unpadding)], nil
}

//Key length verification
func validKey(key []byte) bool {
	k := len(key)
	switch k {
	default:
		return false
	case 16, 24, 32:
		return true
	}
}
