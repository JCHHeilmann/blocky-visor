package blocky

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// ReadConfig reads the raw Blocky config file.
func ReadConfig(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return data, nil
}

// ValidateYAML checks that data is valid YAML.
func ValidateYAML(data []byte) error {
	var out interface{}
	if err := yaml.Unmarshal(data, &out); err != nil {
		return fmt.Errorf("invalid YAML: %w", err)
	}
	return nil
}

// BackupConfig creates a timestamped backup of the config file.
func BackupConfig(path string) (string, error) {
	backupPath := fmt.Sprintf("%s.bak.%s", path, time.Now().Format("20060102-150405"))
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read for backup: %w", err)
	}
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return "", fmt.Errorf("write backup: %w", err)
	}
	return backupPath, nil
}

// WriteConfig validates YAML, backs up existing config, and writes new config.
func WriteConfig(path string, data []byte) (backupPath string, err error) {
	if err := ValidateYAML(data); err != nil {
		return "", err
	}

	backupPath, err = BackupConfig(path)
	if err != nil {
		return "", fmt.Errorf("backup failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return backupPath, fmt.Errorf("write config: %w", err)
	}

	return backupPath, nil
}
