package routes

import (
	"database/sql"
	"fmt"

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

//NodoSQL structure get nodo sql
type NodoSQL struct {
	NodeName sql.NullString `json:"nodo_name"`
	Version  sql.NullString `json:"version"`
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

type CommandsDeploy struct {
	Path       sql.NullString `json:"path"`
	Repository sql.NullString `json:"repository"`
	Host       sql.NullString `json:"host"`
}

type CommandsResult struct {
	Command sql.NullString `json:"command"`
}

type ResponseDeploy struct {
	R string `json:"message"`
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
		var nodoSQL NodoSQL
		var commandFromSQL CommadsSQL
		var listCommandFromSQL []CommadsSQL
		var responseData reponseIDs
		var arrLitsCommands []*string
		serviceID := c.Params("service_id")
		nodeID := c.Params("node_id")

		err := sq.Select(
			"nodos.name_nodo",
			"nodos.version",
		).
			From("commands_node").
			LeftJoin("nodos on commands_node.nodo_id=nodos.nodo_id").
			Where(sq.Eq{"commands_node.nodo_id": nodeID, "service_id": serviceID}).
			RunWith(database).
			QueryRow().
			Scan(&nodoSQL.NodeName, &nodoSQL.Version)

		if err != nil {
			fmt.Println(err)
			ErrorI := ErrorResponse{Message: "Ocurrio al obtener el nodo"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		commandSQLStr, errIds := sq.Select("command").
			From("commands_node").
			LeftJoin("commands on commands_node.command_id=commands.command_id").
			Where(sq.Eq{"commands_node.nodo_id": nodeID, "service_id": serviceID}).
			OrderBy("index_position ASC").
			RunWith(database).
			Query()

		if errIds != nil {
			fmt.Println(errIds)
			ErrorI := ErrorResponse{Message: "Ocurrio un error con los commandos"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}
		for commandSQLStr.Next() {
			_ = commandSQLStr.Scan(&commandFromSQL.Comm)
			listCommandFromSQL = append(listCommandFromSQL, commandFromSQL)
		}

		for i := 0; i < len(listCommandFromSQL); i++ {
			arrLitsCommands = append(arrLitsCommands, &listCommandFromSQL[i].Comm.String)
		}

		responseData.NodeName = &nodoSQL.NodeName.String
		responseData.Version = &nodoSQL.Version.String
		responseData.Commands = arrLitsCommands

		c.JSON(responseData)
	})

	app.Post("/service/:service_id/node/:node_id/deploy", ValidateRoute, func(c *fiber.Ctx) {
		var commandsD CommandsDeploy
		var commandRe CommandsResult
		var listCommands []CommandsResult
		var accequibleCommands []*string
		serviceID := c.Params("service_id")
		nodeID := c.Params("node_id")

		errCD := sq.Select("path", "repository", "host").
			From("commands_node").
			LeftJoin("services on commands_node.service_id=services.service_id").
			LeftJoin("hosts on services.hosts_id=hosts.hosts_id").
			RunWith(database).
			QueryRow().
			Scan(&commandsD.Path, &commandsD.Repository, &commandsD.Host)

		if errCD != nil {
			fmt.Println(errCD)
			ErrorI := ErrorResponse{Message: "Ocurrio al obtener el nodo"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		cmd, errCMD := sq.Select("command").
			From("commands_node").
			LeftJoin("commands on commands_node.command_id=commands.command_id").
			Where(sq.Eq{"commands_node.nodo_id": nodeID, "service_id": serviceID}).
			OrderBy("index_position ASC").
			RunWith(database).
			Query()

		if errCMD != nil {
			fmt.Println(errCMD)
			ErrorI := ErrorResponse{Message: "Ocurrio al obtener el nodo"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		for cmd.Next() {
			_ = cmd.Scan(&commandRe.Command)
			fmt.Println(commandRe.Command.String, &commandRe.Command.String)
			listCommands = append(listCommands, commandRe)
		}

		for i := 0; i < len(listCommands); i++ {
			accequibleCommands = append(accequibleCommands, &listCommands[i].Command.String)
		}

		ExecuteDeploy(&commandsD.Path.String, accequibleCommands, &commandsD.Repository.String, &commandsD.Host.String)

		c.JSON(ResponseDeploy{R: "Success"})
	})
}
