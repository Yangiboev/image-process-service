package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

func sign(path string, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(path))
	encodedHash := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	truncatedHash := encodedHash[:20]
	replacedHash := strings.ReplaceAll(truncatedHash, "+", "-")
	replacedHash = strings.ReplaceAll(replacedHash, "/", "_")
	return replacedHash + "/" + path
}

// possible requests:
// /fit-in/128x64/filters:format(webp)/https://raw.githubusercontent.com/cshum/imagor/master/testdata/gopher.png
// /fit-in/128x64/filters:format(jpeg)/https://raw.githubusercontent.com/cshum/imagor/master/testdata/gopher.png
// /fit-in/128x64/filters:format(png)/https://raw.githubusercontent.com/cshum/imagor/master/testdata/gopher.png
func main() {
	signedPath := sign("/fit-in/128x64/filters:format(webp)/https://raw.githubusercontent.com/cshum/imagor/master/testdata/gopher.png", "mysecret")
	fmt.Println(signedPath)
}
