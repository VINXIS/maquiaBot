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
	ErrRead(err, "13", "downloader.go")
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	ErrRead(err, "18", "downloader.go")
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	ErrRead(err, "23", "downloader.go")
}
