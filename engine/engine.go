package engine

import (
	"fmt"

	"../db"
	"../routes"
	"github.com/dgrijalva/jwt-go"
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
		passport := true
		for _, v := range arrays {
			fmt.Println(c.Path(), v.Url, "vamos")
			if c.Path() == v.Url && c.Method() == v.Method {
				passport = false
			}
		}

		if passport == false {
			if c.Get("token") != "" {
				token, err := jwt.Parse(c.Get("token"), func(token *jwt.Token) (interface{}, error) {
					return []byte("secret"), nil
				})
				if token.Valid {
					c.Next()
					return
				} else {
					if ve, ok := err.(*jwt.ValidationError); ok {
						if ve.Errors&jwt.ValidationErrorMalformed != 0 {
							c.JSON(ErroRespnse{MESSAGE: "Token structure not valid"})
							c.Status(401)
						} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
							c.JSON(ErroRespnse{MESSAGE: "Token is expired"})
							c.Status(401)
						} else {
							c.JSON(ErroRespnse{MESSAGE: "Invalid token"})
							c.Status(401)
						}
					}
					return
				}
			}

			c.JSON(ErroRespnse{MESSAGE: "Without token"})
			c.Status(401)
		} else {
			c.Next()
		}
	})

	database := db.Connect()
	redisC := db.Redisdb()
	GetUser := db.GetUserId()

	routes.Routes(app, database, redisC, GetUser)
	app.Listen(3000)
}
