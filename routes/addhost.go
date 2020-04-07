package routes

import (
	"database/sql"
	"fmt"
	"log"

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

// BodyConnect struct for body of request
type BodyConnect struct {
	HOST string `json:"host"`
}

//ResponseSuccessJSON This structure is for send json custom
// type ResponseSuccessJSON struct {
// 	Hosts []HostResponse `json:"hosts"`
// }
type ResponseSuccessJSON struct {
	Hosts []string `json:"hosts"`
}

// HostData struct of data get host
type HostData struct {
	Host        sql.NullString
	CreateAt    sql.NullString
	ServiceName sql.NullString
	Repository  sql.NullString
	Path        sql.NullString
	NameNodo    sql.NullString
	Version     sql.NullString
}

type ResponseHostData struct {
	Host        string `json:"host"`
	CreateAt    string `json:"create_at"`
	ServiceName string `json:"service_name"`
	Repository  string `json:"respository"`
	Path        string `json:"path"`
	NameNodo    string `json:"name_nodo"`
	Version     string `json:"version"`
}

/*AddHost - Function of cotainer for AddHost*/
func AddHost() {
	app.Post("/addhost", ValidateRoute, func(c *fiber.Ctx) {
		var hostToadd Host
		if err := c.BodyParser(&hostToadd); err != nil {
			fmt.Println(err)
		}

		idFil, errorInsert := sq.Insert("hosts").
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

	app.Get("/hosts", ValidateRoute, func(c *fiber.Ctx) {
		var hostsArr HostResponse
		var hostList []string

		host, err := sq.Select("host").
			From("hosts").
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

	app.Get("/host/:host_id", ValidateRoute, func(c *fiber.Ctx) {
		var hostDataSet HostData
		hostID := c.Params("host_id")

		sql, _, _ := sq.
			Select("host", "create_at", "service_name", "repository", "path", "name_nodo", "version").
			From("hosts").
			LeftJoin("services on hosts.hosts_id=services.hosts_id").
			LeftJoin("nodos on nodos.nodo_id=services.nodo_id").
			Where(sq.Eq{"hosts_id.host_id": hostID}).
			ToSql()

		fmt.Println(sql)
		err := sq.
			Select("host", "create_at", "service_name", "repository", "path", "name_nodo", "version").
			From("hosts").
			LeftJoin("services on hosts.hosts_id=services.hosts_id").
			LeftJoin("nodos on nodos.nodo_id=services.nodo_id").
			Where(sq.Eq{"hosts.hosts_id": hostID}).
			RunWith(database).
			QueryRow().
			Scan(&hostDataSet.Host, &hostDataSet.CreateAt, &hostDataSet.ServiceName, &hostDataSet.Repository, &hostDataSet.Path, &hostDataSet.NameNodo, &hostDataSet.Version)
			// Scan(&hostDataSet.Host.String, &hostDataSet.CreateAt.String, &hostDataSet.ServiceName.String, &hostDataSet.Repository.String, &hostDataSet.Path.String, &hostDataSet.NameNodo.String, &hostDataSet.Version.String)
		if err != nil {
			fmt.Println(err)
			ErrorI := ErrorResponse{Message: "Ocurrio un error con los hosts"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		c.JSON(ResponseHostData{
			Host:        hostDataSet.Host.String,
			CreateAt:    hostDataSet.CreateAt.String,
			ServiceName: hostDataSet.ServiceName.String,
			Repository:  hostDataSet.Repository.String,
			Path:        hostDataSet.Path.String,
			NameNodo:    hostDataSet.NameNodo.String,
			Version:     hostDataSet.Version.String,
		})
	})
}
