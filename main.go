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
	version    bool
	stop       bool
	start      bool
	interval   int
	gather     bool
	configfile string
	status     bool
}

func (c *config) setup() {
	flag.BoolVar(&c.version, "v", false, "Output version information")
	flag.BoolVar(&c.start, "start", false, "Start gathering data")
	flag.BoolVar(&c.stop, "stop", false, "Stop gathering data")
	flag.IntVar(&c.interval, "t", 30, "Gathering interval in seconds")
	flag.BoolVar(&c.gather, "g", false, "Gather one-time")
	flag.StringVar(&c.configfile, "c", "gdg.cfg", "gdg.cfg file location")
	flag.BoolVar(&c.status, "status", false, "Get current status")
}

const progName string = "gdg"
const ver string = "0.9.0"

func main() {

	c := config{}
	c.setup()
	flag.Parse()

	// User requests version
	if c.version == true {
		fmt.Println(progName + " v" + ver + " (https://github.com/rfparedes/gdg)")
		return
	}

	// User requests status
	if c.status == true {
		value, err := util.GetConfigKeyValue("status", "")
		if err != nil {
			fmt.Println("Cannot get status. No status value present")
		}
		fmt.Printf("Status: %s\n", value)
		return
	}
	// Everything but getting version requires root user
	user, _ := user.Current()
	if user.Uid != "0" {
		fmt.Println("NOT RUNNING AS ROOT")
		os.Exit(1)
	}

	if c.interval < 30 {
		log.Println("Interval cannot be less than 30s. Setting to 30s")
		c.interval = 30
	}
	if c.start == true {
		status, _ := util.GetConfigKeyValue("status", "")
		if status != "started" {
			log.Print("Setting up Granular Data Gatherer")
			setup.CreateOrLoadConfig(strconv.Itoa(c.interval))
			log.Print("Creating and Enabling systemd service and timer in /etc/systemd/system/")
			setup.CreateSystemd(strconv.Itoa(c.interval), "/home/rich/mdata/git/gdg/")
			setup.EnableSystemd()
		} else {
			fmt.Println("gdg is already started")
		}
	}
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

	}
	if c.gather == true {
		action.Gather(c.configfile)
	}
}
