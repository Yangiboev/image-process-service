package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func image(w http.ResponseWriter, req *http.Request) {
	images := dell()
	// Just a simple GET request to the image URL
	// We get back a *Response, and an error
	for _, image := range images {
		fmt.Println(image)
		res, err := http.Get(image)

		if err != nil {
			log.Fatalf("http.Get -> %v", err)
		}

		// We read all the bytes of the image
		// Types: data []byte
		data, err := ioutil.ReadAll(res.Body)

		if err != nil {
			log.Fatalf("ioutil.ReadAll -> %v", err)
		}

		// You have to manually close the body, check docs
		// This is required if you want to use things like
		// Keep-Alive and other HTTP sorcery.
		res.Body.Close()
		// You can now save it to disk or whatever...
		ioutil.WriteFile("web."+format(image), data, 0666)
	}
	log.Println("I saved your image buddy!")
}

func upload() {
	url := "https://storage.bunnycdn.com/storageZoneName/path/fileName"

	req, _ := http.NewRequest("PUT", url, nil)

	req.Header.Add("content-type", "application/octet-stream")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}

// possible requests:
// /fit-in/128x64/filters:format(webp)/https://raw.githubusercontent.com/cshum/imagor/master/testdata/gopher.png
// /fit-in/128x64/filters:format(jpeg)/https://raw.githubusercontent.com/cshum/imagor/master/testdata/gopher.png
// /fit-in/128x64/filters:format(png)/https://raw.githubusercontent.com/cshum/imagor/master/testdata/gopher.png

// https://vediagames.b-cdn.net/1641834775838.jpeg
// https://vediagames.b-cdn.net/download%20(1).jpeg
// https://vediagames.b-cdn.net/download.jpeg
func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/image", image)

	http.ListenAndServe(":8090", nil)
}

func dell() []string {
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
