package route

import (
	ctrl "github.com/binarstrike/magic-url/internal/delivery/http/controller"
	"github.com/binarstrike/magic-url/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoute(app *fiber.App, uc *ctrl.UserController, m *middleware.Middleware) *fiber.App {
	app.Post(ctrl.UserRegisterRoute, uc.Register)
	app.Post(ctrl.UserLoginRoute, uc.Login)

	app.Use(ctrl.BaseUserRoute, m.AuthSessionMiddleware)
	app.Get(ctrl.UserCurrentRoute, uc.Me)
	app.Post(ctrl.UserLogoutRoute, uc.Logout)
	// TODO: karena beberapa route perlu hak yang lebih tinggi untuk mengelola sumber daya aplikasi
	// maka diperlukan implementasi role based authentication untuk menyaring akses antara pengguna biasa dan admin
	// app.Delete(ctrl.UserDeleteRoute, app.UserController.Delete)

	return app
}
