package main

import (
	"template/internal/app"
)

// @title		Test API
// @version	1.0
// @BasePath	/
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	app.Start()
}
