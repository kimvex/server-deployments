package routes

import (
	"fmt"
	"reflect"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber"
	"golang.org/x/crypto/bcrypt"
)

func Register() {
	fmt.Println(userID)
	app.Post("/register", func(c *fiber.Ctx) {
		fmt.Println(c.Body())
		fmt.Println(c.Body("email"), c.Body("Email"))
		var bodyRequest BodyRequest

		if err := c.BodyParser(&bodyRequest); err != nil {
			fmt.Println(err)
		}

		exist := sq.Select("email").From("usersS").Where(sq.Eq{"email": bodyRequest.Email})

		var emailExist AlreadyRegistered

		errorResponse := exist.RunWith(database).QueryRow().Scan(&emailExist.Email)

		if errorResponse == nil {
			textReason := fmt.Sprintf("El correo %s ya esta registrado", emailExist.Email)
			errorOf := ErroRespnse{MESSAGE: textReason}
			c.JSON(errorOf)
			c.SendStatus(401)
			return
		}

		if bodyRequest.Password != bodyRequest.Same_password {
			error := ErroRespnse{MESSAGE: "La contraseña no coincide"}
			c.JSON(error)
			c.SendStatus(401)
			return
		}

		password_hash, _ := bcrypt.GenerateFromPassword([]byte(bodyRequest.Password), 14)

		id, errorInsert := sq.Insert("usersS").
			Columns("email", "password").
			Values(bodyRequest.Email, string(password_hash)).
			RunWith(database).
			Exec()

		fmt.Println(reflect.TypeOf(id), errorInsert, "con el valor de los terroistas")
		r, _ := id.LastInsertId()
		fmt.Println(r, errorInsert, "con el valor de los terroistas")

		if errorInsert != nil {
			fmt.Println(errorInsert)
			errorI := ErroRespnse{MESSAGE: "No se puedo registrar al usuario"}
			c.JSON(errorI)
			c.SendStatus(400)
			return
		}

		success := SuccessResponse{MESSAGE: "El usuario se ha creado correctamente"}
		c.JSON(success)
	})
}
