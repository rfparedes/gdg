package setup

import (
	"fmt"
	"github.com/rfparedes/gdg/util"
	"os"
	"os/exec"

	"gopkg.in/ini.v1"
)

// FindSupportedUtilities - Determine supported binaries and path
func FindSupportedUtilities() map[string]string {

	utilities := []string{"iostat", "top", "mpstat", "vmstat", "ss", "nstat", "cat", "ps", "nfsiostat", "ethtool", "ip", "pidstat", "rtmon", "iofake"}
	u := make(map[string]string)

	for _, utility := range utilities {
		path, err := exec.LookPath(utility)
		if err != nil {
			fmt.Printf("Cannot find %s. Excluding.\n", utility)
		} else {
			fmt.Printf("%s is supported. Path is %s.\n", utility, path)
			u[utility] = path
		}
	}
	return u
}

// CreateOrLoadConfig - Create a configuration file
func CreateOrLoadConfig() int {

	argMap := map[string]string{
		"iostat":    " 1 3 -t -k -x -N",
		"top":       " -c -b -n 1",
		"mpstat":    " 1 2 -P ALL",
		"vmstat":    " -d",
		"ss":        " -neopa",
		"meminfo":   "/proc/meminfo",
		"slabinfo":  "/proc/slabinfo",
		"ps":        " -eo user,pid,ppid,%cpu,%mem,vsz,rss,tty,stat,start,time,wchan:32,args",
		"nfsiostat": " 1 3",
		"ethtool":   " -S",
		"ip":        " -s -s addr",
		"pidstat":   "",
		"nstat":     " -asz",
	}

	const interval = "30"
	const configFile = "gdg.cfg"
	const logDir = "/var/log/gdg/"

	// Create gdg configuration file
	if err := util.CreateFile(configFile); err != nil {
		fmt.Println("File creation failed with error: " + err.Error())
		os.Exit(1)
	}
	// Create parent log directory
	if err := util.CreateDir(logDir); err != nil {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(1)
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
		var call string

		//Create child log directory for utility
		if err := util.CreateDir(logDir + u); err != nil {
			fmt.Println("Directory creation failed with error: " + err.Error())
			os.Exit(1)
		}
		if _, ok := argMap[u]; ok {
			call = p + argMap[u]
		} else {
			call = p
		}
		cfg.Section("utility").NewKey(u, call)
	}

	cfg.SaveTo(configFile)
	return 0
}

// Find network interfaces

// Setup systemd timer
