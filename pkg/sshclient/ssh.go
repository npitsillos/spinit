package sshclient

import (
	"fmt"
	"os"

	"github.com/npitsillos/spinit/config"
	"golang.org/x/crypto/ssh"
)

var (
	PORT = 22
)

type NodeSSHClient struct {
	Config *ssh.ClientConfig
	Client *ssh.Client
}

func getSSHClientConfig(username, privKeyPath string) (*ssh.ClientConfig, error) {
	sshConfig := &ssh.ClientConfig{
		User: username,
	}

	privKey, err := os.ReadFile(privKeyPath)
	if err != nil {
		return sshConfig, fmt.Errorf("unable to read private key file: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(privKey)
	if err != nil {
		return sshConfig, fmt.Errorf("failed parsing private key: %w", err)
	}

	sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))

	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey() // TODO: revisit this to add host key callback auth
	return sshConfig, nil
}

func NewSSHClient(node *config.Node, privKeyPath string) (*NodeSSHClient, error) {

	nodeSSHClient := &NodeSSHClient{}
	sshClientConfig, err := getSSHClientConfig(node.Username, privKeyPath)
	if err != nil {
		return nil, err
	}
	nodeSSHClient.Config = sshClientConfig
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", node.IP, PORT), sshClientConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to dial client: %w", err)
	}
	nodeSSHClient.Client = client
	return nodeSSHClient, nil
}

func (n *NodeSSHClient) RunCommands(commands []string) error {

	for _, command := range commands {
		session, err := n.Client.NewSession()
		if err != nil {
			return err
		}

		if err := session.Run(command); err != nil {
			return err
		}
		session.Close()
	}

	return nil
}
