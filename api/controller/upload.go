package controller

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "image/jpeg" // import for JPEG decoding
	_ "image/png"  // import for PNG decoding
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const (
	imgKey = 0
)

type ImageData struct {
	ImgUrl string `json:"imgUrl"`
}

// ImgUrlConverter is a middleware that converts the image url data into an actual png file.
func ImgUrlConverter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var imgData ImageData

		err := json.NewDecoder(r.Body).Decode(&imgData)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Split the Data URL to extract base64 data
		parts := strings.Split(imgData.ImgUrl, ";base64,")
		if len(parts) != 2 {
			fmt.Println(err)
			http.Error(w, "Invalid data URL format", http.StatusBadRequest)
			return
		}
		base64Data := parts[1]

		// Decode the base64 data
		decoded, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to decode base64 data", http.StatusBadRequest)
			return
		}

		// At this point, 'decoded' contains the raw bytes of the image.
		imgReader := bytes.NewReader(decoded)

		// Create a new file in write mode
		file, err := os.Create("../script_sight/api/img/img.png")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to create image file", http.StatusInternalServerError)
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(file)

		// Copy the bytes from the reader into the file
		_, err = io.Copy(file, imgReader)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}

		// Reset the file pointer
		_, err = file.Seek(io.SeekStart, 0)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "An Error occurred resetting file pointer", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), imgKey, file)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ImgUploader(w http.ResponseWriter, r *http.Request) {
	file := r.Context().Value(imgKey).(*os.File)

	// Create a buffer and a multipart writer
	var requestBody bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Create a form file and write to it
	fileWriter, err := multiPartWriter.CreateFormFile("uploadFile", file.Name())
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to write to form file", http.StatusInternalServerError)
		return
	}

	_, err = file.Seek(0, 0) // ensure file pointer is at start
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to seek to start of the file", http.StatusInternalServerError)
		return
	}

	// Copy the image file to the form writer
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to copy file to form file", http.StatusInternalServerError)
		return
	}

	// Ensure the writer is closed to write the trailing boundary
	err = multiPartWriter.Close()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to close the multipart writer", http.StatusInternalServerError)
		return
	}

	//Set up the HTTP request to the django server that runs the AI.
	url := "http://localhost:8000/ai/"
	request, err := http.NewRequest(http.MethodPost, url, &requestBody)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to create request", http.StatusInternalServerError)
		return
	}

	request.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}

	if response.StatusCode != http.StatusOK {
		fmt.Println(err)
		http.Error(w, fmt.Sprintf("Server returned non-OK status: %v", response.Status), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write([]byte("Success"))
}
