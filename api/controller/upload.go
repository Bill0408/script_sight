package controller

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

// ImgUrlConverter converts the data url that was sent over from the browser to a png file.
func ImgUrlConverter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var imgData ImageData

		// Decode the json data to ImageData struct.
		err := json.NewDecoder(r.Body).Decode(&imgData)
		if err != nil {
			fmt.Println("Error decoding JSON: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Since data urls have two parts, the mime and the actual base64 encoding,
		// split the string at ";base64," to get the mime type and base64 encoding.
		parts := strings.Split(imgData.ImgUrl, ";base64,")
		if len(parts) != 2 {
			fmt.Println("Invalid data URL format.")
			http.Error(w, "Invalid data URL format", http.StatusBadRequest)
			return
		}

		// Get the bytes represented by the base64 encoding.
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			fmt.Println("Error decoding base64: ", err)
			http.Error(w, "Failed to decode base64 data", http.StatusBadRequest)
			return
		}

		// Get an io reader containing the png bytes from the decoded byte array,
		// and create a file called img.png.
		imgReader := bytes.NewReader(decoded)
		file, err := os.Create("img/img.png")
		if err != nil {
			fmt.Println("Error creating file: ", err)
			http.Error(w, "Failed to create image file", http.StatusInternalServerError)
			return
		}

		// Copy the bytes from the io reader to the file.
		_, err = io.Copy(file, imgReader)
		if err != nil {
			fmt.Println("Error copying to file: ", err)
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}

		// Reset the file pointer.
		_, err = file.Seek(0, 0)
		if err != nil {
			fmt.Println("Error seeking file: ", err)
			http.Error(w, "An error occurred resetting file pointer", http.StatusInternalServerError)
			return
		}

		// Save the png file in the request's context, so that
		// it can be accessed from the ImgUploader handler.
		ctx := context.WithValue(r.Context(), imgKey, file)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ImgUploader uploads a multipart form containing the
// image to the django server that's running the AI.
func ImgUploader(w http.ResponseWriter, r *http.Request) {
	// Retrieve the png file from the request's handler, and cast
	// it as an os file because r.Context().Value returns "any" datatype.
	file, ok := r.Context().Value(imgKey).(*os.File)
	if !ok {
		fmt.Println("File not found in context")
		http.Error(w, "File not found in context", http.StatusInternalServerError)
		return
	}

	// Create a byte buffer for the multipart form writer.
	var requestBody bytes.Buffer

	multiPartWriter := multipart.NewWriter(&requestBody)

	// Create a form file with the field name of "uploadFile".
	// The django server uses the field name to access the image file
	fileWriter, err := multiPartWriter.CreateFormFile("uploadFile", file.Name())
	if err != nil {
		fmt.Println("Error creating form file: ", err)
		http.Error(w, "Unable to write to form file", http.StatusInternalServerError)
		return
	}

	// Copy the contents of the png file to the multipart form file writer.
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		fmt.Println("Error copying to form file: ", err)
		http.Error(w, "Unable to copy file to form file", http.StatusInternalServerError)
		return
	}

	err = multiPartWriter.Close()
	if err != nil {
		fmt.Println("Error closing multiPartWriter: ", err)
		http.Error(w, "Unable to close the multipart writer", http.StatusInternalServerError)
		return
	}

	// The url endpoint for on the django server.
	url := "http://django:8000/ai/"

	// Create a new http request with the body of the byte buffer
	// that contains the multipart form data.
	request, err := http.NewRequest(http.MethodPost, url, &requestBody)
	if err != nil {
		fmt.Println("Error creating new request: ", err)
		http.Error(w, "Unable to create request", http.StatusInternalServerError)
		return
	}

	// Set the content type of the request to multipart form.
	request.Header.Set("Content-Type", multiPartWriter.FormDataContentType())
	client := &http.Client{}            // Create an http client to send the request.
	response, err := client.Do(request) // Send the request to the django server.
	if err != nil {
		fmt.Println("Error during client.Do: ", err)
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}

	// Response from the django is not in the 200 range.
	if response.StatusCode != http.StatusOK {
		fmt.Println("Non-OK status: ", response.Status)
		http.Error(w, fmt.Sprintf("Server returned non-OK status: %v", response.Status), http.StatusInternalServerError)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing the response body", err)
		}
	}(response.Body)

	// Read the data of the response body.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Write the read data from the response body back to the client.
	// This data contains the AI's predicted number.
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		fmt.Println("Error writing response: ", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
