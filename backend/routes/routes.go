package routes

import (
	"github.com/gofiber/fiber/v2"

	"backend/handlers"
)

func Register(app *fiber.App) {
	api := app.Group("/api")

	// List documents
	api.Get("/documents", handlers.ListDocuments)

	// approve documents
	api.Post("/documents/approval", handlers.ApproveDocuments)

	// reject documents
	api.Post("/documents/rejection", handlers.RejectDocuments)

	// mock documents
	api.Post("/documents/seed", handlers.SeedMockDocuments)
}
