package tools

import "os"

// DeleteFile deletes a file
func DeleteFile(path string) {
	_ = os.Remove(path)
	return
}
