package encrypts

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"errors"
)

//3DESC初始化向量
const encrypt3DESIV = "boolbool"

//PKCS5Padding
func PKCS5Padding(src []byte, blocksize int) []byte {
	padnum := blocksize - len(src)%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	return append(src, pad...)
}

func PKCS5UnPadding(src []byte) (bytes []byte, err error) {
	n := len(src)
	unpadnum := int(src[n-1])

	if n-unpadnum < 0 || n-unpadnum >= n {
		return nil, errors.New("invalid src bytes")
	}

	return src[:n-unpadnum], nil
}

//3des加密
//src:待加密的明文   key:密钥
//返回值：加密之后的密文
func DES3Encrypt(src, key []byte) (bytes []byte, err error) {
	//创建加密块
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	length := block.BlockSize()

	//填充最后一组数据
	src = PKCS5Padding(src, length)

	//初始化向量
	//iv := key[:block.BlockSize()]
	iv := []byte(encrypt3DESIV)

	//创建CBC加密模式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	dst := make([]byte, len(src))

	//加密
	blockMode.CryptBlocks(dst, src)

	return dst, nil
}

//3des解密
//src：待解密的密文  key:密钥，和加密时使用的密钥相同
//返回值:解密之后的明文
func DES3Decrypt(src, key []byte) (bytes []byte, err error) {
	//创建解密的块
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	//初始化向量
	//iv := key[:block.BlockSize()]
	iv := []byte(encrypt3DESIV)

	//创建CBC解密模式
	blockMode := cipher.NewCBCDecrypter(block, iv)
	dst := make([]byte, len(src))

	//解密
	blockMode.CryptBlocks(dst, src)

	//去除尾部填充数据
	dst, err = PKCS5UnPadding(dst)

	if err != nil {
		return nil, err
	}

	return dst, nil
}

//对用户输入的密钥进行梳理，如果用户输入的密钥太长
//保证key是24位
func GenKey(key []byte) []byte {
	//创建切片，用于存储最终的密钥
	kkey := make([]byte, 0, 24)
	length := len(key)

	//密钥长度大于24字节
	if length > 24 {
		kkey = append(kkey, key[:24]...)
	} else {
		div := 24 / length
		mod := 24 % length
		for i := 0; i < div; i++ {
			kkey = append(kkey, key...)
		}
		kkey = append(kkey, key[:mod]...)
	}

	return kkey
}
