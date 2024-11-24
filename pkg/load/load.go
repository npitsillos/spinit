package load

import (
	"context"
	"fmt"
	"os"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/npitsillos/spinit/config"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

var (
	PORT           = 22
	PERMISSIONS    = "0666"
	LOAD_IMAGE_CMD = "sudo ctr images import"
)

func getSSHClient(node *config.Node, sshKeyPath string) (*ssh.Client, error) {
	clientConfig, _ := auth.PrivateKey(node.Username, sshKeyPath, ssh.InsecureIgnoreHostKey())
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", node.IP, PORT), &clientConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func LoadImageToNodes(image string, nodes []string) error {

	cfg := viper.Get("config").(*config.Config)
	// eg, context := errgroup.WithContext(appcontext.Context())

	for _, node := range cfg.Nodes {
		nodeSSHClient, err := getSSHClient(node, cfg.SSH)
		if err != nil {
			return err
		}
		defer nodeSSHClient.Close()
		client, err := scp.NewClientBySSH(nodeSSHClient)
		if err != nil {
			return err
		}
		defer client.Close()
		err = client.Connect()
		if err != nil {
			return err
		}

		f, _ := os.Open(image)
		err = client.CopyFromFile(context.Background(), *f, fmt.Sprintf("/home/%s/%s", node.Username, image), PERMISSIONS)

		if err != nil {
			return err
		}

		session, err := nodeSSHClient.NewSession()
		if err != nil {
			return err
		}
		defer session.Close()
		if err := session.Run(fmt.Sprintf("%s %s", LOAD_IMAGE_CMD, image)); err != nil {
			return err
		}

		if err := session.Run(fmt.Sprintf("rm %s", image)); err != nil {
			return err
		}
	}
	return nil
}
