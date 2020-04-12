package routes

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber"
)

type BodyCommand struct {
	Command string `json:"command"`
}

//Commands namespace
func Commands() {
	app.Post("/command", ValidateRoute, func(c *fiber.Ctx) {
		var body BodyCommand
		if err := c.BodyParser(&body); err != nil {
			fmt.Println(err)
		}

		_, errorInsert := sq.Insert("commands").
			Columns("command").
			Values(body.Command).
			RunWith(database).
			Exec()

		if errorInsert != nil {
			fmt.Println(errorInsert)
			ErrorI := ErrorResponse{Message: "No se pudo agregar el host"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		success := SuccessResponse{MESSAGE: "Comando agregado correctamente"}
		c.JSON(success)
	})
}
