package engine

import (
	"unsafe"

	"../db"
	"../routes"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
)

type BodyRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	Same_password string `json:"same_password"`
}

type ErroRespnse struct {
	MESSAGE string `json:"message"`
}

type SuccessResponse struct {
	MESSAGE string `json:"message_success"`
}

type AlreadyRegistered struct {
	Email string `json:"email"`
}

type UserBody struct {
	UserId int `json:"userId"`
}

var getString = func(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func ServerExecute() {
	app := fiber.New()
	app.Use(middleware.Logger())

	database := db.Connect()
	redisC := db.Redisdb()
	GetUser := db.GetUserId()

	routes.Routes(app, database, redisC, GetUser)
	app.Listen(3000)
}
