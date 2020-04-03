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

/*HostsStruct slice of string*/
type HostsStruct []string

/*HostResponse structure of response hosts*/
type HostResponse struct {
	HOSTS string `json:"hosts"`
}

//ResponseSuccessJSON This structure is for send json custom
// type ResponseSuccessJSON struct {
// 	Hosts []HostResponse `json:"hosts"`
// }
type ResponseSuccessJSON struct {
	Hosts []string `json:"hosts"`
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

	app.Get("/hosts", func(c *fiber.Ctx) {
		var hostsArr HostResponse
		var hostList []string

		host, err := sq.Select("host").
			From("hots").
			Where(sq.Eq{"user_id": userID}).
			RunWith(database).
			Query()

		if err != nil {
			fmt.Println(err)
			ErrorI := ErrorResponse{Message: "Ocurrio un error con los hosts"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		for host.Next() {
			_ = host.Scan(&hostsArr.HOSTS)
			hostList = append(hostList, hostsArr.HOSTS)
		}

		var response ResponseSuccessJSON
		response.Hosts = hostList

		c.JSON(response)
	})
}
