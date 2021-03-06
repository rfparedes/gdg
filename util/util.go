package util

import (
	"os"
)

// CreateDir function
func CreateDir(dirName string) (err error) {

	err = os.MkdirAll(dirName, 0750)
	if err != nil {
		return err
	}
	return nil
}

// CreateFile function
func CreateFile(fileName string) (err error) {

	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0744)
	defer file.Close()
	if err != nil {
		return err
	}
	return nil
}
