package main

import (
	"net/http"

	"github.com/cshum/imagor"
)

func resizeHandler(w http.ResponseWriter, r *http.Request) {
	// Read the image from the request body
	img, err := imagor.Decode(r.Body)
	if err != nil {
		http.Error(w, "Error decoding image", http.StatusBadRequest)
		return
	}

	// Resize the image
	resizedImg, err := imagor.Resize(img, 400, 400)
	if err != nil {
		http.Error(w, "Error resizing image", http.StatusInternalServerError)
		return
	}

	// Write the resized image to the response
	if err := imagor.Encode(w, resizedImg); err != nil {
		http.Error(w, "Error encoding image", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/resize", resizeHandler)
	http.ListenAndServe(":8080", nil)
}
