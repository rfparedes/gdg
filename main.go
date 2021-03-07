package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/rfparedes/gdg/action"
	"github.com/rfparedes/gdg/setup"
)

type config struct {
	version   bool
	uninstall bool
	install   bool
	stop      bool
	start     bool
	interval  int
	gather    bool
}

func (c *config) setup() {
	flag.BoolVar(&c.version, "v", false, "Output version information")
	flag.BoolVar(&c.install, "i", false, "Install Granular Data Gatherer")
	flag.BoolVar(&c.uninstall, "u", false, "Uninstall Granular Data Gatherer")
	flag.BoolVar(&c.start, "s", false, "Start gathering data")
	flag.BoolVar(&c.stop, "p", false, "Stop gathering data")
	flag.IntVar(&c.interval, "t", 30, "Gathering interval")
	flag.BoolVar(&c.gather, "g", false, "Gather one-time")

}

const progName string = "gdg"
const ver string = "0.0.9"

func main() {

	c := config{}
	c.setup()
	flag.Parse()
	if c.version == true {
		fmt.Println(progName + " v" + ver + " (https://github.com/rfparedes/gdg)")
		return
	}
	if c.install == true {
		log.Print("Installing Granular Data Gatherer")
		setup.CreateOrLoadConfig()
		log.Print("Creating and Enabling systemd service and timer in /etc/systemd/system/")
		setup.EnableSystemd(strconv.Itoa(c.interval), "/home/rich/mdata/git/gdg/")
	}

	if c.uninstall == true {
		log.Print("Uninstalling Granular Data Gatherer")
		log.Print("Deleting systemd service and timer in /etc/systemd/system/")
		setup.DisableSystemd()
	}
	if c.gather == true {
		action.Gather()
	}
}
