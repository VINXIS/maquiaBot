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

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	ErrRead(err)
	resp.Body.Close()
	out.Close()
}
