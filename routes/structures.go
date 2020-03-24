package routes

import (
	"database/sql"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber"
)

/*Login*/
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserID   string `json:userId`
}

type TokenResponse struct {
	Token string `json:"token"`
}

/*Register*/
type BodyRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	Same_password string `json:"same_password"`
}

type AlreadyRegistered struct {
	Email string `json:"email"`
}

type UserBody struct {
	UserId int `json:"userId"`
}

/*Response to client*/
type ErroRespnse struct {
	MESSAGE string `json:"message"`
}

type SuccessResponse struct {
	MESSAGE string `json:"message_success"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

var (
	app      *fiber.App
	database *sql.DB
	redisC   *redis.Client
	userID   string
)

func Routes(App *fiber.App, Database *sql.DB, RedisCl *redis.Client, UserIDC string) {
	app = App
	database = Database
	redisC = RedisCl
	userID = UserIDC

	Register()
	Login()
	AddHost()

	app.Get("/", func(c *fiber.Ctx) {
		var respuesta SuccessResponse
		respuesta.MESSAGE = "Estas en la raiz del proyecto"
		c.JSON(respuesta)
	})
}
