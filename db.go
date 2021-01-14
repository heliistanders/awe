package main

// TODO: commits / rollback - defer .. very bad code atm

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
)
import _ "github.com/mattn/go-sqlite3"

func setUpDatabase() error {

	if _, err := os.Stat("awe.sqlite"); os.IsNotExist(err) {
		// File does not exist -> create one
		log.Println("Database does not exists. Creating Database ...")
		file, err := os.Create("awe.sqlite")
		if err != nil {
			return err
		}
		file.Close()
		log.Println("Database created")
	}

	db, err := sql.Open("sqlite3", "awe.sqlite")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	createFlagsTable := `CREATE TABLE IF NOT EXISTS"flags" (
							"id" integer,
							"image" text NOT NULL UNIQUE ,
							"flag" text NOT NULL UNIQUE , 
							PRIMARY KEY (id)
	)`
	createOwnsTable := `CREATE TABLE IF NOT EXISTS "owned" (
							"id" integer,
							"image" varchar NOT NULL UNIQUE,
							"owned_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, 
							PRIMARY KEY (id)
	)`

	_, err = db.Exec(createFlagsTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(createOwnsTable)
	if err != nil {
		return err
	}

	return nil
}

func GetOwnedMachines() map[string]string {
	db, err := sql.Open("sqlite3", "awe.sqlite")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	rows, err := db.Query("SELECT image, owned_at from owned")
	if err != nil {
		fmt.Printf(err.Error())
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	machines := make(map[string]string)

	for rows.Next() {
		var image, ownedAt string
		err = rows.Scan(&image, &ownedAt)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		machines[image] = ownedAt
	}
	return machines
}

func OwnMachine(machine Machine) bool {
	if machine.Image == "" {
		return false
	}
	imageName := machine.Image
	fmt.Println("Saving Own of: " + imageName)
	db, err := sql.Open("sqlite3", "awe.sqlite")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return false
	}
	defer func() {
		err = tx.Commit()
		if err != nil {
			log.Println(err)
		}
	}()

	stmt, err := tx.Prepare("insert into owned(image) values(?)")
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = stmt.Exec(imageName)



	return true
}

func createFlag(machine Machine, flag string) bool {
	if machine.Image == "" {
		return false
	}
	imageName := machine.Image
	fmt.Println("Saving Own of: " + imageName)
	db, err := sql.Open("sqlite3", "awe.sqlite")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return false
	}
	defer func() {
		err = tx.Commit()
		if err != nil {
			log.Println(err)
		}
	}()

	stmt, err := tx.Prepare("insert or replace into flags(image, flag) values(?, ?)")
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = stmt.Exec(imageName, flag)



	return true
}

func deleteFlag(machine Machine) bool{
	if machine.Image == "" {
		return false
	}
	imageName := machine.Image
	fmt.Println("Saving Own of: " + imageName)
	db, err := sql.Open("sqlite3", "awe.sqlite")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return false
	}
	defer func() {
		err = tx.Commit()
		if err != nil {
			log.Println(err)
		}
	}()

	stmt, err := tx.Prepare("delete from flags where image = ?")
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = stmt.Exec(imageName)


	return true
}

func checkFlag(machine Machine, flag string) bool {

	if machine.Image == "" {
		return false
	}
	imageName := machine.Image

	db, err := sql.Open("sqlite3", "awe.sqlite")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return false
	}
	defer func() {
		err = tx.Commit()
		if err != nil {
			log.Println(err)
		}
	}()

	stmt, err := tx.Prepare("select count(*) from flags where image = ? and flag = ?")
	if err != nil {
		log.Println(err)
		return false
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	var count int
	var ctx context.Context
	ctx = context.Background()
	err = stmt.QueryRowContext(ctx, imageName, flag).Scan(&count)
	if err != nil {
		log.Println(err)
		return false
	}

	if count == 1 {
		return true
	}


	return true
}
