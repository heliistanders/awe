package service

import (
	"awe/aweDocker"
	"awe/database"
	"awe/model"
	"database/sql"
	"errors"
	"log"
	"strings"
)

type MachineService struct {
	awe          *aweDocker.AweDocker
	machineStore *database.OwnedMachineStore
	flagStore    *database.TempFlagStore
}

func NewMachineService(awe *aweDocker.AweDocker, db *sql.DB) *MachineService {
	return &MachineService{
		machineStore: database.NewOwnedMachineStore(db),
		flagStore:    database.NewTempFlagStore(db),
		awe:          awe,
	}
}

func (ms *MachineService) GetAllMachines() ([]model.Machine, error) {
	var machines []model.Machine

	allImages := ms.awe.GetAllAweImages()

	allOwns, err := ms.machineStore.GetAll()
	if err != nil {
		return machines, err
	}

	allContainers, err := ms.awe.GetAllAweContainers()
	if err != nil {
		return machines, err
	}

	for _, image := range allImages {
		var machine model.Machine
		for key, value := range image.Labels {
			switch key {
			case "awe":
				machine.Name = value
			case "difficulty":
				machine.Difficulty = value
			case "ports":
				machine.InternalPorts = strings.Split(value, ",")
			case "hint":
				machine.Hint = value
			}
		}
		machine.Status = "not running"
		machine.Image = strings.Split(image.RepoTags[0], ":")[0]

		for _, v := range allOwns {
			if v.Image == machine.Image {
				machine.Owned = true
				machine.OwnedAt = v.OwnedAt
			}
		}

		for _, v := range allContainers {
			log.Println(v.Image + ":" + machine.Image)
			if v.Image == machine.Image {
				log.Println("Container Status: " + v.Status)
				switch v.Status {
				case "":
					machine.Status = "not running"
				default:
					machine.Status = "not running"
				}
				machine.Status = v.Status
				ports, err := ms.awe.GetOpenPorts(v)
				if err != nil {
					return machines, err
				}
				machine.Ports = ports
			}
		}

		if machine.Validate() {
			machines = append(machines, machine)
		}
	}

	return machines, nil
}

func (ms *MachineService) getMachineByName(name string) (model.Machine, error) {
	log.Println("Name: " + name)
	var machine model.Machine
	machines, err := ms.GetAllMachines()
	if err != nil {
		return machine, err
	}

	for _, v := range machines {
		log.Println(name + ":" + v.Image)
		if v.Image == name {
			return v, nil
		}
	}

	return machine, errors.New("cannot find machine with given name")
}

func (ms *MachineService) StartMachine(name string) error {
	return ms.StartMachineWithFlag(name, "")
}

func (ms *MachineService) StartMachineWithFlag(name string, flag string) error {
	log.Println("Name: " + name + ", flag: " + flag)
	machine, err := ms.getMachineByName(name)
	if err != nil {
		return err
	}

	log.Printf("Machine: %v", machine)

	tempFlag := model.TempFlag{
		Image: machine.Image,
		Flag:  flag,
	}
	if _, err := ms.flagStore.Insert(&tempFlag); err != nil {
		return err
	}

	return ms.awe.StartMachineWithFlag(&machine, flag)
}

func (ms *MachineService) StopMachine(name string) error {
	machine, err := ms.getMachineByName(name)
	if err != nil {
		return err
	}

	tempFlag, err := ms.flagStore.FindTempFlagByImage(machine.Image)
	if err != nil {
		return err
	}
	if err = ms.flagStore.Delete(tempFlag); err != nil {
		return err
	}

	return ms.awe.StopMachine(&machine)
}

func (ms *MachineService) RestartMachine(name string) error {
	machine, err := ms.getMachineByName(name)
	if err != nil {
		return err
	}

	return ms.awe.RestartMachine(&machine)
}

func (ms *MachineService) PauseMachine(name string) error {
	machine, err := ms.getMachineByName(name)
	if err != nil {
		return err
	}

	return ms.awe.PauseMachine(&machine)
}

func (ms *MachineService) ResumeMachine(name string) error {
	machine, err := ms.getMachineByName(name)
	if err != nil {
		return err
	}

	return ms.awe.ResumeMachine(&machine)
}

func (ms *MachineService) ResetMachine(name string, flag string) error {
	if err := ms.StopMachine(name); err != nil {
		return err
	}
	if err := ms.StartMachineWithFlag(name, flag); err != nil {
		return err
	}
	return nil
}

func (ms *MachineService) AddMachine(path string) error {
	if err := ms.awe.AddMachine(path); err != nil {
		return err
	}

	return nil
}

func (ms *MachineService) Solve(flag string) error {
	tempFlag, err := ms.flagStore.FindTempFlagByFlag(flag)
	if err != nil {
		log.Printf("cannot find tempFlag: %s", err)
	}
	if tempFlag.Image == "" {
		return errors.New("wrong flag")
	}
	machine, err := ms.getMachineByName(tempFlag.Image)
	if err != nil {
		return err
	}

	ownedMachine, err := ms.machineStore.Insert(&machine)
	if err != nil {
		return err
	}

	if ownedMachine.OwnedAt == "" {
		return errors.New("this should not happen ...")
	}

	return nil
}