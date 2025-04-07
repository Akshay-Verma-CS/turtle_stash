package routes

import (
	"turtle-stash/controllers"

	"github.com/gofiber/fiber/v2"
)

func RegisterFileRoutes(app *fiber.App) {
	api := app.Group("/file")
	api.Post("/", controllers.UploadFile)
	api.Get("/:fileId", controllers.GetDownloadURL)
	api.Delete("/:fileId", controllers.DeleteFile)

	app.Get("/folders/:folderId", controllers.GetFolderSnapshot)
}
