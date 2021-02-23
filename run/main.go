package main

import (
	"os"

	railsassets "github.com/avarteqgmbh/rvm-rails-assets"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/fs"
	"github.com/paketo-buildpacks/packit/pexec"
)

func main() {
	logEmitter := railsassets.NewLogEmitter(os.Stdout)

	packit.Run(
		railsassets.Detect(railsassets.NewGemfileParser()),
		railsassets.Build(
			railsassets.NewPrecompileProcess(
				pexec.NewExecutable("bash"),
				logEmitter,
			),
			fs.NewChecksumCalculator(),
			railsassets.NewDirectorySetup(),
			railsassets.NewEnvironment(),
			logEmitter,
			chronos.DefaultClock,
		),
	)
}
