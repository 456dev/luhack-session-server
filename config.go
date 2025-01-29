package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server struct {
		Host     string `yaml:"host"`
		Domain   string `yaml:"domain"`
		Protocol string `yaml:"protocol"`
	} `yaml:"server"`
	Session struct {
		Title      string `yaml:"title"`
		BackendMap string `yaml:"backendMap"`
	} `yaml:"session"`
	Security struct {
		JwtSecret string `yaml:"jwtSecret"`
		Server    string `yaml:"server"`
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
