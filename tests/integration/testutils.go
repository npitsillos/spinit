package integration

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	NETWORK_PREFIX = "10.10.10"
	START_IP       = 100
)

func init() {
	os.Setenv("E2E_NETWORK_PREFIX", NETWORK_PREFIX)
	os.Setenv("E2E_START_IP", strconv.Itoa(START_IP))
}

type NodeError struct {
	Cmd  string
	Node string
	Err  error
}

type K8sClient struct {
	*kubernetes.Clientset
}

func (ne *NodeError) Error() string {
	return fmt.Sprintf("failed creating cluster: %s: %v", ne.Cmd, ne.Err)
}

func (ne *NodeError) Unwrap() error {
	return ne.Err
}

func newNodeError(cmd, node string, err error) *NodeError {
	return &NodeError{
		Cmd:  cmd,
		Node: node,
		Err:  err,
	}
}

func CreateCluster(serverNodeName, agentNodeName string) error {

	cmd := fmt.Sprintf("vagrant up %s &> vagrant.log", serverNodeName)

	if _, err := RunCommand(cmd); err != nil {
		return newNodeError(cmd, serverNodeName, err)
	}

	cmd = fmt.Sprintf("vagrant up %s &> vagrant.log", agentNodeName)
	if _, err := RunCommand(cmd); err != nil {
		return newNodeError(cmd, agentNodeName, err)
	}
	time.Sleep(5 * time.Second)

	return nil
}

func DestroyCluster() error {
	if _, err := RunCommand("vagrant destroy -f"); err != nil {
		return err
	}
	return os.Remove("vagrant.log")
}

func GenKubeConfigFile(serverName, serverIP string) (string, error) {
	kubeConfigFile := fmt.Sprintf("kubeconfig-%s", serverName)
	cmd := fmt.Sprintf("vagrant scp %s:/etc/rancher/k3s/k3s.yaml ./%s", serverName, kubeConfigFile)
	_, err := RunCommand(cmd)
	if err != nil {
		return "", err
	}

	kubeConfig, err := os.ReadFile(kubeConfigFile)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`(?m)==> vagrant:.*\n`)
	modifiedKubeConfig := re.ReplaceAllString(string(kubeConfig), "")

	modifiedKubeConfig = strings.Replace(modifiedKubeConfig, "127.0.0.1", serverIP, 1)
	if err := os.WriteFile(kubeConfigFile, []byte(modifiedKubeConfig), 0644); err != nil {
		return "", err
	}

	if err := os.Setenv("E2E_KUBECONFIG", kubeConfigFile); err != nil {
		return "", err
	}
	return kubeConfigFile, nil
}

func GetK8SClient(kubeConfigFile string) (*K8sClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile)
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &K8sClient{clientSet}, nil
}

func CopyKubeConfigFileToDefaultPath(kubeConfigFile string) error {
	home := homedir.HomeDir()
	kubeDir := filepath.Join(home, ".kube")
	err := os.Mkdir(kubeDir, 0o700)
	if err != nil {
		return err
	}
	kubeConfig, err := os.ReadFile(kubeConfigFile)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(kubeDir, "config"), []byte(kubeConfig), 0644)
}

func GetJournalLogs(node string) (string, error) {
	cmd := "journalctl -u k3s* --no-pager"
	return RunCmdOnNode(cmd, node)
}

func GetVagrantLog(cErr error) string {
	var nodeErr *NodeError
	nodeJournal := ""
	if errors.As(cErr, &nodeErr) {
		nodeJournal, _ = GetJournalLogs(nodeErr.Node)
		nodeJournal = "\nNode Journal Logs:\n" + nodeJournal
	}

	log, err := os.Open("vagrant.log")
	if err != nil {
		return err.Error()
	}
	bytes, err := io.ReadAll(log)
	if err != nil {
		return err.Error()
	}

	return string(bytes) + nodeJournal
}

func RunCmdOnNode(cmd string, node string) (string, error) {
	runCmd := "vagrant ssh " + node + " -c \"sudo GOCOVERDIR=/tmp/k3scov " + cmd + "\""
	out, err := RunCommand(runCmd)
	out = strings.ReplaceAll(out, "[fog][WARNING] Unrecognized arguments: libvirt_ip_command\n", "")
	if err != nil {
		return out, fmt.Errorf("failed to run command: %s on node %s: %s, %v", cmd, node, out, err)
	}
	return out, nil
}

func RunCommand(cmd string) (string, error) {
	c := exec.Command("bash", "-c", cmd)
	if kc, ok := os.LookupEnv("E2E_KUBECONFIG"); ok {
		c.Env = append(os.Environ(), "KUBECONFIG="+kc)
	}
	out, err := c.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("failed to run command: %s, %v", cmd, err)
	}
	return string(out), nil
}

func CreateTempProjectDir() (string, error) {
	dir, err := os.MkdirTemp(".", "sample-app")
	if err != nil {
		return "", err
	}

	file, err := os.Create(filepath.Join(dir, "Dockerfile"))
	if err != nil {
		return "", err
	}

	if _, err := file.WriteString("FROM python:3.13.0-bullseye"); err != nil {
		return "", err
	}

	return dir, nil
}

func CreateSSHKeyPair() error {

	privKeyFile := "test_key"
	pubKeyFile := "test_key.pub"
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	if err := privKey.Validate(); err != nil {
		return err
	}

	pubKey, err := ssh.NewPublicKey(&privKey.PublicKey)
	if err != nil {
		return err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(pubKey)

	privKeyBytes := encodePrivateKey(privKey)

	if err := writeToFile(privKeyFile, privKeyBytes); err != nil {
		return err
	}

	if err := writeToFile(pubKeyFile, pubKeyBytes); err != nil {
		return err
	}

	privKeyAbsPath, err := filepath.Abs(privKeyFile)
	if err != nil {
		return err
	}
	if err := os.Setenv("E2E_PRIV_KEY_PATH", privKeyAbsPath); err != nil {
		return err
	}

	pubKeyAbsPath, err := filepath.Abs(pubKeyFile)
	if err != nil {
		return err
	}
	if err := os.Setenv("E2E_PUB_KEY_PATH", pubKeyAbsPath); err != nil {
		return err
	}

	return nil
}

func encodePrivateKey(privateKey *rsa.PrivateKey) []byte {
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	return pem.EncodeToMemory(&privBlock)
}

func writeToFile(file string, keyBytes []byte) error {
	return os.WriteFile(file, keyBytes, 0600)
}
