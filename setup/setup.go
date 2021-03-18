package setup

import (
	"fmt"
	"github.com/rfparedes/gdg/util"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
)

// Utility type storing all utility info
type Utility struct {
	Name      string
	Binary    string
	Path      string
	Arg       string
	Dup       bool
	Supported bool
}

// function to populate initial utility struct
func utilityPopulate() []Utility {
	utilities := []Utility{
		{Name: "iostat", Binary: "iostat", Path: "", Arg: "1 3 -t -k -x -N", Dup: false, Supported: false},
		{Name: "top", Binary: "top", Path: "", Arg: "-c -b -w 512 -n 1", Dup: false, Supported: false},
		{Name: "mpstat", Binary: "mpstat", Path: "", Arg: "1 2 -P ALL", Dup: false, Supported: false},
		{Name: "vmstat", Binary: "vmstat", Path: "", Dup: false, Supported: false},
		{Name: "vmstat-d", Binary: "vmstat", Path: "", Arg: "-d", Dup: true, Supported: false},
		{Name: "ss", Binary: "ss", Path: "", Arg: "-neopa", Dup: false, Supported: false},
		{Name: "nstat", Binary: "nstat", Path: "", Arg: "-asz", Dup: false, Supported: false},
		{Name: "ps", Binary: "ps", Path: "", Arg: "-eo user,pid,ppid,%cpu,%mem,vsz,rss,tty,stat,start,time,wchan:32,args", Dup: false, Supported: false},
		{Name: "nfsiostat", Binary: "nfsiostat", Path: "", Arg: " 1 3", Dup: false, Supported: false},
		{Name: "ip", Binary: "ip", Path: "", Arg: "-s -s addr", Dup: false, Supported: false},
		{Name: "pidstat", Binary: "pidstat", Path: "", Arg: "-udrRsvwt -p ALL", Dup: false, Supported: false},
		{Name: "meminfo", Binary: "cat", Path: "", Arg: "/proc/meminfo", Dup: false, Supported: false},
		{Name: "slabinfo", Binary: "cat", Path: "", Arg: "/proc/slabinfo", Dup: true, Supported: false},
		{Name: "numastat", Binary: "numastat", Path: "", Arg: "-m -n -c", Dup: false, Supported: false},
		{Name: "sar", Binary: "sar", Path: "", Arg: "-A 0", Dup: false, Supported: false},
	}
	nics := getNICs()
	addl := Utility{}
	for i, n := range nics {
		if i == 0 {
			addl = Utility{
				Name: "ethtool-" + n, Binary: "ethtool", Path: "", Arg: "-S" + " " + n, Dup: false, Supported: false,
			}
		} else {
			addl = Utility{
				Name: "ethtool-" + n, Binary: "ethtool", Path: "", Arg: "-S" + " " + n, Dup: true, Supported: false,
			}
		}
		utilities = append(utilities, addl)
	}
	return utilities
}

// FindSupportedUtilities returns supported binaries with path
func FindSupportedUtilities() []Utility {

	utilities := utilityPopulate()
	var supported []string

	fmt.Println("~ Finding supported utilities ~")
	for i, utility := range utilities {

		var path string
		var err error
		path, err = exec.LookPath(utility.Binary)
		if err != nil {
			fmt.Printf("~ %s not found ~\n", utility.Binary)
		} else {
			utilities[i].Supported = true
			utilities[i].Path = path
			if utilities[i].Dup == false {
				supported = append(supported, utility.Binary)
			}
		}
	}
	fmt.Println(supported)
	return utilities
}

// CreateOrLoadConfig - Create configuration file and directories
func CreateOrLoadConfig(interval string) int {

	rtmon := "stopped"

	// Check rtmon status if config exists as this is set separately
	if _, err := os.Stat(util.ConfigFile); err == nil {
		rtmon, err = util.GetConfigKeyValue("rtmon", "")
		if err != nil {
			err = util.SetConfigKey("rtmon", "unknown", "")
			if err != nil {
				fmt.Println("Cannot set key 'rtmon'")
			}
		}
	}
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

	err := util.SetConfigKey("status", "stopped", "")
	if err != nil {
		fmt.Println("Cannot set key 'status'")
	}
	err = util.SetConfigKey("hostname", util.GetShortHostname(), "")
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
	err = util.SetConfigKey("rtmon", rtmon, "")
	if err != nil {
		fmt.Println("Cannot set key 'rtmon'")
	}
	// Set default values for dstate
	err = util.SetConfigKey("dstate", "stopped", "d-state")
	if err != nil {
		fmt.Println("Cannot set key 'dstate'")
	}
	err = util.SetConfigKey("numprocs", "0", "d-state")
	if err != nil {
		fmt.Println("Cannot set key 'numprocs'")
	}

	for _, utility := range utilities {

		//Only perform if utility is supported
		if utility.Supported == true {
			//Create child log directory for utility
			if err := util.CreateDir(util.DataDir + utility.Name); err != nil {
				fmt.Println("Directory creation failed with error: " + err.Error())
				os.Exit(1)
			}
			if utility.Arg != "" {
				err = util.SetConfigKey(utility.Name, utility.Path+" "+utility.Arg, "utility")
			} else {
				err = util.SetConfigKey(utility.Name, utility.Path, "utility")
			}
			if err != nil {
				fmt.Println("Cannot set key ", utility.Name)
			}
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

	unitPath := ("/etc/systemd/system/" + service)

	systemctl, err := exec.LookPath("systemctl")
	if err != nil {
		fmt.Println("! Cannot find 'systemctl' executable !")
		os.Exit(1)
	}
	fmt.Println("~ Disabling systemd timer ~")
	// Check for systemd timer file
	if _, err := os.Stat(unitPath); err == nil {
		disableCmd := exec.Command(systemctl, "disable", service, "--now")
		err = disableCmd.Run()
		if err != nil {
			fmt.Printf("~ Cannot disable '%s' ~\n", service)
		}
	} else {
		fmt.Printf("~ Cannot disable nonexistent service '%s'\n", service)
	}
}

// DeleteSystemd function to delete the gdg systemd service or timer
func DeleteSystemd(service string, key string) {

	unitPath := ("/etc/systemd/system/" + service)

	if _, err := os.Stat(unitPath); err == nil {
		fmt.Printf("~ Removing systemd %s ~\n", service)
		err := os.Remove(unitPath)
		if err != nil {
			fmt.Printf("~ Cannot remove %s ~\n", service)
		}
		err = util.SetConfigKey(key, "stopped", "")
		if err != nil {
			fmt.Printf("~ Cannot set key %s ~\n", key)
		}
	} else {
		fmt.Printf("~ Cannot delete nonexistent service '%s'\n", service)
	}
}

// Find physical network interfaces
func getNICs() []string {

	var NICs []string
	interfaces, _ := net.Interfaces()

	// get list virtual nics
	dir, err := ioutil.ReadDir("/sys/devices/virtual/net/")
	if err != nil {
		fmt.Println("Cannot get list of virtual NICs")
	}

	//check if interface is virtual
	for _, nic := range interfaces {
		isVirtual := false
		for _, vnic := range dir {
			if nic.Name == vnic.Name() {
				//is virtual interface
				isVirtual = true
				break
			}
		}
		if isVirtual == false {
			NICs = append(NICs, nic.Name)
		}
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
	// Create rtmon systemd service
	CreateSystemd("service", rtmonService, "rtmon")
	// Enable rtmon systemd service
	EnableSystemd("rtmon.service", "rtmon")
	// Add rtmon_status to configfile as started/stopped

}

// DisableRtmon will disable rtmon
func DisableRtmon() {
	DisableSystemd("rtmon.service")
	DeleteSystemd("rtmon.service", "rtmon")
}
