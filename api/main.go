package main

import (
	"os"

	"github.com/NdoleStudio/httpmock/docs"

	"github.com/gofiber/fiber/v2"

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
// @host     api.httpmock.dev
// @schemes  https
//
// @securitydefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if len(os.Args) == 1 {
		di.LoadEnv()
	}

	if Version != "" {
		docs.SwaggerInfo.Version = Version
	}

	container := di.NewContainer(Version, os.Getenv("GCP_PROJECT_ID"))

	app := container.App()
	err := make(chan error, 1)

	go serveHTTP(app, err)

	if os.Getenv("TLS_CERT_FILE") != "" {
		go serveHTTPS(app, err)
	}

	container.Logger().Error(<-err)
}

func serveHTTP(app *fiber.App, err chan<- error) {
	err <- app.Listen(":8000")
}

func serveHTTPS(app *fiber.App, err chan<- error) {
	err <- app.ListenTLS(":8443", os.Getenv("TLS_CERT_FILE"), os.Getenv("TLS_KEY_FILE"))
}
