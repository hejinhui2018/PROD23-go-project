package config

import (
	"encoding/json"
	"flag"
	"os"
)

type Config struct {
	HTTPAddr      string `json:"http_addr"`
	DataDir       string `json:"data_dir"`
	SnapshotEvery int    `json:"snapshot_every"`
	LogLevel      string `json:"log_level"`
}

func Default() Config {
	return Config{HTTPAddr: ":8080", DataDir: "./data", SnapshotEvery: 50, LogLevel: "info"}
}
func (c Config) Clone() Config       { return c }
func (c Config) IsJSONLogging() bool { return c.LogLevel == "json" || c.LogLevel == "structured" }
func (c Config) String() string      { return c.HTTPAddr + " " + c.DataDir }

func Load() Config {
	c := Config{HTTPAddr: ":8080", DataDir: "./data", SnapshotEvery: 50, LogLevel: "info"}
	p := flag.String("config", "", "config file")
	flag.Parse()
	if *p != "" {
		b, _ := os.ReadFile(*p)
		_ = json.Unmarshal(b, &c)
	}
	if v := os.Getenv("RECOVERY_DATA"); v != "" {
		c.DataDir = v
	}
	return c
}
