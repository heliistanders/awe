package main

import (
	"log"
)

func startupChecks(){
	err := setUpDatabase()
	if err != nil {
		log.Fatal("Database not found")
	}
	err = checkWebPages()
	if err != nil {
		log.Fatal("WebContent not found")
	}
	err = checkDockerInstallation()
	if err != nil {
		log.Fatal("Docker not found")
	}
}