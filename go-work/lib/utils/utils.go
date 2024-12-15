package utils

import (
	"log"
	"os"
)

// Must1 is a helper function that takes a value of any type and an error.
// If the error is nil, it returns the value; if the error is non-nil, it panics.
func Must1[T any](value T, err error) T {
	if err != nil {
		log.Fatalln(err)
	}
	return value
}

func Must2[T any, U any](vT T, vU U, err error) (T, U) {
	if err != nil {
		log.Fatalln(err)
	}
	return vT, vU
}

func Must3[T any, U any, V any](vT T, vU U, vV V, err error) (T, U, V) {
	if err != nil {
		log.Fatalln(err)
	}
	return vT, vU, vV
}

func Expect(condition bool, errorMessage string) {
	if !condition {
		log.Fatalln(errorMessage)
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
