package main

import (
	"os"

	_ "github.com/NdoleStudio/httpmock/docs"
	"github.com/NdoleStudio/httpmock/pkg/di"
)

// Version is the version of the API
var Version string

// @title       HTTP Mock API
// @version     1.0
// @description Backend HTTP Mock API server.
//
// @contact.name  Acho Arnold
// @contact.email arnold@httpmock.dev
//
// @license.name AGPL-3.0
// @license.url  https://raw.githubusercontent.com/NdoleStudio/httpmock/main/LICENSE
//
// @host     localhost:8000
// @schemes  http https
//
// @securitydefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if len(os.Args) == 1 {
		di.LoadEnv()
	}

	container := di.NewContainer(Version, os.Getenv("GCP_PROJECT_ID"))
	container.Logger().Info(container.App().Listen(":8000").Error())
}
