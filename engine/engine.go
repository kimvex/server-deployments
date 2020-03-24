package engine

import (
	"fmt"

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

func ServerExecute() {
	app := fiber.New()
	app.Use(middleware.Logger())
	app.Use("*", func(c *fiber.Ctx) {
		arrays := Con()
		fmt.Println(arrays)
		for _, v := range arrays {
			if c.Path() == v {
				if c.Get("token") != "" {
					fmt.Println("lo encontro")
				} else {
					fmt.Println("no lo encontro")
				}
			}
		}
		fmt.Println(c.Get("token"))
		fmt.Println("si pasamos", c.Path())
		c.Next()
	})

	database := db.Connect()
	redisC := db.Redisdb()
	GetUser := db.GetUserId()

	routes.Routes(app, database, redisC, GetUser)
	app.Listen(3000)
}
