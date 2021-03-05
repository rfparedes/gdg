package setup

import (
	"fmt"
	"os/exec"

	"github.com/spf13/viper"
)

// FindSupportedBinaries - Determine supported binaries and path
func FindSupportedBinaries() []string {

	binaries := []string{"iostat", "top", "mpstat", "vmstat", "ss", "nstat", "cat", "ps", "nfsiostat", "ethtool", "ip", "pidstat", "rtmon"}
	var binariesFound []string

	for _, binary := range binaries {
		path, err := exec.LookPath(binary)
		if err != nil {
			fmt.Printf("Cannot find %s. Excluding.\n", binary)
		} else {
			fmt.Printf("%s is supported. Path is %s.\n", binary, path)
			binariesFound = append(binariesFound, path)
		}
	}
	return binariesFound
}

// CreateOrLoadConfig - Create a configuration file if not already present
func CreateOrLoadConfig() {
	viper.SetConfigName("gdg.cfg")
	viper.AddConfigPath("/etc/")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			viper.SetDefault("interval", "30")
			viper.SetDefault("logdir", "/var/log/gdg")
			viper.SetDefault("iostat", "")
			viper.SetDefault("top", "")
			viper.SetDefault("mpstat", "")
			viper.SetDefault("vmstat", "")
			viper.SetDefault("ss", "")
			viper.SetDefault("nstat", "")
			viper.SetDefault("cat", "")
			viper.SetDefault("ps", "")
			viper.SetDefault("nfsiostat", "")
			viper.SetDefault("ethtool", "")
			viper.SetDefault("ip", "")
			viper.SetDefault("pidstat", "")
			viper.SetDefault("rtmon", "")
		}
	}
}

// Setup systemd timer

// Store configuration4
