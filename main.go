package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/template/html/v2"
	"log"
	"os"
)

const DevMode = true

func main() {
	engine := html.New("./", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(flogger.New())

	app.Get("/gif", proxy.Forward("https://i.imgur.com/IWaBepg.gif"))

	// API
	api := app.Group("/api")

	api.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	if DevMode {
		app.Get("*", proxy.BalancerForward([]string{"http://frontend_server:5173"}))
		app.Get("/", func(c *fiber.Ctx) error {
			return c.Render("index", fiber.Map{
				"Title":   "Test",
				"DevMode": DevMode,
			})
		})
	} else {
		// Parse JSON Vite manifest
		manifest := map[string]Chunk{}
		data, err := os.ReadFile("./frontend/dist/.vite/manifest.json")
		if err != nil {
			log.Fatal(err)
		}
		if err == nil {
			err = json.Unmarshal(data, &manifest)
		}

		app.Static("/", "./frontend/dist")
		//app.Static("*", "./frontend/dist/index.html")
		app.Get("/", func(c *fiber.Ctx) error {
			return c.Render("index", fiber.Map{
				"Title":   "Test",
				"DevMode": DevMode,
				"MainCSS": manifest["src/main.tsx"].Css[0],
				"MainJS":  manifest["src/main.tsx"].File,
			})
		})
	}

	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
