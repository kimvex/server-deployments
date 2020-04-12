package routes

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber"
)

//DataService struct of dataservice
type DataService struct {
	ServiceName string `json:"service_name"`
	Repository  string `json:"repository"`
	Path        string `json:"path"`
	HostsID     int    `json:"hosts_id"`
	NodoID      int    `json:"nodo_id"`
}

//NodoSql structure get nodo sql
type NodoSql struct {
	NodeName  sql.NullString `json:"nodo_name"`
	Version   sql.NullString `json:"version"`
	CommandID sql.NullString `json:"command_id"`
}

//Example for null int
//CommandID sql.NullInt32  `json:"command_id"`

type CommadsSQL struct {
	Comm sql.NullString
}

//CommandsSql structure of commands
type CommandsSql struct {
	Command *string `json:"command"`
}

type reponseIDs struct {
	NodeName *string   `json:"nodo_name"`
	Version  *string   `json:"version"`
	Commands []*string `json:"commands"`
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

	app.Get("/service/:service_id/node/:node_id", ValidateRoute, func(c *fiber.Ctx) {
		var nodoSql NodoSql
		var listIds []string
		var listIdsArr []int
		var commandFromSQL CommadsSQL
		var listCommandFromSQL []CommadsSQL
		var listCommand []*string
		var responseData reponseIDs
		// var nodeList []NodoSql
		// var commands CommandsSql
		serviceID := c.Params("service_id")
		nodeId := c.Params("node_id")

		err := sq.Select(
			"nodos.name_nodo",
			"nodos.version",
			"group_concat(command_id) as commands_id",
		).
			From("commands_node").
			LeftJoin("nodos on commands_node.nodo_id=nodos.nodo_id").
			Where(sq.Eq{"commands_node.nodo_id": nodeId, "service_id": serviceID}).
			GroupBy("service_id, name_nodo, version").
			RunWith(database).
			QueryRow().
			Scan(&nodoSql.NodeName, &nodoSql.Version, &nodoSql.CommandID)

		if err != nil {
			fmt.Println(err)
			ErrorI := ErrorResponse{Message: "Ocurrio al obtener el nodo"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		listIds = strings.Split(nodoSql.CommandID.String, ",")
		for i := 0; i < len(listIds); i++ {
			intValue, errStr := strconv.Atoi(listIds[i])
			if errStr != nil {
				fmt.Println(errStr, "Error al castear el id a int")
			}

			listIdsArr = append(listIdsArr, intValue)
		}

		fmt.Println(listIdsArr)

		commandSqlStr, errIds := sq.Select("command").
			From("commands").
			Where(sq.Eq{"command_id": []int{1, 2}}).
			RunWith(database).
			Query()

		if errIds != nil {
			fmt.Println(errIds)
			ErrorI := ErrorResponse{Message: "Ocurrio un error con los commandos"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}
		for commandSqlStr.Next() {
			_ = commandSqlStr.Scan(&commandFromSQL.Comm)

			listCommandFromSQL = append(listCommandFromSQL, commandFromSQL)
		}

		fmt.Println("el-", listCommandFromSQL)

		for i := 0; i < len(listCommandFromSQL); i++ {
			listCommand = append(listCommand, &listCommandFromSQL[i].Comm.String)
		}

		fmt.Println(listCommand, "ella es cayaita")

		responseData.NodeName = &nodoSql.NodeName.String
		responseData.Version = &nodoSql.Version.String
		responseData.Commands = listCommand

		c.JSON(responseData)
	})

	app.Post("/service/deploy", ValidateRoute, func(c *fiber.Ctx) {

	})
}
