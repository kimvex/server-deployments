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
	Host     sql.NullString
	CreateAt sql.NullString
}

type ServiceHostData struct {
	ServiceName sql.NullString
	Repository  sql.NullString
	Path        sql.NullString
	NameNodo    sql.NullString
	Version     sql.NullString
}

type ResponseServiceHost struct {
	ServiceName *string `json:"service_name"`
	Repository  *string `json:"respository"`
	Path        *string `json:"path"`
	NameNodo    *string `json:"name_nodo"`
	Version     *string `json:"version"`
}

//ResponseHostData structure of response
type ResponseHostData struct {
	Host     string                `json:"host"`
	CreateAt string                `json:"create_at"`
	Services []ResponseServiceHost `json:"services"`
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
		var serviceToadd ServiceHostData
		var listServices []ServiceHostData
		var componentServiceVal ResponseServiceHost
		var listServicesVal []ResponseServiceHost
		hostID := c.Params("host_id")

		err := sq.
			Select("host", "create_at").
			From("hosts").
			Where(sq.Eq{"hosts.hosts_id": hostID}).
			RunWith(database).
			QueryRow().
			Scan(&hostDataSet.Host, &hostDataSet.CreateAt)

		if err != nil {
			fmt.Println(err)
			ErrorI := ErrorResponse{Message: "Ocurrio un error con los hosts"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		services, errServices := sq.
			Select(
				"service_name",
				"repository",
				"path",
				"name_nodo",
				"version",
			).
			From("services").
			LeftJoin("nodos on nodos.nodo_id=services.nodo_id").
			RunWith(database).
			Query()

		if errServices != nil {
			fmt.Println(errServices)
			ErrorI := ErrorResponse{Message: "Ocurrio un error con los nodos"}
			c.JSON(ErrorI)
			c.SendStatus(400)
			return
		}

		for services.Next() {
			_ = services.Scan(&serviceToadd.ServiceName, &serviceToadd.Repository, &serviceToadd.Path, &serviceToadd.NameNodo, &serviceToadd.Version)
			listServices = append(listServices, serviceToadd)
		}

		for i := 0; i < len(listServices); i++ {
			fmt.Println(listServices[i].ServiceName.String)
			componentServiceVal.ServiceName = &listServices[i].ServiceName.String
			componentServiceVal.Repository = &listServices[i].Repository.String
			componentServiceVal.Path = &listServices[i].Path.String
			componentServiceVal.NameNodo = &listServices[i].NameNodo.String
			componentServiceVal.Version = &listServices[i].Version.String
			listServicesVal = append(listServicesVal, componentServiceVal)
		}

		c.JSON(ResponseHostData{
			Host:     hostDataSet.Host.String,
			CreateAt: hostDataSet.CreateAt.String,
			Services: listServicesVal,
		})
	})
}
