package action

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rfparedes/gdg/util"
	"gopkg.in/ini.v1"
)

const configFile = "/etc/gdg.cfg"
const logDir = "/var/log/gdg/"

// Gather the data
func Gather() {

	var gatherCmd *exec.Cmd
	cfg, err := ini.Load(configFile)
	if err != nil {
		fmt.Printf("Failed to read config file: %v", err)
		os.Exit(1)
	}
	// Get all supported utilities
	keys := cfg.Section("utility").KeyStrings()

	// Gather for each
	for _, k := range keys {
		// Create dat file if it doesn't exist
		datFile := (logDir + k + "/" + util.CurrentDatFile(k))
		f, err := os.OpenFile(datFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0744)
		util.Check(err)
		defer f.Close()
		v := cfg.Section("utility").Key(k).Value()
		_, err = f.WriteString(util.CreateHeader() + "\n")
		util.Check(err)
		gatherCmd = exec.Command("bash", "-c", v)
		gatherCmd.Stdout = f
		err = gatherCmd.Start()
		util.Check(err)
		gatherCmd.Wait()
	}

}
