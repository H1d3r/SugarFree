package main

import (
	"SugarFree/Packages/Arguments"
	"SugarFree/Packages/Utils"
	"log"
	"os"
)

// main function
func main() {
	logger := log.New(os.Stderr, "[!] ", 0)

	// Call function named CheckGoVersion
	Utils.CheckGoVersion()

	//  SugarFreeCli Execute
	err := Arguments.SugarFreeCli.Execute()
	if err != nil {
		logger.Fatal("Error: ", err)
		return
	}
}
