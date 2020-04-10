package routes

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber"
)

type DataService struct {
	ServiceName string `json:"service_name"`
	Repository  string `json:"repository"`
	Path        string `json:"path"`
	HostsID     int    `json:"hosts_id"`
	NodoID      int    `json:"nodo_id"`
}

//Services Namespace for endpoint of services
func Services() {
	app.Post("/service", ValidateRoute, func(c *fiber.Ctx) {
		var service DataService
		if err := c.BodyParser(&service); err != nil {
			fmt.Println(err)
		}

		_, errorInsert := sq.Insert("services").
			Columns(
				"service_name",
				"repository",
				"path",
				"hosts_id",
				"nodo_id",
			).
			Values(
				service.ServiceName,
				service.Repository,
				service.Path,
				service.HostsID,
				service.NodoID,
			).
			RunWith(database).
			Exec()

		if errorInsert != nil {
			fmt.Println(errorInsert)
			ErrorI := ErrorResponse{Message: "No se pudo guardar el servicio"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		success := SuccessResponse{MESSAGE: "Servicio agregado correctamente"}
		c.JSON(success)
	})

	app.Post("/service/deploy", ValidateRoute, func(c *fiber.Ctx) {

	})
}
