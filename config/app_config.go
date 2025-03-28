package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

const (
	DefaultAppConfigFile = "spinit.toml"
)

type AppConfig struct {
	AppName string `toml:"app"`
	Build   *Build `toml:"build"`

	configFilePath string
}

type Build struct {
	Image      string `toml:"image"`
	Dockerfile string `toml:"dockerfile"`
	Version    string `toml:"version"`
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}

func LoadAppConfig(path string) (*AppConfig, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	appCfg, err := unMarshalAppConfig(buf)
	if err != nil {
		return nil, err
	}

	appCfg.configFilePath = path
	return appCfg, nil
}

func (appCfg *AppConfig) SetConfigFilePath(path string) {
	appCfg.configFilePath = path
}

func (appCfg *AppConfig) WriteToFile() error {
	var file *os.File
	var err error

	if file, err = os.Create(appCfg.configFilePath); err != nil {
		return err
	}

	defer func() {
		err = file.Close()
	}()

	cfgBytes, err := appCfg.Marshal()
	_, err = bytes.NewBuffer(cfgBytes).WriteTo(file)
	return err
}

func (appCfg *AppConfig) Marshal() ([]byte, error) {
	var b bytes.Buffer
	encoder := toml.NewEncoder(&b)
	encoder.SetIndentTables(true)
	if err := encoder.Encode(appCfg); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func unMarshalAppConfig(buf []byte) (*AppConfig, error) {
	cfg := &AppConfig{}
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
