package utils

import (
	"log"
	"os"
)

// Must is a helper function that takes a value of any type and an error.
// If the error is nil, it returns the value; if the error is non-nil, it panics.
func Must[T any](value T, err error) T {
	if err != nil {
		log.Fatalln(err)
	}
	return value
}

func Expect(result bool, msg string) {
	if !result {
		log.Fatalln(msg)
	}
}

func Assert(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func AssertFileExists(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("File %s does not exist\n", filename)
	}
}
