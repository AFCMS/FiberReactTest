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

func main() {
	DevMode := os.Getenv("FIBER_REACT_DEV_MODE") == "true"
	FrontendDevServer := os.Getenv("FIBER_REACT_FRONTEND_SERVER")
	GoogleSiteVerification := os.Getenv("FIBER_REACT_GOOGLE_SITE_VERIFICATION")

	engine := html.New("./", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(flogger.New())

	// API
	api := app.Group("/api")

	api.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api.All("*", func(ctx *fiber.Ctx) error {
		return ctx.Status(404).SendString("Not Found")
	})

	if DevMode {
		templateHandler := func(c *fiber.Ctx) error {
			return c.Render("index", fiber.Map{
				"Title":   "Test",
				"DevMode": DevMode,
			})
		}
		app.Get("*", proxy.BalancerForward([]string{FrontendDevServer}), templateHandler)
		app.Get("/", templateHandler)
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
		app.Get("*", func(c *fiber.Ctx) error {
			return c.Render("index", fiber.Map{
				"Title":                  "Test",
				"DevMode":                DevMode,
				"MainCSS":                manifest["src/main.tsx"].Css[0],
				"MainJS":                 manifest["src/main.tsx"].File,
				"GoogleSiteVerification": GoogleSiteVerification,
			})
		})
	}

	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
