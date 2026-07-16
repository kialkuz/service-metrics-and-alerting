package sign

import (
	"crypto/hmac"
	"crypto/sha256"
)

func Generate(jsonData []byte, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(jsonData)
	return h.Sum(nil)
}
