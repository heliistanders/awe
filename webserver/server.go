package webserver

import (
	"awe/aweDocker"
	"awe/service"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"log"
	"math/rand"
	"time"
)

func NewServer(awe *aweDocker.AweDocker, db *sql.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		UnescapePath: true,
		ReadTimeout: time.Second * 30,
		WriteTimeout: time.Second * 30,
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