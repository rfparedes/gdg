package util

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gopkg.in/ini.v1"
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

// SetConfigKey sets the configuration file key value
func SetConfigKey(key string, value string, section string) (err error) {
	// Get current working directory to store config file and dataDir
	pwd, err := os.Getwd()
	if err != nil {
		log.Print("Cannot get current working directory")
		os.Exit(1)
	}
	configFile := pwd + "/gdg.cfg"
	cfg, err := ini.Load(configFile)
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	cfg.Section(section).NewKey(key, value)
	err = cfg.SaveTo(configFile)
	return err
}

// GetLocations get the locations of configfile and datadir
func GetLocations() (gdgPath string, configFile string, dataDir string) {
	// Get current working directory to store config file and dataDir
	pwd, err := os.Getwd()
	if err != nil {
		log.Print("Cannot get current working directory")
		os.Exit(1)
	}
	gdgPath = pwd + "/"
	configFile = pwd + "/gdg.cfg"
	dataDir = pwd + "/gdg-data/"
	return gdgPath, configFile, dataDir
}

// GetConfigKeyValue gets the configuration file key value
func GetConfigKeyValue(key string, section string) (value string, err error) {
	// Get current working directory to store config file and dataDir
	pwd, err := os.Getwd()
	if err != nil {
		log.Print("Cannot get current working directory")
		os.Exit(1)
	}
	configFile := pwd + "/gdg.cfg"
	cfg, err := ini.Load(configFile)
	if err != nil {
		return
	}

	value = cfg.Section(section).Key(key).String()
	return value, err
}

// DirSizeMB gets the size of the datadir
func DirSizeMB(path string) float64 {
	var dirSize int64 = 0
	readSize := func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			dirSize += file.Size()
		}
		return nil
	}
	filepath.Walk(path, readSize)
	sizeMB := float64(dirSize) / 1024.0 / 1024.0
	return sizeMB
}
