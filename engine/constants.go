package engine

type Paths struct {
	Url    string
	Method string
}

func Con() []Paths {
	validate := []Paths{}
	validate = append(validate, Paths{Url: "/addhost", Method: "POST"})
	validate = append(validate, Paths{Url: "/hosts", Method: "GET"})
	validate = append(validate, Paths{Url: "/connect", Method: "POST"})
	validate = append(validate, Paths{Url: "/host/:host_id", Method: "GET"})
	// validate = append(validate, Paths{Url: "/host/:host_id/carros/:e", Method: "GET"})
	// validate = append(validate, Paths{Url: "/", Method: "GET"})
	return validate
}
