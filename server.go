package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func generateURL() []string {
	images := []string{
		"fit-in/128x64/filters:format(webp)/https://raw.githubusercontent.com/cshum/imagor/master/testdata/gopher.png",
		"fit-in/128x64/filters:format(jpeg)/https://vediagames.b-cdn.net/download%20(1).jpeg",
		"fit-in/128x64/filters:format(png)/https://vediagames.b-cdn.net/download.jpeg",
	}
	var res = make([]string, 0, 3)
	for _, image := range images {
		res = append(res, sign(image, "mysecret"))
	}
	return res
}

func processImage() {
	for index, image := range generateURL() {
		res, err := http.Get(image)

		if err != nil {
			log.Fatalf("http.Get -> %v", err)
		}

		upload(fmt.Sprintf("%d.test.%s", index, format(image)), res.Body)
	}
	log.Println("I saved your image buddy!")
}

func sign(path, secret string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(path))
	s := base64.StdEncoding.EncodeToString(hash.Sum(nil))[:40]
	s = strings.Replace(s, "+", "-", -1)
	s = strings.Replace(s, "/", "_", -1)
	return "http://localhost:8000/" + s + "/" + path
}

func format(str string) string {
	var buffer bytes.Buffer
	var found bool
	for i := 0; i < len(str); i++ {
		if str[i] == '(' {
			found = true
		} else if str[i] == ')' {
			return buffer.String()
		} else if found {
			buffer.WriteByte(str[i])
		}
	}
	return ""
}

func upload(fileName string, body io.Reader) {
	url := "https://storage.bunnycdn.com/vediagames/vediagames/" + fileName

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("AccessKey", "88800820-6c69-4dba-9a3919b168f5-dbc5-4262")
	req.Header.Add("content-type", "application/octet-stream")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	response, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(res)
	fmt.Println(string(response))
}

func main() {
	start := time.Now()
	defer func() {
		fmt.Printf("took %v\n", time.Since(start))
	}()

	processImage()
}
