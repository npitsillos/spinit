package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/npitsillos/spinit/errors"
	"github.com/npitsillos/spinit/helpers"
)

var (
	configDir  = ".spinit"
	configFile = "spinit.toml"
)

func InitSpinitConfig() (*Config, error) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Looks like this is your first time running spinit!")
	choices := "Y/n"
	fmt.Printf("Follow the instructions here or simply run `spinit init`. Continue? %s: ", choices)

	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())
	if strings.ToLower(input) == "n" || strings.ToLower(input) == "no" {
		os.Exit(0)
	}

	dir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	cfg := newConfig()

	fmt.Print("How many nodes does your cluster have? ")
	scanner.Scan()
	numNodesStr := scanner.Text()
	numNodes, err := strconv.Atoi(numNodesStr)
	if err != nil {
		return nil, err
	}

	for i := 0; i < numNodes; i++ {
		fmt.Print("Provide the node's namd & IP address and the user's username (<node-name>,<node-ip>,<username>): ")
		scanner.Scan()
		nodeInfo := scanner.Text()
		nodeInfoArr := strings.Split(nodeInfo, ",")
		cfg.Nodes = append(cfg.Nodes, &Node{
			Name:     nodeInfoArr[0],
			IP:       nodeInfoArr[1],
			Username: nodeInfoArr[2],
		})
	}

	fmt.Print("Provide the path to the SSH key used to access the nodes: ")
	scanner.Scan()
	sshKeyPath := scanner.Text()

	cfg.SSH = sshKeyPath

	if err := initConfigDir(dir); err != nil {
		return nil, err
	}

	cfg.ConfigFilePath = filepath.Join(dir, configFile)

	return cfg, cfg.WriteToFile()
}

func ConfigDirExists() (bool, error) {
	spinitRootPath, err := getConfigDir()
	if err != nil {
		return false, err
	}

	return !helpers.DirectoryExists(spinitRootPath), nil
}

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.ErrAccessingHomeDir
	}
	return filepath.Join(homeDir, configDir), nil
}

func initConfigDir(dir string) error {
	if !helpers.DirectoryExists(dir) {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return err
		}
	}

	return nil
}
