package typical

// AppContext project contenxt
type AppContext struct {
	Name        string
	Version     string
	Description string
}

// Context instance of Context
var Context = AppContext{
	Name:        "Typical-RESTful-Server",
	Version:     "0.1.0",
	Description: "Example of typical and scalable RESTful API Server for Go",
}
