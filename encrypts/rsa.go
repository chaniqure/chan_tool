package encrypts

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/wenzhenxi/gorsa"
	"time"
)

const (
	CHAR_SET               = "UTF-8"
	BASE_64_FORMAT         = "UrlSafeNoPadding"
	RSA_ALGORITHM_KEY_TYPE = "PKCS8"
	RSA_ALGORITHM_SIGN     = crypto.SHA256
)

type XRsa struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// 生成密钥对
func CreateKeys(keyLength int) (publicKeyStr, privateKeyStr string, err error) {

	start := time.Now().UnixNano()

	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return publicKeyStr, privateKeyStr, err
	}

	derStream := MarshalPKCS8PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}

	bufferPrivate := new(bytes.Buffer)

	err = pem.Encode(bufferPrivate, block)

	if err != nil {
		return
	}

	privateKeyStr = bufferPrivate.String()

	fmt.Println("私钥：", privateKeyStr)

	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return privateKeyStr, privateKeyStr, err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	bufferPublic := new(bytes.Buffer)

	err = pem.Encode(bufferPublic, block)

	if err != nil {
		return
	}

	publicKeyStr = bufferPublic.String()

	fmt.Println("公钥：", publicKeyStr)

	end := time.Now().UnixNano()

	fmt.Println("耗时", (start-end)/1000000)

	return publicKeyStr, privateKeyStr, nil
}

func NewXRsa(publicKey []byte, privateKey []byte) (*XRsa, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	block, _ = pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pri, ok := priv.(*rsa.PrivateKey)
	if ok {
		return &XRsa{
			publicKey:  pub,
			privateKey: pri,
		}, nil
	} else {
		return nil, errors.New("private key not supported")
	}
}

// 公钥加密
func (r *XRsa) PublicEncrypt(data string) (string, error) {
	partLen := r.publicKey.N.BitLen()/8 - 11
	chunks := split([]byte(data), partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(bytes)
	}

	return base64.RawURLEncoding.EncodeToString(buffer.Bytes()), nil
}

// 私钥解密
func (r *XRsa) PrivateDecrypt(encrypted string) (string, error) {
	partLen := r.publicKey.N.BitLen() / 8
	raw, err := base64.RawURLEncoding.DecodeString(encrypted)
	chunks := split([]byte(raw), partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(decrypted)
	}

	return buffer.String(), err
}

// 数据加签
func (r *XRsa) Sign(data string) (string, error) {
	h := RSA_ALGORITHM_SIGN.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	sign, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, RSA_ALGORITHM_SIGN, hashed)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(sign), err
}

// 数据验签
func (r *XRsa) Verify(data string, sign string) error {
	h := RSA_ALGORITHM_SIGN.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	decodedSign, err := base64.RawURLEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(r.publicKey, RSA_ALGORITHM_SIGN, hashed, decodedSign)
}

func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)

	k, _ := asn1.Marshal(info)
	return k
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

//RSA公钥私钥产生
func GenRsaKey(bits int) (publicKeyStr, privateKeyStr string, err error) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}

	derStream := x509.MarshalPKCS1PrivateKey(privateKey)

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	bufferPrivate := new(bytes.Buffer)

	err = pem.Encode(bufferPrivate, block)

	if err != nil {
		return
	}

	privateKeyStr = bufferPrivate.String()
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	bufferPublic := new(bytes.Buffer)
	err = pem.Encode(bufferPublic, block)
	if err != nil {
		return
	}
	publicKeyStr = bufferPublic.String()
	return

}

// 私钥加密
func PriKeyEncrypt(privateKey, str string) (priEncrypt string, err error) {

	//byte用指针的方式转string类型
	//strPriKey := *(*string)(unsafe.Pointer(&r.privateKey))
	//加密
	priEncrypt, err = gorsa.PriKeyEncrypt(str, privateKey)
	if err != nil {
		return priEncrypt, errors.New("私钥加密失败！")
	}

	return priEncrypt, err
}

// 公钥解密
func PublicDecrypt(publicKey, priEncrypt string) (pubDecrypt string, err error) {

	//byte用指针的方式转string类型
	//strPubKey := *(*string)(unsafe.Pointer(&r.publicKey))

	//解密
	pubDecrypt, err = gorsa.PublicDecrypt(priEncrypt, publicKey)
	if err != nil {
		return pubDecrypt, errors.New("公钥解密失败！")
	}

	return pubDecrypt, err
}
