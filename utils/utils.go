package utils

import (
	"fmt"
	"os"
)

// Panicf formats s, and panics.
func Panicf(s string, i ...interface{}) { panic(fmt.Sprintf(s, i...)) }

// PathExists checks whether a fully qualified path is present on the server's filesystem.
func PathExists(fp string) (bool, error) {
	_, err := os.Stat(fp)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
