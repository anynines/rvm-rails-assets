package main

import (
	"os"

	railsassets "github.com/avarteqgmbh/rvm-rails-assets"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/fs"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/scribe"
)

func main() {
	logger := scribe.NewLogger(os.Stdout)

	packit.Run(
		railsassets.Detect(railsassets.NewGemfileParser()),
		railsassets.Build(
			railsassets.NewPrecompileProcess(
				pexec.NewExecutable("bash"),
				logger,
			),
			fs.NewChecksumCalculator(),
			railsassets.NewDirectorySetup(),
			logger,
			chronos.DefaultClock,
		),
	)
}
