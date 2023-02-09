package encrypts

import "encoding/hex"

var (
	EmptyString = &hexError{"empty hex string"}
)

type hexError struct {
	msg string
}

func (h *hexError) Error() string {
	return h.msg
}

// Encode encodes bytes as a hex string.
func Encode(bytes []byte) string {
	encode := make([]byte, len(bytes)*2)
	hex.Encode(encode, bytes)

	return string(encode)
}

// Decode hex string as bytes
func Decode(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, EmptyString
	}

	return hex.DecodeString(input[:])
}

//判断是否16进制
func IsHex(str string) (isHex bool) {
	b := []byte(str)

	isHex = true
	for _, v := range b {
		if v >= 48 && v <= 57 || v >= 65 && v <= 70 || v >= 97 && v <= 102 {

		} else {
			isHex = false
			break
		}
	}

	return isHex
}
