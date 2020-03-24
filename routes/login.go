package routes

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"golang.org/x/crypto/bcrypt"
)

/*Login function */
func Login() {
	singSecret := []byte("secret")
	app.Post("/login", func(c *fiber.Ctx) {
		var userLogin User
		if err := c.BodyParser(&userLogin); err != nil {
			fmt.Println(err)
		}

		var userSelected User

		errorGetUser := sq.Select("email", "password", "user_id").
			From("usersS").
			Where(sq.Eq{"email": userLogin.Email}).
			RunWith(database).
			QueryRow().
			Scan(&userSelected.Email, &userSelected.Password, &userSelected.UserID)

		if errorGetUser != nil {
			fmt.Println(errorGetUser)
			error := ErrorResponse{Message: "Ocurrio un error al obtener el usuario"}
			c.JSON(error)
			c.Status(401)
			return
		}

		compare := bcrypt.CompareHashAndPassword([]byte(string(userSelected.Password)), []byte(userLogin.Password))

		if compare == nil {

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"kmv": "kimvex",
				"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
			})

			tokenS, setErr := token.SignedString(singSecret)
			if setErr != nil {
				panic(setErr)
			}

			fmt.Println(tokenS)
			err := redisC.Set(tokenS, userSelected.UserID, 0).Err()
			if err != nil {
				panic(err)
			}
			c.JSON(TokenResponse{Token: tokenS})
			return
		}

		error := ErrorResponse{Message: "Usuario ó contraseña incorrectos"}
		c.JSON(error)
		c.Status(401)
	})
}
