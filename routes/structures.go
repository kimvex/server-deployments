package routes

import (
	"database/sql"

	"github.com/dgrijalva/jwt-go"
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
	Services()
	Nodos()
	Commands()

	app.Get("/", func(c *fiber.Ctx) {
		var respuesta SuccessResponse
		respuesta.MESSAGE = "Estas en la raiz del proyecto"
		c.JSON(respuesta)
	})
}

//ValidateRoute of token
func ValidateRoute(c *fiber.Ctx) {
	if c.Get("token") != "" {
		token, err := jwt.Parse(c.Get("token"), func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if token.Valid {
			c.Next()
			return
		}

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

	c.JSON(ErroRespnse{MESSAGE: "Without token"})
	c.Status(401)
	return
}
