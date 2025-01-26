package main

import (
	"github.com/NdoleStudio/httpmock/pkg/di"
)

func main() {
	di.LoadEnv("../../.env")

	container := di.NewLiteContainer()
	container.InitializeTraceProvider()
	logger := container.Logger()

	logger.Info("Starting the application")
}
