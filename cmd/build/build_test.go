package build

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/shlex"
	"github.com/npitsillos/spinit/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

func createTempProjectDir() (string, func(), error) {
	dir, err := os.MkdirTemp("", "sample-project")
	if err != nil {
		return "", nil, err
	}

	_, err = os.Create(filepath.Join(dir, "Dockerfile"))
	if err != nil {
		return "", nil, err
	}

	return dir, func() {
		os.RemoveAll(dir)
	}, nil
}

func TestBuildCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test build command")
}

var _ = Describe("Build Test", func() {
	var buildCmd *cobra.Command

	BeforeEach(func() {
		buildCmd = NewBuildCommand()
		buildCmd.Run = func(*cobra.Command, []string) {}
	})

	When("No arguments are passed to `build`", func() {
		It("the name flag should be set to the project's directory name", func() {
			buildCmd.SetArgs([]string{})
			err := buildCmd.Execute()
			Expect(err).ToNot(BeNil())
			name, err := buildCmd.Flags().GetString("name")
			Expect(err).To(BeNil())
			Expect(name).To(Equal("build"))
		})
	})

	When("Arguments are passed to `build`", func() {
		It("should return an error if no Dockerfile exists", func() {
			buildCmd.SetArgs([]string{"/tmp"})
			err := buildCmd.Execute()
			Expect(err).To(Equal(errors.ErrNoDockerFile))
		})

		It("should not return an error if a Dockerfile is present", func() {
			dir, rmDir, err := createTempProjectDir()
			Expect(err).To(BeNil())
			buildCmd.SetArgs([]string{dir})
			err = buildCmd.Execute()
			Expect(err).To(BeNil())
			rmDir()
		})

		It("should return an error when more arguments are sent", func() {
			args := "../../ test"
			argv, _ := shlex.Split(args)
			buildCmd.SetArgs(argv)
			err := buildCmd.Execute()
			Expect(err).ToNot(BeNil())
		})
	})
})
