package rpctool

import (
	"bytes"
	"encoding/base64"
	"sync"
)

const (
	BASE64_LENGTH = 64
)

var once sync.Once
var gEncryptHelperInstance *EncryptHelper

func GetEncryptHelperInstance() *EncryptHelper {
	once.Do(func() {
		ec := &encoder{
			codeStr: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/",
			verify: func(b byte) bool {
				if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '+' || b == '/' {
					return true
				}
				return false
			},
		}

		gEncryptHelperInstance = &EncryptHelper{
			encoding: base64.NewEncoding(string(ec.getEncoder())),
			ecoder:   ec,
		}
	})

	return gEncryptHelperInstance
}

type EncryptHelper struct {
	encoding *base64.Encoding
	ecoder   *encoder
}

func (e *EncryptHelper) base64Encode(src []byte) string {
	return e.encoding.EncodeToString(src)
}

func (e *EncryptHelper) base64Decode(s string) []byte {
	quotient := len(s) / 4
	remainder := len(s) % 4
	var decLen int
	if remainder != 0 { //调整，使得待解密字符串长度为4的倍数 ，若余数为1直接丢弃，否则补充‘=’
		decLen = quotient*3 + remainder - 1
		switch remainder {
		case 1:
			s = s[:len(s)-1]
		case 2:
			s += "=="
		case 3:
			s += "="
		}
	} else {
		decLen = quotient * 3
	}

	if len(s) == 0 {
		return nil
	}

	decData, err := e.encoding.DecodeString(s)
	if err != nil {
		return nil
	}

	return decData[:decLen]
}

func (e *EncryptHelper) Encrypt(user string, pwd string, rand string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(pwd); i++ {
		if !e.ecoder.contains(pwd[i]) {
			buffer.WriteByte(e.ecoder.get(int(pwd[i]) % BASE64_LENGTH))
		} else {
			buffer.WriteByte(pwd[i])
		}
	}
	//fmt.Println("pwd2:",buffer.String())
	pwd2 := e.base64Decode(buffer.String())
	//fmt.Println("pwd2_2:",pwd2,"\nlen(pwd2_2):",len(pwd2))
	maxlen := len(pwd2)
	if len(user) > maxlen {
		maxlen = len(user)
	}
	if len(rand) > maxlen {
		maxlen = len(rand)
	}

	buffer.Reset()
	for idx := 0; idx < maxlen; idx++ {
		tmp := 0
		if idx < len(user) {
			tmp += int(user[idx])
		}
		if idx < len(pwd2) {
			tmp += int(pwd2[idx])
		}
		if idx < len(rand) {
			tmp += int(rand[idx])
		}
		tmp %= 256
		buffer.WriteByte(byte(tmp))
	}
	//raw := buffer.Bytes()
	//fmt.Println("raw:",raw,"\nlen(raw):",len(raw))
	return e.base64Encode(buffer.Bytes())
}
