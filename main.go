package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

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
	rtmon    bool
	dstate   int
}

func (c *config) setup() {
	flag.BoolVar(&c.version, "v", false, "Output version information")
	flag.BoolVar(&c.start, "start", false, "Start gathering data")
	flag.BoolVar(&c.stop, "stop", false, "Stop gathering data")
	flag.BoolVar(&c.reload, "reload", false, "Reload after interval or utility change")
	flag.IntVar(&c.interval, "t", 30, "Gathering interval in seconds")
	flag.BoolVar(&c.gather, "g", false, "Gather oneshot")
	flag.BoolVar(&c.status, "status", false, "Get current status")
	flag.BoolVar(&c.rtmon, "rtmon", false, "Toggle rtmon")
	flag.IntVar(&c.dstate, "d", 1, "Trigger sysrq-t on this many D-state procs")

}

const (
	progName string = "gdg"
	ver      string = "0.9.0"
	gdgDir   string = "/usr/local/sbin"
)

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
	if exPath != gdgDir {
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
		util.GetStatus(progName, ver)
		os.Exit(0)
	}
	// Everything but getting version requires root user
	user, _ := user.Current()
	if user.Uid != "0" {
		fmt.Println("NOT RUNNING AS ROOT")
		os.Exit(1)
	}

	// User enters interval less than 30s (NOT ALLOWED)
	if c.interval < 30 || c.interval > 3600 {
		fmt.Println("! Interval cannot be less than 30s or more than 3600s. Setting to 30s !")
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

	if c.rtmon == true {
		rtmon, err := util.GetConfigKeyValue("rtmon", "")
		if err != nil {
			fmt.Println("~ Cannot get rtmon status ~")
			os.Exit(1)
		}
		if rtmon == "stopped" {
			setup.EnableRtmon()
		} else if rtmon == "started" {
			setup.DisableRtmon()
		} else {
			fmt.Println("~ Cannot determine rtmon status ~")
			os.Exit(1)
		}
		os.Exit(0)
	}

	if c.dstate >= 1 {
		util.SetConfigKey("numprocs", strconv.Itoa(c.dstate), "d-state")
		util.SetConfigKey("dstate", "started", "d-state")
	} else {
		fmt.Println("Number of procs has to be greater than 1")
		os.Exit(1)
	}

	dstate, err := util.GetConfigKeyValue("dstate", "d-state")
	if err != nil {
		fmt.Println("~ Cannot get dstate status ~")
		return
	}
	if dstate == "started" {
		numprocs, err := util.GetConfigKeyValue("numprocs", "d-state")
		if err != nil {
			fmt.Println("~ Cannot get d-state numprocs ~")
			return
		}
		procs, err := strconv.ParseInt(strings.TrimSpace(numprocs), 10, 64)
		dprocs := util.DStateCount()
		if dprocs >= procs {
			action.TriggerSysrq()
		}
	}

	if c.gather == true {
		action.Gather()
		os.Exit(0)
	}
}

func start() {

	gdgTimer := `[Unit]
Description=Granular Data Gatherer Timer
Requires=gdg.service
	
[Timer]
OnActiveSec=0
OnUnitActiveSec=` + strconv.Itoa(c.interval) + "\n" +
		`AccuracySec=500msec
	
[Install]
WantedBy=timers.target`

	gdgService := `[Unit]
Description=Granular Data Gatherer
Wants=gdg.timer
	
[Service]
Type=oneshot
ExecStart=` + gdgDir + "/gdg -g" + "\n" +
		`
[Install]
WantedBy=multi-user.target`

	setup.CreateOrLoadConfig(strconv.Itoa(c.interval))
	setup.CreateSystemd("service", gdgService, "gdg")
	setup.CreateSystemd("timer", gdgTimer, "gdg")
	setup.EnableSystemd("gdg.timer", "status")
}

func stop() {
	setup.DisableSystemd("gdg.timer")
	setup.DeleteSystemd("gdg.timer", "status")
	setup.DeleteSystemd("gdg.service", "status")
}
