package routes

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber"
)

type BodyVariable struct {
	VariableName  string `json:"variable_name"`
	ValueVariable string `json:"value_variable"`
	ServiceID     int    `json:"service_id"`
}

func EnvoromentVariables() {
	app.Post("/variable", ValidateRoute, func(c *fiber.Ctx) {
		var body BodyVariable

		if err := c.BodyParser(&body); err != nil {
			fmt.Print(err)
		}

		_, errInsert := sq.Insert("variables_enviroment").
			Columns("variable_name", "variable_value", "service_id").
			Values(body.VariableName, body.ValueVariable, body.ServiceID).
			RunWith(database).
			Exec()

		if errInsert != nil {
			fmt.Println(errInsert)
			ErrorI := ErrorResponse{Message: "Error al agregar la variable"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		success := SuccessResponse{MESSAGE: "La variable se a agregado correctamente"}
		c.JSON(success)
	})
}
