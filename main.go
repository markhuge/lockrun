package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	// resolve options
	opts := ParseCLI()

	// read pid from file if exists
	pid, files := resolveFiles(opts)
	fmt.Printf("PID: %d", pid)

	// if pid/lock exists, check to see if process is still running
	if pid != 0 && isRunning(pid) {
		log.Printf("Operation \"%s\" already running with PID: %d", opts.Command, pid)
		os.Exit(0)
	}

	// start process
	err := runCommand(opts, files)

	// TODO handle error
	if err != nil {
		panic(err)
	}

}
func createFile(file string, value int) (err error) {
	f, err := os.Create(file)
	if err != nil {
		return
	}

	// Expose error if f.Close() fails
	// See https://github.com/kisielk/errcheck/issues/55
	defer func() {
		err = f.Close()
	}()

	_, err = f.WriteString(fmt.Sprintf("%d", value))
	return
}

func runCommand(opts Opts, files []string) (err error) {
	// exec
	cmd := exec.Command(opts.Command, opts.CommandArgs...)
	cmd.Stdout = os.Stdout

	err = cmd.Start()
	//TODO handle error
	if err != nil {
		log.Println("Error in start")
		return
	}

	// create pid/lock
	pid := cmd.Process.Pid

	// Walk file list for creation
	for _, file := range files {
		err = createFile(file, pid)

		if err != nil {
			log.Println("Error in createfile")
			return
		}
	}

	err = cmd.Wait()
	// TODO handle inavlid commands
	// Always clenup on this error
	if err != nil {
		log.Println("Error in wait")
		return
	}

	// cleanup pid/lock
	err = cleanupFiles(files)
	return
}

func cleanupFiles(files []string) (err error) {
	for _, file := range files {
		if file != "" {
			// guard against deleting directories
			if info, err := os.Stat(file); err == nil && !info.IsDir() {
				err = os.Remove(file)
				if err != nil {
					return err
				}

			}
		}
	}
	return
}

// TODO make this work cross platform
func isRunning(pid int) bool {
	return fileExists(fmt.Sprintf("/proc/%d", pid))
}

func readFile(filename string) (pid int) {
	data, err := ioutil.ReadFile(filename)

	// errors shouldn't be fatal here
	if err != nil {
		return
	}

	// TODO handle this error
	pid, err = strconv.Atoi(string(data))
	return
}

func fileExists(filename string) (found bool) {
	// treat directories as not found
	if info, err := os.Stat(filename); err == nil && !info.IsDir() {
		found = true
	}
	return
}

// resolve pid in order: pidfile, lockfile
func resolveFiles(opts Opts) (pid int, files []string) {
	// Build array of files to check (in order)
	files = []string{opts.Lockfile}

	// Prepend Pidfile if specified
	if opts.usePidfile() && opts.Pidfile != "" {
		files = append([]string{opts.Pidfile}, files...)
	}

	for _, f := range files {
		if fileExists(f) {
			if pid = readFile(f); pid != 0 {
				return
			}
		}
	}

	return
}
