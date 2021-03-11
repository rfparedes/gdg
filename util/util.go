package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

// Constants for file/direction locations following FHS 3.0 Specification
const (
	ConfigFile = "/etc/gdg.cfg"
	DataDir    = "/var/log/gdg-data/"
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
		fmt.Println("Cannot find hostname executable")
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
	return ("\nzzz ***" + t.Format("Mon Jan 2 03:04:05 MST 2006"))
}

// SetConfigKey sets the configuration file key value
func SetConfigKey(key string, value string, section string) (err error) {

	cfg, err := ini.Load(ConfigFile)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	cfg.Section(section).NewKey(key, value)
	err = cfg.SaveTo(ConfigFile)
	return err
}

// GetConfigKeyValue gets the configuration file key value
func GetConfigKeyValue(key string, section string) (value string, err error) {

	cfg, err := ini.Load(ConfigFile)
	if err != nil {
		return
	}
	value = cfg.Section(section).Key(key).String()
	return value, err
}

// DirSizeMB gets the size of the datadir
func DirSizeMB(dir string) (sizeMB float64, err error) {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return 0.0, err
	}
	var dirSize int64 = 0
	readSize := func(dir string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			dirSize += file.Size()
		}
		return nil
	}
	filepath.Walk(dir, readSize)
	sizeMB = float64(dirSize) / 1024.0 / 1024.0
	return sizeMB, err
}

// DStateCount will return the number of processes in D state
func DStateCount() int64 {

	cmd := "ps -eo stat | grep D | wc -l"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", cmd)
		return 0
	}
	s := string(out)
	i, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		fmt.Printf("Failed to convert DState string\n")
		return 0
	}
	return i
}

// GetStatus will get the current status in config
func GetStatus(progName string, ver string) {
	status, err := GetConfigKeyValue("status", "")
	if err != nil {
		fmt.Println("Cannot get status. Try running '-start' if this is first time running")
		os.Exit(1)
	}
	interval, err := GetConfigKeyValue("interval", "")
	if err != nil {
		fmt.Println("Cannot get interval. Try running '-stop', then '-start'")
		os.Exit(1)
	}
	rtmon, err := GetConfigKeyValue("rtmon", "")
	if err != nil {
		fmt.Println("~ Cannot get rtmon status. ~")
		os.Exit(1)
	}
	dstate, err := GetConfigKeyValue("dstate", "d-state")
	if err != nil {
		fmt.Println("~ Cannot get dstate status. ~")
		os.Exit(1)
	}
	numprocs, err := GetConfigKeyValue("numprocs", "d-state")
	if err != nil {
		fmt.Println("~ Cannot get dstate numprocs ~")
		os.Exit(1)
	}
	fmt.Println("~~~~~~~~~~~~~~~")
	fmt.Println("  gdg status")
	fmt.Println("~~~~~~~~~~~~~~~")
	fmt.Printf("VERSION: %s-%s\n", progName, ver)
	fmt.Printf("STATUS: %s\n", status)
	fmt.Printf("RTMON: %s\n", rtmon)
	fmt.Printf("INTERVAL: %ss\n", interval)
	fmt.Printf("DATA LOCATION: %s\n", DataDir)
	fmt.Printf("CONFIG LOCATION: %s\n", ConfigFile)
	dirSize, err := DirSizeMB(DataDir)
	if err != nil {
		fmt.Printf("CURRENT DATA SIZE: N/A\n")
	} else {
		fmt.Printf("CURRENT DATA SIZE: %.0fMB\n", dirSize)
	}
	fmt.Println("~~~~~~~~~~~~~~~")
	fmt.Printf("DSTATE: %s\n", dstate)
	fmt.Printf("NUMPROCS: %s\n", numprocs)
}
