package core

import (
	"contra/src/collectors"
	"contra/src/configuration"
	"contra/src/utils/git"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const version = "1.0.0"

// Application holds global application data and functions for kicking off execution.
type Application struct {
	config *configuration.Config
}

// Start is the main entrance to the application.
func (a *Application) Start() {
	// Parse the config, which brings in flags.
	a.config = configuration.GetConfig()

	// Display banner first, in case we crash on parsing the config.
	a.DisplayBanner()

	// Determine what to do.
	a.Route()
}

// StandardRun if there are no special cases designated by the configuration.
func (a *Application) StandardRun() {
	// Initialize our main worker.
	worker := collectors.CollectorWorker{
		RunConfig: a.config,
	}

	// Collect everything
	worker.RunCollectors()

	// And check for any necessary commits.
	utils.GitOps(a.config)
}

// Route determines what to do, and kicks off the doing.
func (a *Application) Route() {
	// Determine if the config designates some special run process, otherwise handle our main handler.
	if a.config.Copyrights {
		a.DisplayCopyrights()
	} else if a.config.Debug {
		a.DisplayDebugInfo()
	} else {
		// Standard operating proceedure.
		a.StandardRun()
	}
}

// DisplayBanner with basic information about this application.
func (a *Application) DisplayBanner() {
	// Print something.
	if !a.config.Quiet {
		fmt.Printf("\n=== Contra ===\n"+
			" - Network Device Configuration Tracking\n"+
			" - AJA Video Systems, Inc. Version: %s\n\n", version)
	}
}

// DisplayCopyrights simply dumping the COPYRIGHTS file.
func (a *Application) DisplayCopyrights() {
	log.Println("COPYRIGHT Information")
	data, err := ioutil.ReadFile("COPYRIGHTS")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println(string(data))
	os.Exit(0)
}

// DisplayDebugInfo may print sensitive passwords to the screen.
func (a *Application) DisplayDebugInfo() {
	log.Println("DEBUG ENABLED: Dumping config and exiting.")
	log.Println(a.config)
	os.Exit(0)
}
