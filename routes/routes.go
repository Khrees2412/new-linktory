package routes

import (
	"github.com/gofiber/fiber/v2"

	"linktory/controllers"
)

func Setup(app *fiber.App) {

	// User routes
	app.Post("/api/v1/user/register", controllers.Register)
	app.Post("/api/v1/user/login", controllers.Login)
	app.Post("/api/v1/user/logout", controllers.Logout)
	app.Post("/api/v1/change_password", controllers.ChangePassword)
	app.Post("/api/v1/delete_account", controllers.DeleteAccount)


	//Link routes
	app.Get("api/v1/link/get/:id", controllers.GetLink)
	app.Get("api/v1/link/get_all", controllers.GetLinks)
	app.Post("/api/v1/link/create", controllers.CreateLink)
	app.Post("api/v1/link/update/:id", controllers.UpdateLink)
	app.Post("api/v1/link/delete/:id", controllers.DeleteLink)

}