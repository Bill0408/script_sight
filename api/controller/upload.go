package controller

import (
	"bytes"
	"context"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	_ "image/jpeg" // import for JPEG decoding
	_ "image/png"  // import for PNG decoding
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	// maxMemory is the maximum amount of data that I have decided that can be received from the multipart form.
	// 10 << 20 is a bit wise operator that basically shifts the bits of the number to the left by the
	// number of positions specified on the right side. Left shifting is multiplying the left operand
	// by 2 raised to the right operand, so 10 << 20 is equivalent to 10 * 2^20, which is 10,485,760.
	maxMemory = 10 << 20
	// key corresponds to the name attribute of the file input field in the form.
	key = "file"

	// fnImgPairKey is the key used to store and retrieve the slice containing the file name and image data.
	fnImgPairKey = 0
)

// ImageValidator ensures that the request is a multipart form, retrieves the file associated
// with the name field of the multipart form, and checks if the file is png or jpeg.
func ImageValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			if err == http.ErrNotMultipart {
				http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
				return
			}

			http.Error(w, "Error parsing form", http.StatusInternalServerError)
			return
		}

		// Retrieve the first file associated with the key.
		file, handler, err := r.FormFile(key)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			return
		}

		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				http.Error(w, "Error closing the file", http.StatusInternalServerError)
				return
			}
		}(file)

		// Create a slice of bytes to stores the first 4 bytes of the file.
		buffer := make([]byte, 4)

		// Read the first 4 bytes of the file into the slice of bytes.
		_, err = file.Read(buffer)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		// Reset the file pointer back to the very beginning of the file,
		// so that other read/write operations cannot return unexpected data.
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, "Error resetting the read pointer", http.StatusInternalServerError)
			return
		}

		// Check the bytes that were written to the byte slice "buffer" for the magic numbers of jpeg
		// and png files. jpeg files begin with 0xFF 0xD8, and png files begin with 0x89 0x50 0x4E 0x47.
		// If the file does not begin with those values, it not a png or jpeg file.
		if !(buffer[0] == 0xFF && buffer[1] == 0xD8) && // Check if jpeg.
			!(buffer[0] == 0x89 && buffer[1] == 0x50 && buffer[2] == 0x4E && buffer[3] == 0x47) { // Check if png.
			http.Error(w, "Only png and jpeg files are allowed", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), fnImgPairKey, []any{handler.Filename, file})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ImageModifier resizes the image to 28 * 28, and converts it to grayscale.
// This is because the AI model accepts 28 * 28 input features, and since
// the only requirements is to identify numbers, converting the image to
// grayscale won't have an impact on the AI's predictions, and it also
// helps the AI train faster because grayscale are less complicated than RDB images.
func ImageModifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the slice containing the file name and image from
		// the request's context, and cast the returned value as a []any
		// because r.Context().Value returns any, which doesn't match []any.
		fnAndImg := r.Context().Value(fnImgPairKey).([]any)
		// Retrieve the image file from fnAndImg and cast it as multipart.File
		// because fnAndImg contains values with the "any" type,
		// which cannot be decoded.
		imgFile := fnAndImg[1].(multipart.File)

		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				http.Error(w, "Error closing the file", http.StatusInternalServerError)
				return
			}
		}(imgFile)

		// Decode the image so that it can be modified.
		img, _, err := image.Decode(imgFile)
		if err != nil {
			http.Error(w, "Error decoding the image", http.StatusInternalServerError)
			return
		}

		// Resize the image to 28 * 28, so that it matches the model's input features.
		dstImg128 := imaging.Resize(img, 28, 28, imaging.Lanczos)

		// Grayscale the image so that it is easier for the model to train on.
		modifiedImg := imaging.Grayscale(dstImg128)

		fnAndImg[1] = modifiedImg // Replace the old image in fnAndImg with the new modified image.

		ctx := context.WithValue(r.Context(), fnImgPairKey, fnAndImg)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the slice containing the file name and image from the request's context.
	fnAndImg := r.Context().Value(fnImgPairKey).([]any)
	fn := fnAndImg[0].(string) // Get file name from fnAndImg and cast it as an image.
	// Get modified image from fnAndImg and cast it as *image.NRGBA, so that
	// it can be written to the output file.
	modifiedImg := fnAndImg[1].(*image.NRGBA)

	// Create an output file with the same name as the image's file name.
	outputFile, err := os.Create(fmt.Sprintf("../script_sight/api/img/%s", fn))
	if err != nil {
		http.Error(w, "Error creating output file", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			http.Error(w, "Error closing output file", http.StatusInternalServerError)
			return
		}

	}(outputFile)

	// Get the extension of the image's file name and encode the image
	// with the extension that matches the output file extension.
	fileExtension := filepath.Ext(outputFile.Name())
	if fileExtension == ".jpeg" || fileExtension == ".jpg" {
		err = imaging.Encode(outputFile, modifiedImg, imaging.JPEG)
		if err != nil {
			http.Error(w, "Error encoding image", http.StatusInternalServerError)
			return
		}
	} else if fileExtension == ".png" {
		err = imaging.Encode(outputFile, modifiedImg, imaging.PNG)
		if err != nil {
			http.Error(w, "Error encoding image", http.StatusInternalServerError)
			return
		}
	}

	// Reset the pointer so that other file operations can have access to the full contents of the file.
	_, err = outputFile.Seek(0, 0)
	if err != nil {
		http.Error(w, "Error resetting pointer", http.StatusInternalServerError)
		return
	}

	var b bytes.Buffer

	mw := multipart.NewWriter(&b)

	fw, err := mw.CreateFormFile("image", outputFile.Name())
	if err != nil {
		http.Error(w, "Error creating form file", http.StatusInternalServerError)
		return
	}

	// Copy the contents of the output file into the form writer.
	if _, err := io.Copy(fw, outputFile); err != nil {
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	if err := mw.Close(); err != nil {
		http.Error(w, "Error closing the form writer", http.StatusInternalServerError)
		return
	}

	// Create a multipart form request that contains the modified image
	// to the python server that runs the AI model.
	req, err := http.NewRequest("POST", "http://localhost:8000/ai/", &b)
	if err != nil {
		http.Error(w, "Error creating the request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())

	// Make the request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error making the request", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, "Error closing the response body", http.StatusInternalServerError)
			return
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		_, err = w.Write([]byte("Upload successful"))
		if err != nil {
			http.Error(w, "Error occurred", http.StatusInternalServerError)
			return
		}
	}

	return
}
