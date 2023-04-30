package bootstrap

import (
	"github.com/goravel/framework/foundation"
	"goravel/config"
	"os"
	"path/filepath"
)

func Boot() *foundation.Application {
	/*app := foundation.Application{}

	//Bootstrap the application
	app.Boot()

	//Bootstrap the config.
	config.Boot()*/
	app := foundation.NewApplication(getAppBasePath())
	config.Boot()
	if err := app.RegisterConfiguredProviders(); err != nil {
		panic(err)
	}
	app.Boot()

	return app
}

func getAppBasePath() string {
	basePath := os.Getenv("APP_BASE_PATH")
	if basePath == "" {
		workingDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		basePath = filepath.Dir(workingDir)
	}
	return basePath
}
