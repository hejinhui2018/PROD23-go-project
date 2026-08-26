package config

import (
	"os"
	"strconv"
	"strings"
)

func ApplyEnv(c *Config) {
	if v := os.Getenv("RECOVERY_HTTP"); v != "" {
		c.HTTPAddr = v
	}
	if v := os.Getenv("RECOVERY_LOG_LEVEL"); v != "" {
		c.LogLevel = strings.ToLower(v)
	}
	if v := os.Getenv("RECOVERY_SNAPSHOT_EVERY"); v != "" {
		if n, e := strconv.Atoi(v); e == nil && n > 0 {
			c.SnapshotEvery = n
		}
	}
}
func (c Config) Valid() bool { return c.HTTPAddr != "" && c.DataDir != "" && c.SnapshotEvery > 0 }
func (c Config) WithDataDir(dir string) Config {
	if dir != "" {
		c.DataDir = dir
	}
	return c
}
func (c Config) Address() string       { return c.HTTPAddr }
func (c Config) SnapshotEnabled() bool { return c.SnapshotEvery > 0 }
