package build

import (
	"os"
	"testing"

	"github.com/npitsillos/spinit/pkg/build"
	"github.com/npitsillos/spinit/tests/integration"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func Test_E2EDeployment(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Build Test Suite")
}

var (
	dir       string
	buildOpts *build.BuildOpt
)

var _ = BeforeSuite(func() {
	dir, err := integration.CreateTempProjectDir()
	Expect(err).NotTo(HaveOccurred())

	buildOpts = &build.BuildOpt{
		ProjectDir: dir,
		Name:       "sample-app",
		Tag:        "latest",
		Dockerfile: "Dockerfile",
	}
})

var _ = Describe("Build image", func() {
	When("a build command is issued", func() {
		It("should correctly build the image and store it in a .tar file", func() {
			err := build.BuildImage(buildOpts)
			Expect(err).NotTo(HaveOccurred())
			_, err = os.Stat("sample-app.tar")
			Expect(err).To(BeNil())
		})
	})
})

var _ = AfterSuite(func() {
	err := os.RemoveAll(dir)
	Expect(err).NotTo(HaveOccurred())
})
