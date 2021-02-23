package railsassets

import (
	"os"

	"github.com/paketo-buildpacks/packit"
)

type Environment struct{}

func NewEnvironment() Environment {
	return Environment{}
}

func (Environment) Configure(launchEnv packit.Environment) error {
	if val, ok := os.LookupEnv("RAILS_ENV"); ok {
		launchEnv.Default("RAILS_ENV", val)
	} else {
		launchEnv.Default("RAILS_ENV", "production")
		os.Setenv("RAILS_ENV", "production")
	}

	launchEnv.Default("RAILS_SERVE_STATIC_FILES", "true")

	return nil
}
