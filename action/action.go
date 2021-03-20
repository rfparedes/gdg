package action

import (
	"compress/gzip"
	"fmt"
	"github.com/rfparedes/gdg/util"
	"gopkg.in/ini.v1"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Gather the data
func Gather() {

	var gatherCmd *exec.Cmd

	cfg, err := ini.Load(util.ConfigFile)
	if err != nil {
		fmt.Printf("Failed to read config file: %v", err)
		os.Exit(1)
	}
	// Get all supported utilities
	keys := cfg.Section("utility").KeyStrings()
	// Gather for each
	for _, k := range keys {
		// Create dat file if it doesn't exist
		datFile := (util.DataDir + k + "/" + util.CurrentDatFile(k))
		f, err := os.OpenFile(datFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0744)
		util.Check(err)
		defer f.Close()
		_, err = f.WriteString(util.CreateHeader() + "\n")
		util.Check(err)
		v := cfg.Section("utility").Key(k).Value()
		gatherCmd = exec.Command("bash", "-c", v)
		gatherCmd.Stdout = f
		err = gatherCmd.Start()
		util.Check(err)
		gatherCmd.Wait()
	}
}

// TriggerSysrq will trigger a task trace
func TriggerSysrq() {

	util.SetConfigKey("dstate", "stopped", "d-state")
	util.SetConfigKey("numprocs", "0", "d-state")

	echo, err := exec.LookPath("echo")
	if err != nil {
		fmt.Print("Cannot find 'echo' executable.")
		return
	}
	trigger := echo + " t > /proc/sysrq-trigger"
	echoCmd := exec.Command("bash", "-c", trigger)
	err = echoCmd.Run()
	if err != nil {
		fmt.Print("Failed to trigger sysrq")
		return
	}

}

// TidyLogs will delete logs outside the range of days user requests to keep and
// then gzips the log
func TidyLogs(logdays int) {

	var files []string
	var dates []string

	// Get all the log filenames
	err := filepath.Walk(util.DataDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match("*.dat", filepath.Base(path)); err != nil {
			return err
		} else if matched {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Get list of dates in a range that will need to be deleted to maintain logdays
	start := time.Now().AddDate(0, 0, 0)
	end := start.AddDate(0, 0, -logdays)
	for rd := rangeDate(end, start); ; {
		date := rd()
		if date.IsZero() {
			break
		}
		dates = append(dates, date.Format("06.01.02"))
	}

	// Used to exclude the files currently in use
	inUse := start.Format("06.01.02.15") + "00"

	// Delete the files not in the range of dates and not the current in-use file
	for _, file := range files {
		// If the filename doesn't have inuse in its name
		if !strings.Contains(file, inUse) {
			if !contains(dates, file) {
				e := os.Remove(file)
				if e != nil {
					log.Fatal(e)
				}
			} else { //gzip it if not already
				fileExtension := filepath.Ext(file)
				if fileExtension != ".gz" {
					err := gzipit(file, filepath.Dir(file))
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}
}

// rangeDate returns a date range function over start date to end date inclusive.
// After the end of the range, the range function returns a zero date,
// date.IsZero() is true.
func rangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}

// Contains tells whether a contains x.
func contains(a []string, x string) bool {
	var found bool
	for _, n := range a {
		found = strings.Contains(x, n)
		if found {
			return true
		}
	}
	return false
}

// Compress (.gz) a file and delete
func gzipit(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}

	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.gz", filename))
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)

	reader.Close()
	e := os.Remove(source)
	if e != nil {
		log.Fatal(e)
	}
	return err
}
