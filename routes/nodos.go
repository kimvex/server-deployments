package routes

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber"
)

//BodyNodes Struct for body of nodes
type BodyNodes struct {
	NodeName string `json:"node_name"`
	Version  string `json:"version"`
}

//NodosData struct for response of nodes
type NodosData struct {
	NameNodo sql.NullString
	Version  sql.NullString
}

//ResponseSuccessDataJSON response of nodes
type ResponseSuccessDataJSON struct {
	Nodos []NodosData `json:"nodos"`
}

// Nodos is a function for adding nodos and get nodos
func Nodos() {
	app.Post("/nodos", ValidateRoute, func(c *fiber.Ctx) {
		var body BodyNodes
		if err := c.BodyParser(&body); err != nil {
			fmt.Println(err)
		}

		idNode, errorInsert := sq.Insert("nodos").
			Columns("name_nodo", "version").
			Values(body.NodeName, body.Version).
			RunWith(database).
			Exec()

		idNode.LastInsertId()

		if errorInsert != nil {
			fmt.Println(errorInsert)
			ErrorI := ErrorResponse{Message: "No se pudo agregar el host"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		success := SuccessResponse{MESSAGE: "El nodo no pudo ser agregado"}
		c.JSON(success)
	})

	app.Get("/nodos", ValidateRoute, func(c *fiber.Ctx) {
		var nodosData NodosData
		var nodeList []NodosData

		fmt.Println("Si es aquu")

		nodes, err := sq.Select("name_nodo", "version").
			From("nodos").
			RunWith(database).
			Query()

		if err != nil {
			fmt.Println(err)
			ErrorI := ErrorResponse{Message: "Ocurrio un error con los nodos"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		for nodes.Next() {
			_ = nodes.Scan(&nodosData.NameNodo, &nodosData.Version)
			nodeList = append(nodeList, nodosData)
		}

		var response ResponseSuccessDataJSON
		response.Nodos = []NodosData{}
		if nodeList != nil {
			response.Nodos = nodeList
		}
		fmt.Println(response)

		c.JSON(response)
	})
}
