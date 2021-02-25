package database

import (
	"awe/model"
	"database/sql"
	"errors"
	"log"
	"time"
)

type OwnedMachineStore struct {
	db *sql.DB
}

func NewOwnedMachineStore(db *sql.DB) *OwnedMachineStore {
	createStmt := `CREATE TABLE IF NOT EXISTS "owned" (
		"id" integer,
		"image" varchar NOT NULL UNIQUE,
		"owned_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
		)`

	_, err := db.Exec(createStmt)
	if err != nil {
		log.Printf("cannot create table owend: %s", err)
	}
	return &OwnedMachineStore{
		db: db,
	}
}

func (s *OwnedMachineStore) Insert(machine *model.Machine) (*model.OwnedMachine, error) {
	ownedMachine := model.OwnedMachine{
		Image:   machine.Image,
		OwnedAt: time.Now().String(),
	}
	if machine.Image == "" {
		return &ownedMachine, errors.New("machine image is empty")
	}
	result, err := s.db.Exec("INSERT INTO owned(image) VALUES(?)", machine.Image)
	if err != nil {
		log.Printf("cannot insert OwnedMachine: %s", err)
		return &ownedMachine, err
	}

	newId, err := result.LastInsertId()
	if err != nil {
		log.Printf("cannot get last id: %s", err)
		return &ownedMachine, err
	}

	ownedMachine.ID = newId
	return &ownedMachine, nil
}

func (s *OwnedMachineStore) GetAll() ([]model.OwnedMachine, error) {
	var all []model.OwnedMachine
	rows, err := s.db.Query("SELECT id, image, owned_at FROM owned")
	if err != nil {
		log.Println("Cannot query all OwnedMachines: &s", err)
		return all, err
	}
	// no need to close rows -> gets closed automatically when Next() return false

	for rows.Next() {
		var (
			id      int64
			image   string
			ownedAt string
		)

		if err := rows.Scan(&id, &image, &ownedAt); err != nil {
			log.Printf("cannot scan rows: %s", err)
			return all, err
		}
		ownedMachine := model.OwnedMachine{
			ID:      id,
			Image:   image,
			OwnedAt: ownedAt,
		}

		all = append(all, ownedMachine)
	}

	return all, nil
}

func (s *OwnedMachineStore) IsOwned(name string) (bool, error) {
	ownes, err := s.GetAll()
	if err != nil {
		return false, err
	}

	for _, o := range ownes {
		if o.Image == name {
			return true, nil
		}
	}

	return false, nil
}
