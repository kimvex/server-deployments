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

type EnvSQL struct {
	EnvName  sql.NullString `json:"variable_name"`
	EnvValue sql.NullString `json:"variable_value"`
}

type EnvString struct {
	EnvName  *string `json:"variable_name"`
	EnvValue *string `json:"variable_value"`
}

type ResponseServices struct {
	ServiceID   sql.NullInt32
	ServiceName sql.NullString
	Repository  sql.NullString
	Path        sql.NullString
	NameNodo    sql.NullString
	Version     sql.NullString
}

type ResponseCompleteService struct {
	ServiceID   *int32  `json:"service_id"`
	ServiceName *string `json:"service_name"`
	Repository  *string `json:"repository"`
	Path        *string `json:"path"`
	NameNodo    *string `json:"name_nodo"`
	Version     *string `json:"version"`
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
		arrLitsCommands := []*string{}
		serviceID := c.Params("service_id")
		nodeID := c.Params("node_id")

		nodoResult, err := sq.Select(
			"name_nodo",
			"version",
		).
			From("services").
			LeftJoin("nodos on services.nodo_id=nodos.nodo_id").
			Where(sq.Eq{"services.nodo_id": nodeID, "service_id": serviceID}).
			RunWith(database).
			Query()

		if err != nil {
			fmt.Println(err, "El error")
			ErrorI := ErrorResponse{Message: "Ocurrio al obtener el nodo"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		for nodoResult.Next() {
			_ = nodoResult.Scan(&nodoSQL.NodeName, &nodoSQL.Version)
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
		var envs EnvSQL
		var listEnvSQL []EnvSQL
		var envStr EnvString
		var arrListEnv []EnvString

		serviceID := c.Params("service_id")
		nodeID := c.Params("node_id")

		errCD := sq.Select("path", "repository", "host").
			From("services").
			LeftJoin("hosts on services.hosts_id=hosts.hosts_id").
			Where(sq.Eq{"service_id": serviceID, "nodo_id": nodeID}).
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

		envGetSQL, errEnv := sq.Select("variable_name", "variable_value").
			From("variables_enviroment").
			Where(sq.Eq{"service_id": serviceID}).
			RunWith(database).
			Query()

		if errEnv != nil {
			fmt.Println(errEnv)
			ErrorI := ErrorResponse{Message: "Ocurrio un error con las variables"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		for envGetSQL.Next() {
			_ = envGetSQL.Scan(&envs.EnvName, &envs.EnvValue)

			listEnvSQL = append(listEnvSQL, envs)
		}

		for i := 0; i < len(listEnvSQL); i++ {
			envStr.EnvName = &listEnvSQL[i].EnvName.String
			envStr.EnvValue = &listEnvSQL[i].EnvValue.String
			arrListEnv = append(arrListEnv, envStr)
		}

		fmt.Println(commandsD)

		ExecuteDeploy(&commandsD.Path.String, accequibleCommands, arrListEnv, &commandsD.Repository.String, &commandsD.Host.String)

		c.JSON(ResponseDeploy{R: "Success"})
	})

	app.Get("/service/:service_id", ValidateRoute, func(c *fiber.Ctx) {
		serviceID := c.Params("service_id")
		var responseSQL ResponseServices
		var response ResponseCompleteService

		service, errService := sq.Select(
			"services.service_id",
			"service_name",
			"repository",
			"path",
			"name_nodo",
			"version",
		).
			From("services").
			LeftJoin("nodos on services.nodo_id = nodos.nodo_id").
			Where(sq.Eq{"service_id": serviceID}).
			RunWith(database).
			Query()

		if errService != nil {
			fmt.Println(errService)
			ErrorI := ErrorResponse{Message: "Ocurrio un problema al obtener la informacion del servicio."}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		for service.Next() {
			_ = service.Scan(
				&responseSQL.ServiceID,
				&responseSQL.ServiceName,
				&responseSQL.Repository,
				&responseSQL.Path,
				&responseSQL.NameNodo,
				&responseSQL.Version,
			)

			fmt.Println(responseSQL.ServiceID)

			response.ServiceID = &responseSQL.ServiceID.Int32
			response.ServiceName = &responseSQL.ServiceName.String
			response.Repository = &responseSQL.Repository.String
			response.Path = &responseSQL.Path.String
			response.NameNodo = &responseSQL.NameNodo.String
			response.Version = &responseSQL.Version.String
		}

		c.JSON(response)
	})
}
