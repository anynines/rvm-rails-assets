package railsassets_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"io/ioutil"

	railsassets "github.com/avarteqgmbh/rvm-rails-assets"
	"github.com/avarteqgmbh/rvm-rails-assets/fakes"
	"github.com/paketo-buildpacks/packit"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir    string
		gemfileParser *fakes.Parser
		detect        packit.DetectFunc
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		err = os.MkdirAll(filepath.Join(workingDir, "app", "assets"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(workingDir, "Gemfile"), []byte{}, 0600)
		Expect(err).NotTo(HaveOccurred())

		gemfileParser = &fakes.Parser{}

		detect = railsassets.Detect(gemfileParser)
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("when the Gemfile lists rails and the app/assets directory exists", func() {
		it.Before(func() {
			gemfileParser.ParseCall.Returns.HasRails = true
		})

		it("detects", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Plan).To(Equal(packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{
						Name: "rvm-rails-assets",
					},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "rvm-rails-assets",
						Metadata: railsassets.BuildPlanMetadata{
							Build: true,
						},
					},
					{
						Name: "rvm",
						Metadata: railsassets.BuildPlanMetadata{
							Build: true,
						},
					},
					{
						Name: "node",
						Metadata: railsassets.BuildPlanMetadata{
							Build: true,
						},
					},
				},
			}))
		})

		context("when the working directory contains a yarn.lock file", func() {
			it.Before(func() {
				Expect(ioutil.WriteFile(filepath.Join(workingDir, "yarn.lock"), nil, 0600)).To(Succeed())
			})

			it("detects with node_modules", func() {
				result, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Plan).To(Equal(packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{
						{
							Name: "rvm-rails-assets",
						},
					},
					Requires: []packit.BuildPlanRequirement{
						{
							Name: "rvm-rails-assets",
							Metadata: railsassets.BuildPlanMetadata{
								Build: true,
							},
						},
						{
							Name: "rvm",
							Metadata: railsassets.BuildPlanMetadata{
								Build: true,
							},
						},
						{
							Name: "node_modules",
							Metadata: railsassets.BuildPlanMetadata{
								Build: true,
							},
						},
					},
				}))
			})
		})
	})

	context("when the Gemfile does not list rails", func() {
		it.Before(func() {
			gemfileParser.ParseCall.Returns.HasRails = false
		})

		it("detect should fail with error", func() {
			_, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).To(MatchError(packit.Fail))
		})
	})

	context("when the app/assets directory does not exist", func() {
		it.Before(func() {
			gemfileParser.ParseCall.Returns.HasRails = true
			err := os.RemoveAll(filepath.Join(workingDir, "app/assets"))
			Expect(err).NotTo(HaveOccurred())
		})

		it("detect should fail with error", func() {
			_, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).To(MatchError(packit.Fail))
		})
	})

	context("failure cases", func() {
		context("when the gemfile parser fails", func() {
			it.Before(func() {
				gemfileParser.ParseCall.Returns.Err = errors.New("some-error")
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError("failed to parse Gemfile: some-error"))
			})
		})
	})
}
