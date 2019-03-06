package tools

import (
	"io"
	"net/http"
	"os"
)

// DownloadFile downloads a file given the file patch and download URL
func DownloadFile(filepath string, url string) {

	// Get the data
	resp, err := http.Get(url)
	ErrRead(err)
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	ErrRead(err)
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	ErrRead(err)
}
