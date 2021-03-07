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
)

type config struct {
	version    bool
	uninstall  bool
	install    bool
	stop       bool
	start      bool
	interval   int
	gather     bool
	configfile string
}

func (c *config) setup() {
	flag.BoolVar(&c.version, "v", false, "Output version information")
	flag.BoolVar(&c.install, "i", false, "Install Granular Data Gatherer")
	flag.BoolVar(&c.uninstall, "u", false, "Uninstall Granular Data Gatherer")
	flag.BoolVar(&c.start, "s", false, "Start gathering data")
	flag.BoolVar(&c.stop, "p", false, "Stop gathering data")
	flag.IntVar(&c.interval, "t", 30, "Gathering interval in seconds")
	flag.BoolVar(&c.gather, "g", false, "Gather one-time")
	flag.StringVar(&c.configfile, "c", "gdg.cfg", "gdg.cfg file location")
}

const progName string = "gdg"
const ver string = "0.9.0"

func main() {

	user, _ := user.Current()
	if user.Uid != "0" {
		fmt.Println("NOT RUNNING AS ROOT")
		os.Exit(1)
	}
	c := config{}
	c.setup()
	flag.Parse()
	if c.version == true {
		fmt.Println(progName + " v" + ver + " (https://github.com/rfparedes/gdg)")
		return
	}

	if c.interval < 30 {
		log.Println("Interval cannot be less than 30s. Setting to 30s")
		c.interval = 30
	}
	if c.install == true {
		log.Print("Installing Granular Data Gatherer")
		setup.CreateOrLoadConfig(strconv.Itoa(c.interval))
		log.Print("Creating and Enabling systemd service and timer in /etc/systemd/system/")
		setup.CreateSystemd(strconv.Itoa(c.interval), "/home/rich/mdata/git/gdg/")
		setup.EnableSystemd()
	}
	if c.uninstall == true {
		log.Print("Uninstalling Granular Data Gatherer")
		log.Print("Keeping Configuration and Data. Delete manually")
		log.Print("Deleting systemd service and timer in /etc/systemd/system/")
		setup.DisableSystemd()
		setup.DeleteSystemd()
	}
	if c.gather == true {
		action.Gather(c.configfile)
	}
	if c.stop == true {
		setup.DisableSystemd()
	}
	if c.start == true {
		setup.EnableSystemd()
	}
}
