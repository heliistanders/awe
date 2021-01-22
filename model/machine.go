package model

import (
	"fmt"
	"log"
)

type Machine struct {
	Name          string   `json:"name"`
	Difficulty    string   `json:"difficulty"`
	Owned         bool     `json:"owned"`
	OwnedAt       string   `json:"owned_at"`
	Status        string   `json:"status"`
	Ports         []string `json:"ports"`
	InternalPorts []string `json:"-"`
	Image         string   `json:"image"`
	Hint          string   `json:"hint"`
}

func (m *Machine) Validate() bool {
	log.Println("Validating Machine: " + m.Name)
	validated := true
	if m.Name == "" {
		log.Println("Machine Validation Error - Name empty")
		validated = false
	}
	if m.Difficulty == "" {
		log.Println("Machine Validation Error - Difficulty empty")
		validated = false
	}
	if m.Status == "" {
		log.Println("TODO: Machine Validation Error - Status empty")
		//validated = false
	}
	if m.Status == "running" && m.Ports == nil {
		log.Println("Machine Validation Error - Ports empty")
		validated = false
	}
	if m.InternalPorts == nil {
		fmt.Println("Machine Validation Error - No internal ports provided")
		//validated = false
	}
	if m.Image == "" {
		log.Println("Machine Validation Error - Image empty")
		validated = false
	}

	if m.Hint == "" {
		log.Println("Optional: no hint provided")
	}

	return validated
}
