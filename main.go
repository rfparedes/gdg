package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"

	"github.com/rfparedes/gdg/action"
	"github.com/rfparedes/gdg/setup"
	"github.com/rfparedes/gdg/util"
)

type config struct {
	version  bool
	stop     bool
	start    bool
	interval int
	gather   bool
	status   bool
	reload   bool
}

func (c *config) setup() {
	flag.BoolVar(&c.version, "v", false, "Output version information")
	flag.BoolVar(&c.start, "start", false, "Start gathering data")
	flag.BoolVar(&c.stop, "stop", false, "Stop gathering data")
	flag.BoolVar(&c.reload, "reload", false, "Reload after interval or utility change")
	flag.IntVar(&c.interval, "t", 30, "Gathering interval in seconds")
	flag.BoolVar(&c.gather, "g", false, "Gather one-time")
	flag.BoolVar(&c.status, "status", false, "Get current status")
}

const progName string = "gdg"
const ver string = "0.9.0"

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("Nothing to do.")
		os.Exit(0)
	}
	c := config{}
	c.setup()
	flag.Parse()

	// User requests version
	if c.version == true {
		fmt.Println(progName + " v" + ver + " (https://github.com/rfparedes/gdg)")
		os.Exit(0)
	}
	// User requests status
	if c.status == true {
		status, err := util.GetConfigKeyValue("status", "")
		if err != nil {
			fmt.Println("Cannot get status. Try running '-start' if this is first time running")
			os.Exit(1)
		}
		interval, err := util.GetConfigKeyValue("interval", "")
		if err != nil {
			fmt.Println("Cannot get interval. Try running '-stop', then '-start'")
			os.Exit(1)
		}
		fmt.Printf("VERSION: %s-%s\n", progName, ver)
		fmt.Printf("STATUS: %s\n", status)
		fmt.Printf("INTERVAL: %ss\n", interval)
		fmt.Printf("DATA LOCATION: %s\n", util.DataDir)
		fmt.Printf("CONFIG LOCATION: %s\n", util.ConfigFile)
		dirSize, err := util.DirSizeMB(util.DataDir)
		if err != nil {
			fmt.Printf("CURRENT DATA SIZE: N/A\n")
		} else {
			fmt.Printf("CURRENT DATA SIZE: %.0fMB\n", dirSize)
		}
		os.Exit(0)
	}
	// Everything but getting version requires root user
	user, _ := user.Current()
	if user.Uid != "0" {
		fmt.Println("NOT RUNNING AS ROOT")
		os.Exit(1)
	}

	// User enters interval less than 30s (NOT ALLOWED)
	if c.interval < 30 {
		log.Println("Interval cannot be less than 30s. Setting to 30s")
		c.interval = 30
	}

	// User starts gdg
	if c.start == true {
		status, _ := util.GetConfigKeyValue("status", "")
		if status != "started" {
			log.Print("Setting up Granular Data Gatherer")
			setup.CreateOrLoadConfig(strconv.Itoa(c.interval))
			log.Print("Creating and Enabling systemd service and timer in /etc/systemd/system/")
			setup.CreateSystemd(strconv.Itoa(c.interval))
			setup.EnableSystemd()
		} else {
			fmt.Println("gdg is already started")
		}
		os.Exit(0)
	}

	// User stops gdg
	if c.stop == true {
		status, _ := util.GetConfigKeyValue("status", "")
		if status != "stopped" {
			log.Print("Stopping Granular Data Gatherer")
			log.Print("Keeping Configuration and Data. Delete manually")
			log.Print("Deleting systemd service and timer in /etc/systemd/system/")
			setup.DisableSystemd()
			setup.DeleteSystemd()
		} else {
			fmt.Println("gdg is already stopped")
		}
		os.Exit(0)
	}

	// User reloads gdg
	if c.reload == true {
		log.Print("Reloading gdg")
		setup.DisableSystemd()
		setup.DeleteSystemd()
		setup.CreateOrLoadConfig(strconv.Itoa(c.interval))
		setup.CreateSystemd(strconv.Itoa(c.interval))
		setup.EnableSystemd()
		os.Exit(0)
	}
	if c.gather == true {
		action.Gather()
		os.Exit(0)
	}
}
