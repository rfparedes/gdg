package util

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"time"
)

// CreateDir function
func CreateDir(dirName string) (err error) {

	err = os.MkdirAll(dirName, 0755)
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

// Check error
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// GetShortHostname function
func GetShortHostname() string {

	const defaultHostname = "unknown"

	hostnameBinary, err := exec.LookPath("hostname")
	if err != nil {
		log.Print("Cannot find hostname executable")
		return defaultHostname
	}
	cmd := exec.Command(hostnameBinary, "-s")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return defaultHostname
	}
	hostname := out.String()
	hostname = hostname[:len(hostname)-1]
	return hostname
}

// CurrentDatFile function
func CurrentDatFile(utility string) string {
	const fileExt = ".dat"

	t := time.Now()
	file := utility + "_" + t.Format("06.01.02.15") + "00" + fileExt
	return file
}

// CreateHeader will create date header for .dat files
func CreateHeader() string {
	t := time.Now()
	return ("zzz ***" + t.Format("Mon Jan 2 03:04:05 MST 2006"))
}
