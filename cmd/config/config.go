package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

func GetLocal(h *HTTPConfig) *HTTPConfig {
	yamlFile, err := os.ReadFile("local.yaml")
	if err != nil {
		log.Printf("local yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, h)
	if err != nil {
		log.Fatalf("Unmarshal local: %v", err)
	}
	return h
}

func GetConfig(c *Config) *Config {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("config yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal config: %v", err)
	}
	return c
}

type Config struct {
	Storage StorageConfig `yaml:"storage"`
}

type HTTPConfig struct {
	Env        string `yaml:"env" env-default:"development"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type StorageConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}
