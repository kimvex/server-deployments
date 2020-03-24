package routes

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber"
)

/*Host structure for insert to database*/
type Host struct {
	Hosttoadd string `json:"host"`
}

/*AddHost - Function of cotainer for AddHost*/
func AddHost() {
	app.Post("/addhost", func(c *fiber.Ctx) {
		var hostToadd Host
		if err := c.BodyParser(&hostToadd); err != nil {
			fmt.Println(err)
		}

		idFil, errorInsert := sq.Insert("hots").
			Columns("host", "user_id").
			Values(hostToadd.Hosttoadd, userID).
			RunWith(database).
			Exec()

		idFil.LastInsertId()

		if errorInsert != nil {
			fmt.Println(errorInsert)
			ErrorI := ErrorResponse{Message: "No se pudo agregar el host"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		success := SuccessResponse{MESSAGE: "El host se agrego correctamente"}
		c.JSON(success)
	})
}
