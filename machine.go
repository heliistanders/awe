package main

import (
	"errors"
	"fmt"
	"github.com/docker/distribution/context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Machine defines all available fields on a vuln target
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
		validated = false
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

// context for every docker cli
var ctx = context.Background()

func getAllMachines() []Machine {
	images := getAllImages()
	owns := GetOwnedMachines()
	var machines []Machine

	for _, image := range images {
		var machine Machine
		for key, value := range image.Labels {
			switch key {
			case "awe":
				machine.Name = value
			case "difficulty":
				machine.Difficulty = value
			case "ports":
				machine.InternalPorts = strings.Split(value, ",")
			}
		}
		machine.Image = image.RepoTags[0]
		if value, ok := owns[machine.Image]; ok {
			machine.Owned = true
			machine.OwnedAt = value
		}

		if checkStatusByImageName(machine.Image) {
			machine.Status = "running"
		} else {
			machine.Status = "stopped"
		}
		machine.Ports = getPublicPortsByImageName(machine.Image)
		if machine.Validate() {
			machines = append(machines, machine)
		}
	}

	return machines
}

func getMachineByImage(image string) Machine {
	machines := getAllMachines()
	var machine Machine
	for _,m := range machines {
		if m.Image == image {
			machine = m
		}
	}
	return machine
}

func checkStatusByImageName(image string) bool {
	containers := getAllContainers()
	for _, c := range containers {
		if c.Image == image {
			return true
		}
	}
	return false
}

func getPublicPortsByImageName(image string) []string {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer func() {
		err = cli.Close()
		fmt.Println(err)
	}()
	var ports []string
	containers := getAllContainers()
	for _, c := range containers {
		if c.Image == image {
			con, err := cli.ContainerInspect(ctx, c.ID)
			if err != nil {
				panic(err)
			}
			portBinding := con.HostConfig.PortBindings
			for _, v := range portBinding {
				for _, port := range v {
					ports = append(ports, port.HostPort)
				}
			}
		}
	}
	return ports
}

func getAllImages() []types.ImageSummary {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer func() {
		err = cli.Close()
		fmt.Println(err)
	}()

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	var aweImages []types.ImageSummary

	for _, image := range images {
		for key := range image.Labels {
			if key == "awe" {
				aweImages = append(aweImages, image)
			}
		}
	}

	return aweImages
}

func getAllContainers() []types.Container {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer func() {
		err = cli.Close()
		fmt.Println(err)
	}()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		_, err := getImageByName(c.Image)
		if err != nil {
			panic(err)
		}
	}

	return containers
}

func getImageByName(name string) (types.ImageSummary, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer func() {
		err = cli.Close()
		fmt.Println(err)
	}()

	images, err := cli.ImageList(ctx, types.ImageListOptions{})

	var aweImage types.ImageSummary
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		for key := range image.Labels {
			if key == "awe" {
				if image.RepoTags[0] == name {
					return image, nil
				}
			}
		}
	}

	return aweImage, errors.New("no Image found")
}


func (m *Machine) StartMachine(flag string) (bool, error) {
	machine := m

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false, err
	}
	defer func() {
		err = cli.Close()
		fmt.Println(err)
	}()


	// create mapping for ever exposed port of the machine
	exposedPorts := make(map[nat.Port]struct{})
	portBindings := make(map[nat.Port][]nat.PortBinding)
	fmt.Println(machine.InternalPorts)
	for _, port := range machine.InternalPorts {
		// exposed ports
		p, err := nat.NewPort("tcp", port)
		if err != nil {
			fmt.Println(err.Error())
			return false, err
		}
		exposedPorts[p] = struct{}{}

		// portBindings
		portBindings[p] = []nat.PortBinding{
			{
				HostIP: "0.0.0.0",
				HostPort: strconv.Itoa(rand.Intn(64512) + 1024), // port >1024 && < 65 536
			},
		}
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		LogConfig: container.LogConfig{
			Type:   "json-file",
			Config: map[string]string{},
		},
	}

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{},
	}
	gatewayConfig := &network.EndpointSettings{
		Gateway: "gatewayname",
	}
	networkConfig.EndpointsConfig["bridge"] = gatewayConfig

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        machine.Image,
		ExposedPorts: exposedPorts,
	}, hostConfig, networkConfig, nil, "")
	if err != nil {
		return false, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return false, err
	}

	command := types.ExecConfig{
		Cmd: []string{"sh", "-c", "echo " + flag + " > /flag.txt"},
	}

	idResponse, err := cli.ContainerExecCreate(ctx, resp.ID, command)
	if err != nil {
		return false, err
	}

	hi, err := cli.ContainerExecAttach(ctx, idResponse.ID, types.ExecStartCheck{})
	if err != nil {
		return false, err
	}
	defer hi.Close()

	return true, nil
}

func (m *Machine) stopMachine() (bool, error){

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false, err
	}
	defer func() {
		err = cli.Close()
		fmt.Println(err)
	}()

	containerList, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return false, err
	}

	for _, cont := range containerList {
		if cont.Image == m.Image {
			timeout := time.Second * 30
			err = cli.ContainerStop(ctx, cont.ID, &timeout)
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}
	return false, errors.New("can't stop Machine")
}

func removeContainer(name string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer func() {
		err = cli.Close()
		fmt.Println(err)
	}()

	containers := getAllContainers()

	for _, c := range containers {
		if c.Image == name {
			err = cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{Force: true})
			if err != nil {
				panic(err)
			}

		}
	}

}
