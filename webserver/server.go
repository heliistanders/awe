package webserver

import (
	"awe/aweDocker"
	"awe/service"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"io"
	"log"
	"math/rand"
	"net/url"
	"time"
)

func NewServer(awe *aweDocker.AweDocker, db *sql.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		UnescapePath: true,
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
		BodyLimit:    1024 * 1024 * 1024 * 2,
		Concurrency: 1,
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
		// cannot parse Body => 400 BAD REQUEST
		form, err := ctx.MultipartForm()
		if err != nil {
			log.Printf("cannot parse form %s", err)
			return ctx.SendStatus(400)
		}

		flagFields := form.Value["flag"]
		if len(flagFields) != 1 {
			log.Println("not exactly one flag field")
			return ctx.SendStatus(400)
		}

		flag := flagFields[0]
		if err := machineService.Solve(flag); err != nil {
			log.Printf("cannot solve: %s", err)
			return ctx.SendStatus(500)
		}

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

	// Upgrade WebSocket Request
	app.Use("/terminals", func(c *fiber.Ctx) error {
		return c.Next()
	})

	// Upgraded WebSocket request
	app.Get("/terminals", websocket.New(func(conn *websocket.Conn) {
		defer conn.Close()

		encodedMachine := conn.Query("name", "")
		log.Printf("Param: %s", encodedMachine)
		machine, err := url.PathUnescape(encodedMachine)
		if err != nil {
			log.Print(err)
			return
		}
		log.Printf("decodedParam: %s", machine)
		if machine == "" {
			return
		}

		hr, err := machineService.AttachMachine(machine)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage,[]byte(err.Error()))
			log.Print(err)
			return
		}
		defer hr.Close()
		defer func() {
			hr.Conn.Write([]byte("exit\r"))
		}()

		log.Printf("Trying to access machine: %s", machine)


		go func() {
			wsWriterCopy(hr.Conn, conn)
		}()
		wsReaderCopy(conn, hr.Conn)

	}))

	app.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).SendString("pong")
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

func wsWriterCopy(reader io.Reader, writer *websocket.Conn) {
	buf := make([]byte, 8192)
	for {
		nr, err := reader.Read(buf)
		if nr > 0 {
			err := writer.WriteMessage(websocket.BinaryMessage, buf[0:nr])
			if err != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
}

func wsReaderCopy(reader *websocket.Conn, writer io.Writer) {
	for {
		messageType, p, err := reader.ReadMessage()
		if err != nil {
			return
		}
		if messageType == websocket.TextMessage {
			writer.Write(p)
		}
	}
}