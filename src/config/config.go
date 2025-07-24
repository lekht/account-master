package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type SuperuserConf struct {
	Email    string `yaml:"email"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Admin    bool   `yaml:"admin"`
}

type Config struct {
	Server ServerConf    `yaml:"server"`
	Admin  SuperuserConf `yaml:"superuser"`
}

// Load app config. Requires path to yaml config file
func Load(path string, conf *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read from config file: %v", err)
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("failed to close config file: %v", err)
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config data to struct: %v", err)
	}

	return nil
}
