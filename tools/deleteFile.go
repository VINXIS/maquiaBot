package tools

import "os"

// DeleteFile deletes a file
func DeleteFile(path string) {
	var err = os.Remove(path)
	ErrRead(err)
	return
}
