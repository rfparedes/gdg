package setup

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"

	"github.com/rfparedes/gdg/util"
)

// FindSupportedUtilities returns supported binaries with path
func FindSupportedUtilities() map[string]string {

	utilities := []string{"iostat", "top", "mpstat", "vmstat", "ss", "nstat", "ps", "nfsiostat", "ethtool", "ip", "pidstat", "meminfo", "slabinfo"}
	u := make(map[string]string)
	var supportedUtilities []string
	fmt.Println("~ Finding supported utilities ~")
	for _, utility := range utilities {

		var path string
		var err error

		if utility == "meminfo" || utility == "slabinfo" {
			path, err = exec.LookPath("cat")
		} else {
			path, err = exec.LookPath(utility)
		}
		if err != nil {
			fmt.Printf("~ %s not found ~\n", utility)
		} else {
			supportedUtilities = append(supportedUtilities, utility)
			u[utility] = path
		}
	}
	fmt.Printf("~ Supported utilities %s ~\n", supportedUtilities)
	return u
}

// CreateOrLoadConfig - Create configuration file and directories
func CreateOrLoadConfig(interval string) int {

	argMap := map[string]string{
		"iostat":    " 1 3 -t -k -x -N",
		"top":       " -c -b -w 512 -n 1",
		"mpstat":    " 1 2 -P ALL",
		"vmstat":    " -d",
		"ss":        " -neopa",
		"meminfo":   " /proc/meminfo",
		"slabinfo":  " /proc/slabinfo",
		"ps":        " -eo user,pid,ppid,%cpu,%mem,vsz,rss,tty,stat,start,time,wchan:32,args",
		"nfsiostat": " 1 3",
		"ethtool":   " -S ",
		"ip":        " -s -s addr",
		"pidstat":   "",
		"nstat":     " -asz",
	}
	nics := getNICs()

	fmt.Println("~ Setting up gdg ~")
	// Create gdg configuration file
	if err := util.CreateFile(util.ConfigFile); err != nil {
		fmt.Println("File creation failed with error: " + err.Error())
		os.Exit(1)
	}
	// Create parent log directory
	if err := util.CreateDir(util.DataDir); err != nil {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(1)
	}

	utilities := FindSupportedUtilities()

	err := util.SetConfigKey("hostname", util.GetShortHostname(), "")
	if err != nil {
		fmt.Println("Cannot set key 'hostname'")
	}
	err = util.SetConfigKey("interval", interval, "")
	if err != nil {
		fmt.Println("Cannot set key 'interval'")
	}
	err = util.SetConfigKey("configfile", util.ConfigFile, "")
	if err != nil {
		fmt.Println("Cannot set key 'configfile'")
	}
	err = util.SetConfigKey("datadir", util.DataDir, "")
	if err != nil {
		fmt.Println("Cannot set key 'datadir'")
	}
	err = util.SetConfigKey("rtmon", "stopped", "")
	if err != nil {
		fmt.Println("Cannot set key 'rtmon'")
	}

	for u, p := range utilities {
		var call string

		//Create child log directory for utility
		if err := util.CreateDir(util.DataDir + u); err != nil {
			fmt.Println("Directory creation failed with error: " + err.Error())
			os.Exit(1)
		}

		if _, ok := argMap[u]; ok {
			call = p + argMap[u]
		} else {
			call = p
		}
		if u == "ethtool" {
			for i, n := range nics {
				if n == "lo" {
					continue
				}
				err = util.SetConfigKey(u+strconv.Itoa(i), call+n, "utility")
				if err != nil {
					fmt.Println("Cannot set key ", u)
				}
			}
			continue
		}

		err = util.SetConfigKey(u, call, "utility")
		if err != nil {
			fmt.Println("Cannot set key ", u)
		}
	}
	return 0
}

// CreateSystemd will create service and timer files
func CreateSystemd(systemdType string, unitText string, name string) {

	fullPath := ("/etc/systemd/system/" + name + "." + systemdType)

	fmt.Printf("~ Creating systemd %s ~\n", systemdType)
	// Create systemd files
	f, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0644)
	util.Check(err)
	defer f.Close()

	_, err = f.WriteString(unitText)
	util.Check(err)
	f.Sync()
}

// EnableSystemd enables the systemd gdg.timer
func EnableSystemd(service string, key string) {
	systemctl, err := exec.LookPath("systemctl")
	if err != nil {
		fmt.Println("~ Cannot find 'systemctl' executable ~")
		os.Exit(1)
	}
	fmt.Printf("~ Enabling systemd %s ~\n", service)
	enableCmd := exec.Command(systemctl, "enable", service, "--now")
	err = enableCmd.Run()
	if err != nil {
		fmt.Printf("~ Cannot enable '%s' ~\n", service)
		os.Exit(1)
	}
	if len(key) > 0 {
		err = util.SetConfigKey(key, "started", "")
		if err != nil {
			fmt.Printf("~ Cannot set key '%s' ~\n", key)
		}
	}
}

// DisableSystemd disables the sytemd gdg.timer
func DisableSystemd(service string) {

	systemctl, err := exec.LookPath("systemctl")
	if err != nil {
		fmt.Println("! Cannot find 'systemctl' executable !")
		os.Exit(1)
	}
	fmt.Println("~ Disabling systemd timer ~")
	disableCmd := exec.Command(systemctl, "disable", service, "--now")
	err = disableCmd.Run()
	if err != nil {
		fmt.Printf("! Cannot disable '%s' !\n", service)
	}
}

// DeleteSystemd function to delete the gdg systemd service or timer
func DeleteSystemd(name string, key string) {

	fullPath := ("/etc/systemd/system/" + name)

	fmt.Printf("~ Removing systemd %s ~\n", name)
	err := os.Remove(fullPath)
	if err != nil {
		fmt.Printf("~ Cannot remove %s ~\n", fullPath)
	}
	err = util.SetConfigKey(key, "stopped", "")
	if err != nil {
		fmt.Printf("~ Cannot set key %s ~\n", key)
	}
}

// Find network interfaces
func getNICs() []string {
	var NICs []string
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		NICs = append(NICs, inter.Name)
	}
	return NICs
}

// EnableRtmon will enable rtmon logging
func EnableRtmon() {

	rtmon, err := exec.LookPath("rtmon")
	if err != nil {
		fmt.Println("Cannot find 'rtmon' executable")
		os.Exit(1)
	}

	if err := util.CreateDir(util.DataDir + "rtmon"); err != nil {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(1)
	}

	rtmonService := `[Unit]
	Description="RTNetlink Monitor Daemon" 
		
	[Service]
	ExecStart=` + rtmon + " file " + util.DataDir + "rtmon/rtmon.log" + "\n" +
		`
	[Install]
	WantedBy=multi-user.target`

	// Add rtmon to configfile under utility
	err = util.SetConfigKey("rtmon", "started", "")
	if err != nil {
		fmt.Printf("~ Cannot set key %s ~\n", "rtmon")
	}
	fmt.Println(rtmon)
	// Create rtmon systemd service
	CreateSystemd("service", rtmonService, "rtmon")
	// Enable rtmon systemd service
	EnableSystemd("rtmon.service", "rtmon")
	// Add rtmon_status to configfile as enabled/disabled

}

// DisableRtmon will disable rtmon
func DisableRtmon() {
	DisableSystemd("rtmon.service")
	DeleteSystemd("rtmon.service", "rtmon")
}
