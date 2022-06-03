package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

var settings struct {
	Buildpacks struct {
		RvmCnb struct {
			Online string
		}
		RvmBundlerCnb struct {
			Online string
		}
		NodeEngine struct {
			Online string
		}
		RvmRailsAssets struct {
			Online string
		}
		Yarn struct {
			Online string
		}
		YarnInstall struct {
			Online string
		}
	}

	Buildpack struct {
		ID   string
		Name string
	}

	Config struct {
		RvmCnb        string `json:"rvm-cnb"`
		RvmBundlerCnb string `json:"rvm-bundler-cnb"`
		NodeEngine    string `json:"node-engine"`
		Yarn          string `json:"yarn"`
		YarnInstall   string `json:"yarn-install"`
	}

	Pack   occam.Pack
	Docker occam.Docker
}

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect

	root, err := filepath.Abs("./..")
	Expect(err).NotTo(HaveOccurred())

	file, err := os.Open("../integration.json")
	Expect(err).NotTo(HaveOccurred())

	Expect(json.NewDecoder(file).Decode(&settings.Config)).To(Succeed())

	file, err = os.Open("../buildpack.toml")
	Expect(err).NotTo(HaveOccurred())

	// Do not truncate Gomega matcher output
	// The buildpack output text can be large and we often want to see all of it.
	format.MaxLength = 0

	_, err = toml.DecodeReader(file, &settings)
	Expect(err).NotTo(HaveOccurred())
	Expect(file.Close()).To(Succeed())

	buildpackStore := occam.NewBuildpackStore()

	settings.Buildpacks.RvmRailsAssets.Online, err = buildpackStore.Get.
		WithVersion("1.2.3").
		Execute(root)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.RvmCnb.Online, err = buildpackStore.Get.
		Execute(settings.Config.RvmCnb)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.RvmBundlerCnb.Online, err = buildpackStore.Get.
		Execute(settings.Config.RvmBundlerCnb)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.NodeEngine.Online, err = buildpackStore.Get.
		Execute(settings.Config.NodeEngine)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.Yarn.Online, err = buildpackStore.Get.
		Execute(settings.Config.Yarn)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.YarnInstall.Online, err = buildpackStore.Get.
		Execute(settings.Config.YarnInstall)
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(30 * time.Second)

	settings.Pack = occam.NewPack().WithVerbose()
	settings.Docker = occam.NewDocker()

	suite := spec.New("Integration", spec.Parallel(), spec.Report(report.Terminal{}))
	suite("Rails5.0", testRails50)
	suite("Rails6.0", testRails60)
	suite("ReusingLayerRebuild", testReusingLayerRebuild)
	suite.Run(t)
}
