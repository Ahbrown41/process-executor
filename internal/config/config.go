package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Display   Display   `yaml:"display"`
	Processes []Process `yaml:"processes"`
}

type Process struct {
	Name          string        `yaml:"name"`
	PreConditions []Condition   `yaml:"preConditions"`
	WorkDir       string        `yaml:"workDir"`
	Command       string        `yaml:"command"`
	Wait          bool          `yaml:"wait"`
	WaitMax       time.Duration `yaml:"waitMax"`
	Restart       bool          `yaml:"restart"`
	Arguments     []string      `yaml:"arguments"`
	Environment   []Attribute   `yaml:"environment"`
}

type Condition struct {
	Name       string        `yaml:"name"`
	Type       string        `yaml:"type"`
	Wait       bool          `yaml:"wait"`
	Timeout    time.Duration `yaml:"timeout"`
	Attributes []Attribute   `yaml:"attributes"`
}

type Attribute struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Display struct {
	FullScreen bool   `yaml:"fullScreen"`
	BootImage  string `yaml:"bootImage"`
}

// Load reads a configuration file and returns a Config.
func Load(file string) (*Config, error) {
	var config Config
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
