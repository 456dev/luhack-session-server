package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server struct {
		Host       string `yaml:"host"`
		Domain     string `yaml:"domain"`
		Protocol   string `yaml:"protocol"`
		BackendMap string `yaml:"backendMap"`
	} `yaml:"server"`
	Security struct {
		JwtSecret string `yaml:"jwtSecret"`
	} `yaml:"security"`
}

func parseConfig(file string, config **Config) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	return nil
}
