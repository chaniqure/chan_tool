package encrypts

import "encoding/base64"

//加密
func Base64Encrypt(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

//解密
func Base64Decrypt(src string) (string, error) {
	r, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	return string(r), nil
}
