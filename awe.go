package main

import (
	"awe/aweDocker"
	"awe/webserver"
	"database/sql"
	"github.com/docker/docker/client"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	// initialize rng
	rand.Seed(time.Now().UnixNano())
	// initialize logger -> show line number
	log.SetFlags(log.Lshortfile | log.Ldate)

	// setup adminPassword
	adminPassword := webserver.GenerateRandomFlag(10)
	envPW := os.Getenv("AWE_PASS")
	if envPW != "" {
		adminPassword = envPW
	} else {
		log.Printf("ADMIN PASSWORD = %s", adminPassword)
	}

	// check Website
	if err := indexPageExists(); err != nil {
		log.Fatal("index.html not found")
	}

	// setup and check Docker
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	aweDockerInstance := aweDocker.NewAweDocker(cli)

	// check if we can access docker
	err = aweDockerInstance.IsAvailable()
	if err != nil {
		log.Fatal(err)
	}

	// open Database
	db, err := sql.Open("sqlite3", "awe.sqlite")
	if err != nil {
		log.Fatal("Cannot open Database")
	}
	defer db.Close()

	// setup web server
	app := webserver.NewServer(aweDockerInstance, db, adminPassword)

	// start WebServer
	err = app.Listen(":5000")
	if err != nil {
		log.Fatal(err)
	}
}

func indexPageExists() error {
	if _, err := os.Stat("./public"); os.IsNotExist(err) {
		log.Println("public folder not found. \n creating public folder ...")
		err := os.Mkdir("public", 0775)
		if err != nil {
			log.Println("Cannot create public folder")
			return err
		}
	}
	if _, err := os.Stat("./public/index.html"); os.IsNotExist(err) {
		log.Println("index.html not found")
		return err
	}

	return nil
}
