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
	flag.IntVar(&c.dstate, "d", 0, "Trigger sysrq-t on this many D-state procs")

}

const (
	progName string = "gdg"
	ver      string = "0.9.0"
	gdgDir   string = "/usr/local/sbin"
)

var c = config{}

func main() {

	// Make sure user is running gdg out of /usr/local/sbin/
	ex, err := os.Executable()
	if err != nil {
		return
	}
	exPath := filepath.Dir(ex)
	if exPath != gdgDir {
		fmt.Println("gdg binary must be in /usr/local/sbin")
		return
	}

	c.setup()
	flag.Parse()

	// User requests version
	if c.version == true {
		fmt.Println(progName + " v" + ver + " (https://github.com/rfparedes/gdg)")
		return
	}
	// User requests status
	if c.status == true {
		util.GetStatus(progName, ver)
		return
	}
	// Everything but getting version and status requires root user
	user, _ := user.Current()
	if user.Uid != "0" {
		fmt.Println("NOT RUNNING AS ROOT")
		return
	}

	// Only allow -t when start or reload is accompanying -t
	if isFlagPassed("t") == true && (isFlagPassed("start") == false && isFlagPassed("reload") == false) {
		fmt.Println("Use -start or -reload when setting interval")
		return
	}
	// User enters interval less than 30s (NOT ALLOWED)
	if c.interval < 30 {
		fmt.Println("~ Interval cannot be less than 30s ~")
		return
	} else if c.interval > 3600 {
		fmt.Println("~ Interval cannot be more than 3600s ~")
		return
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
		return
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
		return
	}

	// User reloads gdg
	if c.reload == true {
		fmt.Println("~ Reloading gdg ~")
		stop()
		start()
		fmt.Println("~ gdg reloaded ~")
		return
	}

	if c.rtmon == true {

		var rtmon string
		// If rtmon is enabled before gdg started for first time
		if _, err := os.Stat(util.ConfigFile); os.IsNotExist(err) {
			setup.CreateOrLoadConfig(strconv.Itoa(c.interval))
			rtmon = "stopped"
		} else {
			rtmon, err = util.GetConfigKeyValue("rtmon", "")
		}
		if err != nil {
			fmt.Println("~ Cannot determine rtmon status ~")
			return
		}
		if rtmon == "stopped" {
			setup.EnableRtmon()
		} else if rtmon == "started" {
			setup.DisableRtmon()
		} else {
			fmt.Println("~ Cannot determine rtmon status ~")
			return
		}
		return
	}

	if c.dstate >= 1 {
		// Only start dstate if gdg is started
		status, err := util.GetConfigKeyValue("status", "")
		if err != nil {
			fmt.Println("~ Cannot get gdg status. ~")
			os.Exit(1)
		}
		if status == "started" {
			util.SetConfigKey("numprocs", strconv.Itoa(c.dstate), "d-state")
			util.SetConfigKey("dstate", "started", "d-state")
			fmt.Printf("Enabling sysrq-t when D-state procs = %d\n", c.dstate)
			return
		}
	} else if c.dstate < 0 {
		fmt.Println("Number of procs has to be greater than 0")
		return
	}

	if c.gather == true {
		action.Gather()

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
		return
	}

	// print usage
	flag.PrintDefaults()

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
	//gdg stopping reset d-state
	fmt.Println("~ Disabling dstate ~")
	util.SetConfigKey("numprocs", "0", "d-state")
	util.SetConfigKey("dstate", "stopped", "d-state")
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
