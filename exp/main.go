package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/Users/patrickfurtak/desktop/go-gallery/hash"
)

func main() {
	toHash := []byte("string to hash")
	h := hmac.New(sha256.New, []byte("notsosupersecretkey"))
	h.Write(toHash)
	b := h.Sum(nil)
	fmt.Println(base64.URLEncoding.EncodeToString(b))

	hmac := hash.NewHMAC("notsosupersecretkey")
	fmt.Println(hmac.Hash("string to hash"))
}
