package main

import (
	"fmt"
	"github.com/docker/distribution/context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
	"math/rand"
	"net/url"
	"os"
	"sync"
	"time"
)

func messageLogger(messages <-chan events.Message) {
	for msg := range messages {
		fmt.Println("[+] - " + msg.Action)
	}
}

func errorLogger(errors <-chan error) {
	for err := range errors {
		fmt.Println("[+] - " + err.Error())
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	startupChecks()

	// not a TRNG but good enough
	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	messages, errors := cli.Events(ctx, types.EventsOptions{})
	go messageLogger(messages)
	go errorLogger(errors)

	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		fmt.Println("trying to connect ...")
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			ad := []byte("- From Server")
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, append(msg, ad...)); err != nil {
				log.Println("write:", err)
				break
			}
		}

	}))

	app.Static("/", "./public")

	app.Get("/machines", func(c *fiber.Ctx) error {
		// return c.JSON(fiber.Map{
		// 	"ok":       true,
		// 	"machines": machines,
		// })
		return c.JSON(getAllMachines())
	})

	app.Get("/machine/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		decodedName, err := url.QueryUnescape(name)
		if err != nil {
			panic(err)
		}
		machine := getMachineByImage(decodedName)
		if machine.Name == "" {
			return c.SendString("Sorry")
		}

		return c.JSON(machine)
	})

	// preventing multiple instances of the same machine
	var machineMutex sync.Mutex
	app.Get("/start/:name", func(c *fiber.Ctx) error {
		machineMutex.Lock()
		defer machineMutex.Unlock()
		name := c.Params("name")
		decodedName, err := url.QueryUnescape(name)
		if err != nil {
			panic(err)
		}
		flag := generateFlag(15)
		machine := getMachineByImage(decodedName)
		if machine.Status == "running" {
			fmt.Println("machine already running")
			return c.JSON(getAllMachines())
		}
		success, err := machine.StartMachine(flag)
		if err != nil {
			fmt.Println(err)
			return c.SendString("NOK")
		}
		if success {
			createFlag(machine, flag)
		}
		return c.JSON(getAllMachines())
	})

	app.Get("/stop/:name", func(c *fiber.Ctx) error {
		machineMutex.Lock()
		defer machineMutex.Unlock()
		name := c.Params("name")
		decodedName, err := url.QueryUnescape(name)
		if err != nil {
			panic(err)
		}
		fmt.Println("Stopping Machine: " + decodedName)
		machine := getMachineByImage(decodedName)
		removeContainer(decodedName)
		deleteFlag(machine)

		return c.JSON(getAllMachines())
	})

	app.Get("/restart/:name", func(c *fiber.Ctx) error {
		machineMutex.Lock()
		defer machineMutex.Unlock()
		name := c.Params("name")
		decodedName, err := url.QueryUnescape(name)
		if err != nil {
			panic(err)
		}
		fmt.Println("Stopping Machine: " + decodedName)

		machine := getMachineByImage(decodedName)

		flag := generateFlag(15)
		fmt.Println("Flag: " + flag)
		success, err := machine.stopMachine()
		if err != nil || !success {
			return c.JSON(getAllMachines())
		}

		success, err = machine.StartMachine(flag)
		if err != nil || !success {
			return c.JSON(getAllMachines())
		}
		createFlag(machine, flag)

		return c.JSON(getAllMachines())
	})

	type solveMsg struct {
		MachineName string `json:"machineName"`
		Flag string `json:"flag"`
	}
	app.Post("/solve", func(c *fiber.Ctx) error {
		msg := new(solveMsg)
		if err := c.BodyParser(msg); err != nil {
			fmt.Println(err.Error())
			return c.SendString("NOK")
		}
		log.Println(msg.MachineName, msg.Flag)
		machine := getMachineByImage(msg.MachineName)
		fmt.Println("This machine? " + machine.Image)

		if checkFlag(machine, msg.Flag) {
			if OwnMachine(machine) {
				fmt.Println("Solved!")
				return c.SendString("OK")
			}
		}
		fmt.Println("Wrong flag!")

		return c.SendString("NOK")
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		images := getAllImages()
		return c.JSON(images)
	})

	err = app.Listen(":5000")
	if err != nil {
		panic(err)
	}
}

func checkWebPages() error {
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

func generateFlag(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
