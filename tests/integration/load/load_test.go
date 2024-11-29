package load

import (
	"fmt"
	"os"
	"testing"

	"github.com/npitsillos/spinit/config"
	"github.com/npitsillos/spinit/pkg/build"
	"github.com/npitsillos/spinit/pkg/load"
	"github.com/npitsillos/spinit/tests/integration"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

func Test_E2EDeployment(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Deploy Test Suite")
}

var (
	dir            string
	kubeConfigFile string
)

var _ = BeforeSuite(func() {
	err := integration.CreateSSHKeyPair()
	Expect(err).NotTo(HaveOccurred())

	err = integration.CreateCluster("server", "agent")
	Expect(err).NotTo(HaveOccurred(), integration.GetVagrantLog(err))

	dir, err = integration.CreateTempProjectDir()
	Expect(err).NotTo(HaveOccurred())

	err = build.BuildImage(&build.BuildOpt{
		ProjectDir: dir,
		Name:       "sample-app",
		Tag:        "latest",
		Dockerfile: "Dockerfile",
	})
	Expect(err).NotTo(HaveOccurred())

	kubeConfigFile, err = integration.GenKubeConfigFile("server", fmt.Sprintf("%s.%d", integration.NETWORK_PREFIX, integration.START_IP))
	Expect(err).NotTo(HaveOccurred())

	cfg := &config.Config{
		SSHKeyPath: os.Getenv("E2E_PRIV_KEY_PATH"),
	}
	for idx, node := range []string{"server", "agent"} {

		cfg.Nodes = append(cfg.Nodes, &config.Node{
			Name:     node,
			IP:       fmt.Sprintf("%s.%d", integration.NETWORK_PREFIX, integration.START_IP+idx),
			Username: "vagrant",
		})
	}

	viper.Set("config", cfg)
})

var _ = Describe("Load image to cluster", func() {
	When("image is loaded to the cluster", func() {
		It("should be present in the container runtime", func() {
			err := load.LoadImageToNodes("sample-app.tar", []string{"server", "agent"})
			Expect(err).To(BeNil())
		})
	})
})

var _ = AfterSuite(func() {
	err := os.RemoveAll(dir)
	Expect(err).NotTo(HaveOccurred())

	Expect(os.Remove(os.Getenv("E2E_PRIV_KEY_PATH"))).To(Succeed())
	Expect(os.Remove(os.Getenv("E2E_PUB_KEY_PATH"))).To(Succeed())
	Expect(os.Remove("sample-app.tar")).To(Succeed())
	Expect(integration.DestroyCluster()).To(Succeed())
	Expect(os.Remove(kubeConfigFile)).To(Succeed())
})
