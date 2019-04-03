package core

import "log"

// SetupLogs will set up the appropriate log flags.
func SetupLogs() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
