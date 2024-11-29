package load

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/npitsillos/spinit/config"
	"github.com/npitsillos/spinit/pkg/sshclient"
	"github.com/spf13/viper"
)

var (
	PERMISSIONS    = "0666"
	LOAD_IMAGE_CMD = "sudo ctr images import"
)

func LoadImageToNodes(image string, nodes []string) error {

	cfg := viper.Get("config").(*config.Config)
	for _, node := range cfg.Nodes {
		if !slices.Contains(nodes, node.Name) {
			continue
		}
		nodeSSHClient, err := sshclient.NewSSHClient(node, cfg.SSHKeyPath)
		if err != nil {
			return err
		}

		scpClient := scp.NewConfigurer(fmt.Sprintf("%s:%d", node.IP, sshclient.PORT), nodeSSHClient.Config).SSHClient(nodeSSHClient.Client).Create()

		defer scpClient.Close()

		err = scpClient.Connect()
		if err != nil {
			return err
		}

		f, _ := os.Open(image)
		err = scpClient.CopyFromFile(context.Background(), *f, fmt.Sprintf("/home/%s/%s", node.Username, image), PERMISSIONS)

		if err != nil {
			return err
		}

		if err := nodeSSHClient.RunCommands([]string{fmt.Sprintf("%s %s", LOAD_IMAGE_CMD, image), fmt.Sprintf("rm %s", image)}); err != nil {
			return err
		}
	}
	return nil
}
