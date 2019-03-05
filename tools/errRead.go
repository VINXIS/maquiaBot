package tools

import "log"

// ErrRead will check to see if there is an error; it will print the error if there is any
func ErrRead(err error) {
	if err != nil {
		log.Fatal("An error occurred: ", err)
		return
	}
}
