package webserver

import (
	"awe/aweDocker"
	"awe/service"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"math/rand"
	"time"
)

func NewServer(awe *aweDocker.AweDocker, db *sql.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		UnescapePath: true,
		ReadTimeout: time.Second * 60,
		WriteTimeout: time.Second * 60,
		BodyLimit: 1024 * 1024 * 1024 * 2,
	})

	machineService := service.NewMachineService(awe, db)

	// serve our web app
	app.Static("/", "./public")

	app.Get("/machines", func(ctx *fiber.Ctx) error {
		machines, err := machineService.GetAllMachines()
		if err != nil {
			return ctx.SendStatus(500)
		}
		return ctx.JSON(machines)
	})

	app.Get("/machine/:name", func(ctx *fiber.Ctx) error {
		return ctx.SendString("OK")
	})

	app.Get("/start/:name", func(ctx *fiber.Ctx) error {
		name := ctx.Params("name")
		log.Println("Starting: " + name)
		err := machineService.StartMachineWithFlag(name, GenerateRandomFlag(15))
		if err != nil {
			log.Println(err)
		}
		machines, err := machineService.GetAllMachines()
		if err != nil {
			return ctx.SendStatus(500)
		}
		return ctx.JSON(machines)
	})

	app.Get("/stop/:name", func(ctx *fiber.Ctx) error {
		name := ctx.Params("name")
		err := machineService.StopMachine(name)
		if err != nil {
			log.Println(err)
		}
		machines, err := machineService.GetAllMachines()
		if err != nil {
			return ctx.SendStatus(500)
		}
		return ctx.JSON(machines)
	})

	app.Get("/restart/:name", func(ctx *fiber.Ctx) error {
		name := ctx.Params("name")
		machineService.RestartMachine(name)
		return ctx.SendString("OK")
	})

	app.Get("/pause/:name", func(ctx *fiber.Ctx) error {
		name := ctx.Params("name")
		machineService.PauseMachine(name)
		return ctx.Status(200).SendString("OK")
	})

	app.Get("/resume/:name", func(ctx *fiber.Ctx) error {
		name := ctx.Params("name")
		machineService.ResumeMachine(name)
		return ctx.SendString("OK")
	})

	app.Post("/solve", func(ctx *fiber.Ctx) error {
		var flag string

		// cannot parse Body => 400 BAD REQUEST
		if err := ctx.BodyParser(flag); err != nil {
			return ctx.SendStatus(400)
		}

		// TODO: check Flag



		return ctx.SendString("OK")

	})

	app.Post("add", func(ctx *fiber.Ctx) error {
		form, err := ctx.MultipartForm()
		if err != nil {
			log.Println("Cannot get multipart form")
			log.Println(err)
			return ctx.SendStatus(400)
		}

		formFiles := form.File["machine"]
		if len(formFiles) > 1 {
			log.Println("More than one file")
			log.Println(err)
			return ctx.SendStatus(400)
		}

		passwords := form.Value["password"]
		if len(passwords) > 1 {
			log.Println("More than one password")
			log.Println(err)
			return ctx.SendStatus(400)
		}
		if passwords[0] != "123" {
			log.Println("Wrong password")
			log.Println(err)
			return ctx.SendStatus(400)
		}

		formFile := formFiles[0]

		filePath := fmt.Sprintf("/tmp/%s.tar", "test")

		err = ctx.SaveFile(formFile, filePath)
		if err != nil {
			log.Println(err)
			ctx.SendStatus(500)
		}
		if err := machineService.AddMachine(filePath); err != nil {
			return ctx.SendStatus(500)
		}

		return ctx.SendString("OK")
	})



	return app
}

func GenerateRandomFlag(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}