package railsassets_test

import (
	"os"
	"testing"

	railsassets "github.com/avarteqgmbh/rvm-rails-assets"
	"github.com/paketo-buildpacks/packit"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testEnvironment(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect      = NewWithT(t).Expect
		launchEnv   packit.Environment
		environment railsassets.Environment
	)

	it.Before(func() {
		var err error
		Expect(err).NotTo(HaveOccurred())
		launchEnv = packit.Environment{}
	})

	it.After(func() {

	})

	context("Configure", func() {
		it("configures the environment variables if they are not set", func() {
			err := environment.Configure(launchEnv)
			Expect(err).NotTo(HaveOccurred())

			Expect(launchEnv).To(Equal(packit.Environment{
				"RAILS_ENV.default":                "production",
				"RAILS_SERVE_STATIC_FILES.default": "true",
			}))
		})

		context("when RAILS_ENV is set", func() {
			it.Before(func() {
				os.Setenv("RAILS_ENV", "some-rails-env-val")
				os.Setenv("RAILS_SERVE_STATIC_FILES", "some-rails-serve-static-files-val")
			})

			it.After(func() {
				os.Unsetenv("RAILS_ENV")
				os.Unsetenv("RAILS_SERVE_STATIC_FILES")
			})

			it("configures build envs using given value", func() {
				err := environment.Configure(launchEnv)
				Expect(err).NotTo(HaveOccurred())

				Expect(launchEnv["RAILS_ENV.default"]).To(Equal("some-rails-env-val"))
				Expect(launchEnv["RAILS_SERVE_STATIC_FILES.default"]).To(Equal("true"))
			})
		})
	})
}
