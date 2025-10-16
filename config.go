package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AllocationMode string

const (
	AllocationModeCentral AllocationMode = "central"
	AllocationModePerUser AllocationMode = "per-user"
)

func (m *AllocationMode) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	switch s {
	case string(AllocationModeCentral), string(AllocationModePerUser):
		*m = AllocationMode(s)
		return nil
	default:
		return fmt.Errorf("invalid allocationMode: %q", s)
	}
}

type Config struct {
	Server struct {
		Host     string `yaml:"host"`
		Domain   string `yaml:"domain"`
		Protocol string `yaml:"protocol"`
	} `yaml:"server"`
	Session struct {
		Title          string         `yaml:"title"`
		BackendMap     string         `yaml:"backendMap"`
		AllocationMode AllocationMode `yaml:"allocationMode"`
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
