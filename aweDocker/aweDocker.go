package aweDocker

import (
	"awe/model"
	"errors"
	"fmt"
	"github.com/docker/distribution/context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type AweDocker struct {
	cli *client.Client
}

func NewAweDocker(cli *client.Client) *AweDocker {
	return &AweDocker{
		cli: cli,
	}
}

var ctx = context.Background()

func (a *AweDocker) IsAvailable() error {
	_, err := a.cli.Info(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (a *AweDocker) StartMachine(machine *model.Machine) error {
	return a.StartMachineWithFlag(machine, "")
}

func (a *AweDocker) StartMachineWithFlag(machine *model.Machine, flag string) error {
	// if theres already a container running with this image => skip
	containers, err := a.GetAllAweContainers()
	if err != nil {
		return err
	}

	for _, cont := range containers {
		if cont.Image == machine.Image {
			return errors.New("machine is already running")
		}
	}

	// create mapping for ever exposed port of the machine
	exposedPorts := make(map[nat.Port]struct{})
	portBindings := make(map[nat.Port][]nat.PortBinding)

	for _, port := range machine.InternalPorts {
		// exposed ports
		p, err := nat.NewPort("tcp", port)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		exposedPorts[p] = struct{}{}

		// portBindings
		portBindings[p] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
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
		Gateway: "AWE",
	}
	networkConfig.EndpointsConfig["bridge"] = gatewayConfig

	resp, err := a.cli.ContainerCreate(ctx, &container.Config{
		Image:        machine.Image,
		ExposedPorts: exposedPorts,
	}, hostConfig, networkConfig, nil, "")
	if err != nil {
		return err
	}

	if err := a.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	if flag != "" {
		command := types.ExecConfig{
			Cmd: []string{"sh", "-c", "echo " + flag + " > /flag.txt"},
		}

		idResponse, err := a.cli.ContainerExecCreate(ctx, resp.ID, command)
		if err != nil {
			return err
		}

		hi, err := a.cli.ContainerExecAttach(ctx, idResponse.ID, types.ExecStartCheck{})
		if err != nil {
			return err
		}
		defer hi.Close()
	}

	return nil
}

func (a *AweDocker) StopMachine(machine *model.Machine) error {
	containerList, err := a.cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, cont := range containerList {
		if cont.Image == machine.Image {
			timeout := time.Second * 30
			err = a.cli.ContainerStop(ctx, cont.ID, &timeout)
			if err != nil {
				return err
			}
			err = a.cli.ContainerRemove(ctx, cont.ID, types.ContainerRemoveOptions{})
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("cannot stop Machine")
}

func (a *AweDocker) RestartMachine(machine *model.Machine) error {
	containerList, err := a.cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, cont := range containerList {
		if cont.Image == machine.Image {
			timeout := time.Second * 30
			err = a.cli.ContainerRestart(ctx, cont.ID, &timeout)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("cannot restart Machine")
}

func (a *AweDocker) ResumeMachine(machine *model.Machine) error {
	containerList, err := a.cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, cont := range containerList {
		if cont.Image == machine.Image {
			err = a.cli.ContainerUnpause(ctx, cont.ID)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("cannot resume Machine")
}

func (a *AweDocker) PauseMachine(machine *model.Machine) error {
	containerList, err := a.cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, cont := range containerList {
		if cont.Image == machine.Image {
			err = a.cli.ContainerPause(ctx, cont.ID)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("cannot pause Machine")
}

func (a *AweDocker) GetAllAweImages() []types.ImageSummary {
	images, err := a.cli.ImageList(ctx, types.ImageListOptions{})
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

func (a *AweDocker) GetAllAweContainers() ([]types.Container, error) {
	var aweContainers []types.Container
	containers, err := a.cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	// filter container which are not labeled as AWE
	for _, container := range containers {
		if _, ok := container.Labels["awe"]; ok {
			aweContainers = append(aweContainers, container)
		}
	}
	return containers, nil
}

func (a *AweDocker) GetOpenPorts(container types.Container) ([]string, error) {
	var ports []string

	con, err := a.cli.ContainerInspect(ctx, container.ID)
	if err != nil {
		return ports, err
	}
	portBinding := con.HostConfig.PortBindings
	// iterate over every PortBinding configuration (could be multiple)
	for _, v := range portBinding {
		// iterate over each port
		for _, port := range v {
			ports = append(ports, port.HostPort)
		}
	}

	return ports, nil
}

func (a *AweDocker) AddMachine(path string) error {
	log.Println("Adding image from path " + path)
	imageFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer imageFile.Close()

	_, err = a.cli.ImageLoad(ctx, imageFile, false)
	if err != nil {
		return err
	}

	return nil
}
