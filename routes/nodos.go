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
	NameNodo sql.NullString `json:"name_nodo"`
	Version  sql.NullString `json:"version"`
}

//NodoDataParse struct for parsing nodoData
type NodoDataParse struct {
	NameNodo *string `json:"name_nodo"`
	Version  *string `json:"version"`
}

//ResponseSuccessDataJSON response of nodes
type ResponseSuccessDataJSON struct {
	Nodos []NodoDataParse `json:"nodos"`
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

		success := SuccessResponse{MESSAGE: "El nodo se ah agregado correctamente"}
		c.JSON(success)
	})

	app.Get("/nodos", ValidateRoute, func(c *fiber.Ctx) {
		var nodosData NodosData
		var nodeList []NodosData
		var componentNodos NodoDataParse
		var listNodo []NodoDataParse

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
		response.Nodos = []NodoDataParse{}
		if nodeList != nil {
			for i := 0; i < len(nodeList); i++ {
				componentNodos.NameNodo = &nodeList[i].NameNodo.String
				componentNodos.Version = &nodeList[i].Version.String
				listNodo = append(listNodo, componentNodos)
			}
		}

		response.Nodos = listNodo
		c.JSON(response)
	})
}
