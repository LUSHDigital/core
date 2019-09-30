package core

import "log"

// SetupLogs will set up the appropriate log flags.
// DEPRECATED: Import github.com/LUSHDigital/core-lush/lushlogs as a side-effect.
func SetupLogs() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("DEPRECATED: import github.com/LUSHDigital/core-lush/lushlogs as a side-effect.")
}
