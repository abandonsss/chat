package tool

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func GenerateToken() string {
	r := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, r)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(r)
}
