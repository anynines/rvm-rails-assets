package railsassets_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	railsassets "github.com/avarteqgmbh/rvm-rails-assets"
	"github.com/avarteqgmbh/rvm-rails-assets/fakes"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/scribe"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testPrecompileProcess(t *testing.T, context spec.G, it spec.S) {
	var Expect = NewWithT(t).Expect

	context("Execute", func() {
		var (
			workingDir string
			path       string
			executions []pexec.Execution
			executable *fakes.Executable

			precompileProcess railsassets.PrecompileProcess
		)

		it.Before(func() {
			var err error
			workingDir, err = ioutil.TempDir("", "working-dir")
			Expect(err).NotTo(HaveOccurred())

			executions = []pexec.Execution{}
			executable = &fakes.Executable{}
			executable.ExecuteCall.Stub = func(execution pexec.Execution) error {
				executions = append(executions, execution)

				return nil
			}

			path = os.Getenv("PATH")
			os.Setenv("PATH", "/some/bin")

			logger := scribe.NewLogger(bytes.NewBuffer(nil))

			precompileProcess = railsassets.NewPrecompileProcess(executable, logger)
		})

		it.After(func() {
			os.Setenv("PATH", path)

			Expect(os.RemoveAll(workingDir)).To(Succeed())
		})

		it("runs the bundle exec rake:precompile process", func() {
			err := precompileProcess.Execute(workingDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(executions).To(HaveLen(1))
			Expect(executions[0].Args).To(Equal([]string{
				"--login",
				"-c",
				"source profile.d/rvm && RAILS_ENV=production bundle exec rake assets:precompile assets:clean",
			}))
		})

		it("runs the bundle exec assets:precompile process when env variable DB_ADAPTER is set", func() {
			os.Setenv("DB_ADAPTER", "someadapter")
			os.Setenv("RAILS_ENV", "development")

			err := precompileProcess.Execute(workingDir)
			Expect(err).NotTo(HaveOccurred())

			Expect(executions).To(HaveLen(1))
			Expect(executions[0].Args).To(Equal([]string{
				"--login",
				"-c",
				"source profile.d/rvm && DB_ADAPTER=someadapter RAILS_ENV=development bundle exec rake assets:precompile assets:clean",
			}))
		})

		context("failure cases", func() {
			context("when bundle exec fails", func() {
				it.Before(func() {
					executable.ExecuteCall.Stub = func(execution pexec.Execution) error {
						return errors.New("bundle exec failed")
					}
				})
				it("prints the execution output and returns an error", func() {
					err := precompileProcess.Execute(workingDir)
					Expect(err).To(MatchError(ContainSubstring("failed to execute bundle exec")))
					Expect(err).To(MatchError(ContainSubstring("bundle exec failed")))
				})
			})
		})
	})
}
