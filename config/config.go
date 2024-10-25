package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Nodes          []*Node `toml:"nodes"`
	SSH            string  `toml:"ssh_key_path"`
	ConfigFilePath string  `toml:"-"`
}

type Node struct {
	Name     string `toml:"name"`
	IP       string `toml:"ip"`
	Username string `toml:"username"`
}

func newConfig() *Config {
	return &Config{}
}

func LoadConfig() (*Config, error) {

	dir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	buf, err := os.ReadFile(filepath.Join(dir, configFile))
	if err != nil {
		return nil, err
	}

	cfg, err := unMarshal(buf)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) WriteToFile() error {
	var file *os.File
	var err error
	if file, err = os.Create(cfg.ConfigFilePath); err != nil {
		return err
	}

	defer func() {
		err = file.Close()
	}()

	cfgBytes, err := cfg.Marshal()
	_, err = bytes.NewBuffer(cfgBytes).WriteTo(file)

	return err
}

func (cfg *Config) Marshal() ([]byte, error) {
	var b bytes.Buffer
	encoder := toml.NewEncoder(&b)
	encoder.SetIndentTables(true)
	if err := encoder.Encode(cfg); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func unMarshal(buf []byte) (*Config, error) {
	cfg := newConfig()
	if err := toml.Unmarshal(buf, &cfg); err != nil {
		var derr *toml.DecodeError
		if errors.As(err, &derr) {
			row, col := derr.Position()
			return nil, fmt.Errorf("row %d column %d\n%s", row, col, derr.String())
		}
		return nil, err
	}
	return cfg, nil
}
