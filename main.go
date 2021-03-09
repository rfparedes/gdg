package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
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
	flag.BoolVar(&c.gather, "g", false, "Gather oneshot")
	flag.BoolVar(&c.status, "status", false, "Get current status")
}

const progName string = "gdg"
const ver string = "0.9.0"

var c = config{}

func main() {

	// At least one argument is required
	if len(os.Args) <= 1 {
		fmt.Println("Nothing to do.")
		os.Exit(0)
	}
	// Make sure user is running gdg out of /usr/local/sbin/
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	if exPath != setup.GdgDir {
		fmt.Println("gdg binary must be in /usr/local/sbin")
		os.Exit(1)
	}

	//c := config{}
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
			fmt.Println("~ Starting gdg ~")
			start()
			fmt.Println("~ gdg started ~")
		} else {
			fmt.Println("gdg is already started")
		}
		os.Exit(0)
	}

	// User stops gdg
	if c.stop == true {
		status, _ := util.GetConfigKeyValue("status", "")
		if status != "stopped" {
			fmt.Println("~ Stopping gdg ~")
			stop()
			fmt.Println("~ gdg stopped ~")
		} else {
			fmt.Println("gdg is already stopped")
		}
		os.Exit(0)
	}

	// User reloads gdg
	if c.reload == true {
		fmt.Println("~ Reloading gdg ~")
		stop()
		start()
		fmt.Println("~ gdg reloaded ~")
		os.Exit(0)
	}
	if c.gather == true {
		action.Gather()
		os.Exit(0)
	}
}

func start() {
	setup.CreateOrLoadConfig(strconv.Itoa(c.interval))
	setup.CreateSystemd(strconv.Itoa(c.interval))
	setup.EnableSystemd()
}

func stop() {
	setup.DisableSystemd()
	setup.DeleteSystemd()
}
