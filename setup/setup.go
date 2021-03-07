package setup

import (
	"github.com/rfparedes/gdg/util"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
)

// FindSupportedUtilities - Determine supported binaries and path
func FindSupportedUtilities() map[string]string {

	utilities := []string{"iostat", "top", "mpstat", "vmstat", "ss", "nstat", "ps", "nfsiostat", "ethtool", "ip", "pidstat", "meminfo", "slabinfo"}
	u := make(map[string]string)

	log.Print("Finding Supported Utilities", utilities)
	for _, utility := range utilities {

		var path string
		var err error

		if utility == "meminfo" || utility == "slabinfo" {
			path, err = exec.LookPath("cat")
		} else {
			path, err = exec.LookPath(utility)
		}
		if err != nil {
			log.Printf("Cannot find %s. Excluding\n", utility)
		} else {
			u[utility] = path
		}
	}
	return u
}

// CreateOrLoadConfig - Create configuration file and directories
func CreateOrLoadConfig(interval string) int {

	argMap := map[string]string{
		"iostat":    " 1 3 -t -k -x -N",
		"top":       " -c -b -n 1",
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
	configFile, dataDir := util.GetLocations()
	nics := getNICs()
	// Create gdg configuration file
	if err := util.CreateFile(configFile); err != nil {
		log.Println("File creation failed with error: " + err.Error())
		os.Exit(1)
	}
	// Create parent log directory
	if err := util.CreateDir(dataDir); err != nil {
		log.Println("Directory creation failed with error: " + err.Error())
		os.Exit(1)
	}

	utilities := FindSupportedUtilities()

	err := util.SetConfigKey("hostname", util.GetShortHostname(), "")
	if err != nil {
		log.Print("Cannot set key 'hostname'")
	}
	err = util.SetConfigKey("interval", interval, "")
	if err != nil {
		log.Print("Cannot set key 'interval'")
	}
	err = util.SetConfigKey("configfile", configFile, "")
	if err != nil {
		log.Print("Cannot set key 'configfile'")
	}
	err = util.SetConfigKey("datadir", dataDir, "")
	if err != nil {
		log.Print("Cannot set key 'datadir'")
	}
	for u, p := range utilities {
		var call string

		//Create child log directory for utility
		if err := util.CreateDir(dataDir + u); err != nil {
			log.Println("Directory creation failed with error: " + err.Error())
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
				//ethtoolCmd := (call + n)
				err = util.SetConfigKey(u+strconv.Itoa(i), call+n, "utility")
				if err != nil {
					log.Print("Cannot set key ", u)
				}
			}
			continue
		}

		err = util.SetConfigKey(u, call, "utility")
		if err != nil {
			log.Print("Cannot set key ", u)
		}
	}
	return 0
}

// CreateSystemd will create service and timer files
func CreateSystemd(interval string, gdgPath string) {

	configfile, _ := util.GetLocations()
	timer := `[Unit]
Description=Granular Data Gatherer Timer
Requires=gdg.service
	
[Timer]
OnActiveSec=0
OnUnitActiveSec=` + interval + "\n" +
		`AccuracySec=500msec
	
[Install]
WantedBy=timers.target`

	service := `[Unit]
Description=Granular Data Gatherer
Wants=gdg.timer
	
[Service]
Type=oneshot
ExecStart=` + gdgPath + "gdg -g -c " + configfile + "\n" +
		`
[Install]
WantedBy=multi-user.target`

	strings := []string{"timer", "service"}
	// Create systemd files
	for _, s := range strings {
		f, err := os.OpenFile("/etc/systemd/system/gdg."+s, os.O_RDWR|os.O_CREATE, 0755)
		util.Check(err)
		defer f.Close()
		if s == "timer" {
			_, err := f.WriteString(timer)
			util.Check(err)
		} else {
			_, err := f.WriteString(service)
			util.Check(err)
		}
		f.Sync()
	}
}

// EnableSystemd enables the systemd gdg.timer
func EnableSystemd() {
	systemctl, err := exec.LookPath("systemctl")
	if err != nil {
		log.Print("Cannot find 'systemctl' executable")
		os.Exit(1)
	}
	enableCmd := exec.Command(systemctl, "enable", "gdg.timer", "--now")
	err = enableCmd.Run()
	if err != nil {
		log.Print("Cannot enable 'gdg.timer'")
		os.Exit(1)
	}
	err = util.SetConfigKey("status", "started", "")
	if err != nil {
		log.Print("Cannot set key 'status'")
	}
}

// DisableSystemd disables the sytemd gdg.timer
func DisableSystemd() {

	systemctl, err := exec.LookPath("systemctl")
	if err != nil {
		log.Print("Cannot find 'systemctl' executable")
		os.Exit(1)
	}
	disableCmd := exec.Command(systemctl, "disable", "gdg.timer", "--now")
	err = disableCmd.Run()
	if err != nil {
		log.Print("Cannot disable 'gdg.timer'")
	}
	err = util.SetConfigKey("status", "stopped", "")
	if err != nil {
		log.Print("Cannot set key 'status'")
	}

}

// DeleteSystemd function to delete the gdg systemd services
func DeleteSystemd() {

	strings := []string{"timer", "service"}
	for _, s := range strings {
		err := os.Remove("/etc/systemd/system/gdg." + s)
		if err != nil {
			log.Print("Cannot remove '/etc/systemd/system/gdg." + s + "'")
		}
	}
	err := util.SetConfigKey("status", "stopped", "")
	if err != nil {
		log.Print("Cannot set key 'status'")
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
