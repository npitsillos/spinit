package deploy

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/npitsillos/spinit/config"
	"github.com/npitsillos/spinit/pkg/build"
	"github.com/npitsillos/spinit/pkg/deploy"
	"github.com/npitsillos/spinit/pkg/load"
	"github.com/npitsillos/spinit/tests/integration"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	dir            string
	kubeConfigFile string
	appCfg         *config.AppConfig
	k8sClient      *integration.K8sClient
)

func Test_E2EDeployment(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Deploy Test Suite")
}

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

	err = integration.CopyKubeConfigFileToDefaultPath(kubeConfigFile)
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

	err = load.LoadImageToNodes("sample-app.tar", []string{"server", "agent"})
	Expect(err).To(BeNil())

	appCfg = &config.AppConfig{
		AppName: "sample-app",
		Build: &config.Build{
			Image:      "docker.io/library/sample-app",
			Version:    "latest",
			Dockerfile: "Dockerfile",
		},
	}

	k8sClient, err = integration.GetK8SClient(kubeConfigFile)
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("Deploy app to cluster", func() {
	When("an application is deployed to the cluster", func() {
		It("should create the correct resources", func() {
			err := deploy.CreateDeployment(appCfg, false, "sample-app")
			Expect(err).NotTo(HaveOccurred())

			namespace, err := k8sClient.CoreV1().Namespaces().Get(context.Background(), "sample-app", metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(namespace).ToNot(BeNil())
			Expect(namespace.Name).To(Equal("sample-app"))
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
