package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Layout struct {
	ID       string `yaml:"id"`
	Name     string `yaml:"name"`
	Services []struct {
		Description string `yaml:"description"`
		ID          string `yaml:"id"`
		Name        string `yaml:"name"`
		Proxy       string `yaml:"proxy"`
	} `yaml:"services"`
}

type Backend struct {
	ID       string `yaml:"id"`
	Services []struct {
		BoxID     string `yaml:"box_id"`
		Host      string `yaml:"host"`
		ServiceID string `yaml:"service_id"`
	} `yaml:"services"`
}

type BackendMap struct {
	Backends   []Backend `yaml:"backends"`
	Layout     []Layout  `yaml:"layout"`
	LbEndpoint string    `yaml:"lb_endpoint"`
}

func parseBackendMap(file string, backendMap **BackendMap) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &backendMap)
	if err != nil {
		return err
	}
	return nil
}
