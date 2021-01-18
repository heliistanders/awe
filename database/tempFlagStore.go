package database

import (
	"awe/model"
	"database/sql"
	"errors"
	"log"
)
import _ "github.com/mattn/go-sqlite3"

type TempFlagStore struct {
	db *sql.DB
}

func NewTempFlagStore(db *sql.DB) *TempFlagStore {
	createStmt := `CREATE TABLE IF NOT EXISTS "flags" (
					"id" integer,
					"image" text NOT NULL UNIQUE ,
					"flag" text NOT NULL UNIQUE , 
					PRIMARY KEY (id)
	)`

	_, err := db.Exec(createStmt)
	if err != nil {
		log.Println(err)
	}

	return &TempFlagStore{
		db: db,
	}
}

func (s *TempFlagStore) Insert(tempFlag *model.TempFlag) (*model.TempFlag, error) {
	result, err := s.db.Exec("INSERT OR REPLACE INTO flags(image, flag) VALUES(?, ?)", tempFlag.Image, tempFlag.Flag)
	if err != nil {
		log.Println("Cannot insert TempFlag: %s", err)
		return tempFlag, err
	}

	newId, err := result.LastInsertId()
	if err != nil {
		log.Println("Cannot get newId: %s", err)
		return tempFlag, err
	}

	tempFlag.ID = newId
	return tempFlag, nil
}

func (s *TempFlagStore) FindTempFlagByImage(image string) (*model.TempFlag, error) {
	tempFlag := model.TempFlag{Image: image}
	if image == "" {
		return &tempFlag, errors.New("Image name is empty")
	}
	var newId int64
	var newFlag string
	err := s.db.QueryRow("SELECT id, flag FROM flags WHERE image = ?", image).Scan(&newId, &newFlag)
	if err != nil {
		log.Println("Cannot query TempFlag by Image: %s", err)
		return &tempFlag, err
	}

	tempFlag.ID = newId
	tempFlag.Flag = newFlag

	return &tempFlag, nil
}

func (s *TempFlagStore) FindTempFlagByFlag(flag string) (*model.TempFlag, error) {
	tempFlag := model.TempFlag{Flag: flag}
	if flag == "" {
		return &tempFlag, errors.New("flag is empty")
	}
	var newId int64
	var newImage string
	err := s.db.QueryRow("SELECT id, image FROM flags WHERE flag = ?", flag).Scan(&newId, &newImage)
	if err != nil {
		log.Println("Cannot query TempFlag by Flag: %s", err)
		return &tempFlag, err
	}

	tempFlag.ID = newId
	tempFlag.Image = newImage

	return &tempFlag, nil
}

func (s *TempFlagStore) Delete(tempFlag *model.TempFlag) (error) {
	_, err := s.db.Exec("DELETE FROM flags WHERE image = ? OR flag = ?",tempFlag.Image, tempFlag.Image)
	if err != nil {
		log.Println("Cannot delete TempFlag: %s", err)
		return err
	}

	return nil
}