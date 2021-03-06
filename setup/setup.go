package setup

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"gopkg.in/ini.v1"
)

// FindSupportedUtilities - Determine supported binaries and path
func FindSupportedUtilities() map[string]string {

	utilities := []string{"iostat", "top", "mpstat", "vmstat", "ss", "nstat", "cat", "ps", "nfsiostat", "ethtool", "ip", "pidstat", "rtmon"}
	u := make(map[string]string)

	for _, utility := range utilities {
		path, err := exec.LookPath(utility)
		if err != nil {
			fmt.Printf("Cannot find %s. Excluding.\n", utility)
			u[utility] = ""
		} else {
			fmt.Printf("%s is supported. Path is %s.\n", utility, path)
			u[utility] = path
		}
	}
	return u
}

// CreateOrLoadConfig - Create a configuration file if not already present
func CreateOrLoadConfig() int {

	const configFile = "gdg.cfg"
	const interval = "30"
	const logDir = "/Users/rich/Downloads/gdg/"

	// check that logDir and configFile are present, if not create
	file, err := os.OpenFile(configFile, os.O_RDONLY|os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		log.Print("Error: ", err)
		return 1
	}
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0750)
	}
	if err != nil {
		log.Print("Error: ", err)
		return 1
	}

	utilities := FindSupportedUtilities()
	cfg, err := ini.Load(configFile)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	cfg.Section("").NewKey("interval", interval)
	cfg.Section("").NewKey("logdir", logDir)
	for u, p := range utilities {
		cfg.Section("utility").NewKey(u, p)
	}

	cfg.SaveTo(configFile)
	return 0
}

// Setup systemd timer
