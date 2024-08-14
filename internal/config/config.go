package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Screen    Screen    `yaml:"screen"`
	Processes []Process `yaml:"processes"`
}

type Network struct {
	checks []Check `yaml:"checks"`
}

type Check struct {
	Name     string `yaml:"name"`
	Url      string `yaml:"url"`
	WaitTime int    `yaml:"wait_time"`
}

type Process struct {
	Name        string        `yaml:"name"`
	Command     string        `yaml:"command"`
	Wait        bool          `yaml:"wait"`
	WaitTime    time.Duration `yaml:"wait_time"`
	Watch       bool          `yaml:"watch"`
	Arguments   []string      `yaml:"arguments"`
	LogFile     string        `yaml:"log_file"`
	Environment []Environment `yaml:"environment"`
}

type Environment struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Screen struct {
	FullScreen bool   `yaml:"full_screen"`
	BootImage  string `yaml:"boot_image"`
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
