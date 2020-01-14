package config

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const DEFAULT_CONFIG_PATH = "$HOME/.jarvis/config.yaml"

type (
	ConfigurationReader struct {
		path string
	}
)

func InitConfigurationReader(path string) *ConfigurationReader {
	return &ConfigurationReader{path}
}

func (c *ConfigurationReader) ParseConfig() error {
	dir, file := filepath.Split(c.path)
	//viper don't care about the file type
	//we must remove it from the file name
	file = fileNameWithoutExtension(file)
	viper.SetConfigType("yaml")
	viper.SetConfigName(file)
	viper.AddConfigPath(dir)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigurationReader) GetConfigPath() string {
	return c.path
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
