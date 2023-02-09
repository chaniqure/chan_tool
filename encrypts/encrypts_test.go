package encrypts

import (
	"fmt"
	"testing"
)

func TestEncrypts(t *testing.T) {
	password := "123456"
	fmt.Println("md5：", MD5(password))
	fmt.Println("hex：", SHA1(password))
}
