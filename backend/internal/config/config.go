package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Storage   StorageConfig   `yaml:"storage"`
	Ingest    IngestConfig    `yaml:"ingest"`
	Auth      AuthConfig      `yaml:"auth"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type StorageConfig struct {
	Path           string `yaml:"path"`
	ChunkSizeBytes int    `yaml:"chunk_size_bytes"`
	RetentionDays  int    `yaml:"retention_days"`
}

type IngestConfig struct {
	BufferSize    int `yaml:"buffer_size"`
	FlushInterval int `yaml:"flush_interval_ms"`
}

type AuthConfig struct {
	Enabled bool   `yaml:"enabled"`
	APIKey  string `yaml:"api_key"`
}

type RateLimitConfig struct {
	Enabled      bool     `yaml:"enabled"`
	RequestsPer  int      `yaml:"requests_per_minute"`
	Burst        int      `yaml:"burst"`
	IngestOnly   bool     `yaml:"ingest_only"`
	WhitelistIPs []string `yaml:"whitelist_ips"`
	BlacklistIPs []string `yaml:"blacklist_ips"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// Return default config if file not found
		return DefaultConfig(), nil
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Override with environment variables
	if port := os.Getenv("LOKILITE_PORT"); port != "" {
		cfg.Server.Port = port
	}
	if apiKey := os.Getenv("LOKILITE_API_KEY"); apiKey != "" {
		cfg.Auth.APIKey = apiKey
		cfg.Auth.Enabled = true
	}
	if storagePath := os.Getenv("LOKILITE_STORAGE_PATH"); storagePath != "" {
		cfg.Storage.Path = storagePath
	}
	if rateLimitEnabled := os.Getenv("LOGPULSE_RATE_LIMIT_ENABLED"); rateLimitEnabled != "" {
		cfg.RateLimit.Enabled = rateLimitEnabled == "true"
	}
	if rateLimitRequestsPerMinute := os.Getenv("LOGPULSE_RATE_LIMIT_REQUESTS_PER_MINUTE"); rateLimitRequestsPerMinute != "" {
		if val, err := strconv.Atoi(rateLimitRequestsPerMinute); err == nil {
			cfg.RateLimit.RequestsPer = val
		}
	}
	if rateLimitBurst := os.Getenv("LOGPULSE_RATE_LIMIT_BURST"); rateLimitBurst != "" {
		if val, err := strconv.Atoi(rateLimitBurst); err == nil {
			cfg.RateLimit.Burst = val
		}
	}

	return &cfg, nil
}

func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
		},
		Storage: StorageConfig{
			Path:           "./data/logs",
			ChunkSizeBytes: 1024 * 1024, // 1MB
			RetentionDays:  7,
		},
		Ingest: IngestConfig{
			BufferSize:    1000,
			FlushInterval: 5000,
		},
		Auth: AuthConfig{
			Enabled: false,
			APIKey:  "",
		},
		RateLimit: RateLimitConfig{
			Enabled:      true,
			RequestsPer:  1000,
			Burst:        100,
			IngestOnly:   false, // Apply to all endpoints by default
			WhitelistIPs: []string{},
			BlacklistIPs: []string{},
		},
	}
}
