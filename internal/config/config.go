package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	ZomboidFolder       string `json:"zomboid_folder"`
	BackupFolder        string `json:"backup_folder"`
	BackupInterval      int    `json:"backup_interval"`
	MaxBackupFiles      int    `json:"max_backup_files"`
	AutoBackup          bool   `json:"auto_backup"`
	BackupBeforeRestore bool   `json:"backup_before_restore"`
}

// StatePath returns the fixed path of the state file.
// It always lives in ~/ZomboidAutoBackup which is guaranteed to exist.
func StatePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "ZomboidAutoBackup", ".state")
}

// Load reads saved settings from disk, falling back to defaults if the file
// does not exist or cannot be parsed.
func Load() (*Config, error) {
	cfg := Default()
	data, err := os.ReadFile(StatePath())
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, fmt.Errorf("reading state file: %w", err)
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return cfg, fmt.Errorf("parsing state file: %w", err)
	}
	return cfg, nil
}

// Save writes the current settings to the state file.
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding settings: %w", err)
	}
	if err := os.WriteFile(StatePath(), data, 0o644); err != nil {
		return fmt.Errorf("writing state file: %w", err)
	}
	return nil
}

// EnsureBackupFolder creates the backup folder if it does not exist.
func (c *Config) EnsureBackupFolder() error {
	if err := os.MkdirAll(c.BackupFolder, 0o755); err != nil {
		return fmt.Errorf("could not create backup folder %q: %w", c.BackupFolder, err)
	}
	return nil
}

func Default() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		ZomboidFolder:  filepath.Join(home, "Zomboid"),
		BackupFolder:   filepath.Join(home, "ZomboidAutoBackup"),
		BackupInterval: 15,
		MaxBackupFiles: 10,
		AutoBackup:          false,
		BackupBeforeRestore: true,
	}
}
