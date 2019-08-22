package tools

import (
	"fmt"
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
	for {
		if err != nil {
			fmt.Print("An error occured trying to create a file: ")
			fmt.Println(err)
			fmt.Println("Trying again...")
			out, err = os.Create(filepath)
		} else if err == nil {
			break
		}
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	ErrRead(err)
}
