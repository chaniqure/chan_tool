package encrypts

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

//AES解密,key 必须是 16 24 32位
func AESDecrypt(src, key []byte) (bytes []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(src))

	if len(src)%blockSize != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}

	blockMode.CryptBlocks(origData, src)
	origData, err = PKCS7UnPadding(origData)

	if err != nil {
		return nil, err
	}

	return origData, nil
}

//去补码
func PKCS7UnPadding(src []byte) (bytes []byte, err error) {
	length := len(src)
	unpadding := int(src[length-1])

	if length-unpadding < 0 || length-unpadding >= length {
		return nil, errors.New("invalid src bytes")
	}

	return src[:length-unpadding], nil
}

//AES加密,key 必须是 16 24 32位
func AESEncrypt(src, key []byte) (bytes []byte, err error) {
	//获取block块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//补码
	src = PKCS7Padding(src, block.BlockSize())

	//加密模式，
	blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])

	//创建明文长度的数组
	crypted := make([]byte, len(src))

	//加密明文
	blockMode.CryptBlocks(crypted, src)

	return crypted, nil
}

//补码
func PKCS7Padding(src []byte, blockSize int) []byte {
	//计算需要补几位数
	padding := blockSize - len(src)%blockSize

	//在切片后面追加char数量的byte(char)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(src, padtext...)
}
