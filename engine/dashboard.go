package engine

type User struct {
	Name string
}

type UserId struct {
	UserId string `json:"user_id"`
}

func DashBoard() {
	// routes.Register(app)
	// routes.Login(app)
	// routes.AddHost(app)
	// app.Get("/", func(c *fiber.Ctx) {
	// 	// hash, _ := bcrypt.GenerateFromPassword([]byte("elsss"), 14)
	// 	// hashC := bcrypt.CompareHashAndPassword([]byte(string(hash)), []byte("elsss."))
	// 	// hashCS := bcrypt.CompareHashAndPassword([]byte(string(hash)), []byte("elsss"))

	// 	users := sq.Select("*").From("service_type")

	// 	sql, _, _ := users.ToSql()
	// 	r := db.Exec(database, sql)
	// 	c.JSON(r)
	// })
}
