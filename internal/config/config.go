package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type CoreConfig struct {
	PRD struct {
		PRDFile            string `yaml:"prdFile"`
		PRDSharded         bool   `yaml:"prdSharded"`
		PRDShardedLocation string `yaml:"prdShardedLocation"`
	} `yaml:"prd"`
	Architecture struct {
		ArchitectureFile            string `yaml:"architectureFile"`
		ArchitectureSharded         bool   `yaml:"architectureSharded"`
		ArchitectureShardedLocation string `yaml:"architectureShardedLocation"`
	} `yaml:"architecture"`
	DevStoryLocation string `yaml:"devStoryLocation"`
}

func LoadCoreConfig(rootPath string) (*CoreConfig, error) {
	configPath := filepath.Join(rootPath, ".bmad-core", "core-config.yaml")
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return getDefaultConfig(), nil
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return getDefaultConfig(), nil
	}
	
	var config CoreConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return getDefaultConfig(), nil
	}
	
	return &config, nil
}

func getDefaultConfig() *CoreConfig {
	return &CoreConfig{
		PRD: struct {
			PRDFile            string `yaml:"prdFile"`
			PRDSharded         bool   `yaml:"prdSharded"`
			PRDShardedLocation string `yaml:"prdShardedLocation"`
		}{
			PRDFile:            "docs/prd.md",
			PRDSharded:         false,
			PRDShardedLocation: "docs/prd",
		},
		Architecture: struct {
			ArchitectureFile            string `yaml:"architectureFile"`
			ArchitectureSharded         bool   `yaml:"architectureSharded"`
			ArchitectureShardedLocation string `yaml:"architectureShardedLocation"`
		}{
			ArchitectureFile:            "docs/architecture.md",
			ArchitectureSharded:         false,
			ArchitectureShardedLocation: "docs/architecture",
		},
		DevStoryLocation: "docs/stories",
	}
}