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
		log.Println(err)
	}
	return &OwnedMachineStore{
		db: db,
	}
}

func (s *OwnedMachineStore) Insert(machine *model.Machine) (*model.OwnedMachine, error) {
	ownedMachine := model.OwnedMachine{
		Image: machine.Image,
		OwnedAt: time.Now().String(),
	}
	if machine.Image == "" {
		return &ownedMachine, errors.New("Machine Image is empty")
	}
	result, err := s.db.Exec("INSERT INTO owned(image) VALUES(?)", machine.Image)
	if err != nil {
		log.Println("Cannot insert OwnedMachine: &s", err)
		return &ownedMachine, err
	}

	newId, err := result.LastInsertId()
	if err != nil {
		log.Println("Cannot get last id: %s", err)
		return &ownedMachine, err
	}

	ownedMachine.ID = newId
	return &ownedMachine, nil
}

func (s *OwnedMachineStore) GetAll() ([]model.OwnedMachine, error) {
	all := []model.OwnedMachine{}
	rows, err := s.db.Query("SELECT id, image, owned_at FROM owned")
	if err != nil {
		log.Println("Cannot query all OwnedMachines: &s", err)
		return all, err
	}
	// no need to close rows -> gets closed automatically when Next() return false

	for rows.Next() {
		var (
			id int64
			image string
			ownedAt string
		)

		if err := rows.Scan(&id, &image, &ownedAt); err != nil {
			log.Println("Cannot scan rows: %s", err)
			return all, err
		}
		ownedMachine := model.OwnedMachine{
			ID: id,
			Image: image,
			OwnedAt: ownedAt,
		}

		all = append(all, ownedMachine)
	}

	return all, nil
}