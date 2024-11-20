package copy

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
	PORT        = 22
	PERMISSIONS = "0666"
)

func CopyImageToNodes(image string, nodes []string) error {

	cfg := viper.Get("config").(*config.Config)
	// eg, context := errgroup.WithContext(appcontext.Context())

	for _, node := range cfg.Nodes {
		clientConfig, _ := auth.PrivateKey(node.Username, cfg.SSH, ssh.InsecureIgnoreHostKey())
		client := scp.NewClient(fmt.Sprintf("%s:%d", node.IP, PORT), &clientConfig)
		err := client.Connect()
		if err != nil {
			return err
		}

		f, _ := os.Open(image)
		err = client.CopyFromFile(context.Background(), *f, fmt.Sprintf("/home/%s/%s", node.Username, image), PERMISSIONS)

		if err != nil {
			return err
		}

	}
	return nil
}
