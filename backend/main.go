package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	"backend/db"
	_ "backend/docs"
	"backend/routes"
)

// @title           Workflow Approval API
// @version         1.0
// @description     IT03 document approval workflow (Pending/Approved/Rejected) with reason.
// @host            localhost:8080
// @BasePath        /
func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db.Connect()
	defer db.Disconnect()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault)

	routes.Register(app)

	log.Printf("✅ Server: http://localhost:%s", port)
	log.Printf("✅ Swagger: http://localhost:%s/swagger/index.html", port)
	log.Fatal(app.Listen(":" + port))
}
