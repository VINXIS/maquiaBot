package tools

import (
	"fmt"
	"log"
)

// ErrRead will check to see if there is an error; it will print the error if there is any
func ErrRead(err error, line, file string) {
	if err != nil {
		log.Fatalln("An error occurred at line " + line + " in " + file + ": " + fmt.Sprint(err))
		return
	}
}
