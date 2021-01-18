package model

import "fmt"

type Machine struct {
	Name          string   `json:"name"`
	Difficulty    string   `json:"difficulty"`
	Owned         bool     `json:"owned"`
	OwnedAt       string   `json:"owned_at"`
	Status        string   `json:"status"`
	Ports         []string `json:"ports"`
	InternalPorts []string `json:"-"`
	Image         string   `json:"image"`
	//mu sync.Mutex `json:"-"`
}

func (m *Machine) Validate() bool {
	validated := true
	if m.Name == "" {
		fmt.Println("Machine Validation Error - Name empty")
		validated = false
	}
	if m.Difficulty == "" {
		fmt.Println("Machine Validation Error - Difficulty empty")
		validated = false
	}
	if m.Status == "" {
		fmt.Println("Machine Validation Error - Status empty")
		//validated = false
	}
	if m.Status == "running" && m.Ports == nil {
		fmt.Println("Machine Validation Error - Ports empty")
		validated = false
	}
	if m.InternalPorts == nil {
		fmt.Println("Machine Validation Error - No internal ports provided")
		validated = false
	}
	if m.Image == "" {
		fmt.Println("Machine Validation Error - Image empty")
		validated = false
	}

	return validated
}