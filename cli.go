package main

import (
	"flag"
	"fmt"
)

var (
	lockfile  = flag.String("lockfile", "", "create a lockfile at <path>. (default /var/run/<command>")
	pidfile   = flag.String("pidfile", "", "explicitly create a pidfile at <path>. Without this option, lockfile will be used as a pidfile")
	ignorepid = flag.Bool("ignore-pidfile", false, "ignore any pidfiles found for <command>")
)

// Opts contains the CLI options resolved against the defaults
type Opts struct {
	Pidfile       string
	Lockfile      string
	createPidfile bool
	Command       string
	CommandArgs   []string
}

// Return whether or not user has chosen to use a pidfile
func (o *Opts) usePidfile() bool {
	return o.createPidfile
}

// ParseCLI instantiates a new Opts object with options passed from CLI flags
func ParseCLI() Opts {
	flag.Parse()
	args := flag.Args()

	opts := new(Opts)

	opts.Command = args[0]
	opts.CommandArgs = args[1:]

	if *lockfile == "" {
		opts.Lockfile = fmt.Sprintf("/var/run/%s.lock", opts.Command)
	} else {
		opts.Lockfile = *lockfile
	}

	if *pidfile != "" {
		opts.createPidfile = true
		opts.Pidfile = *pidfile
	}

	return *opts
}
