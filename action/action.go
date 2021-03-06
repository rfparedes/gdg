package action

import (
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/ini.v1"
)

const configFile = "gdg.cfg"

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
		v := cfg.Section("utility").Key(k).Value()
		fmt.Println(v)
		gatherCmd = exec.Command("bash", "-c", v)
		gatherCmd.Run()
	}

}
