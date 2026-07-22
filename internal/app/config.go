package app

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	ListenAddr        string   `json:"listen_addr"`
	OracleConnect     string   `json:"oracle_connect"`
	DefaultMaxOpen    int      `json:"default_max_open"`
	DefaultMaxIdle    int      `json:"default_max_idle"`
	AllowedNames      []string `json:"allowed_names"`
	AllowDuplicateFix bool     `json:"allow_duplicate_fix"`
}

func LoadEnvConfig() Config {
	cfg := Config{
		ListenAddr:        getenv("MULTIMEDIASC_LISTEN", "127.0.0.1:8080"),
		OracleConnect:     getenv("ORACLE_CONNECT", "172.16.60.21:1521/prdsgh2"),
		DefaultMaxOpen:    3,
		DefaultMaxIdle:    1,
		AllowedNames:      []string{"013B.pdf", "Epicrisis.pdf", "Consentimiento_Informado.pdf", "Protocolo_Quirurgico.pdf", "Resultados_Laboratorio.pdf"},
		AllowDuplicateFix: false,
	}

	_ = loadKeyValues("oracle.local.env")
	_ = loadKeyValues("oracle.local.env.example")

	if b, err := os.ReadFile("config.json"); err == nil {
		_ = json.Unmarshal(b, &cfg)
		if len(cfg.AllowedNames) == 0 {
			cfg.AllowedNames = []string{"013B.pdf", "Epicrisis.pdf", "Consentimiento_Informado.pdf", "Protocolo_Quirurgico.pdf", "Resultados_Laboratorio.pdf"}
		}
	}

	if v := os.Getenv("MULTIMEDIASC_ALLOW_DUPLICATE_FIX"); v == "1" || v == "true" {
		cfg.AllowDuplicateFix = true
	}
	return cfg
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func (c Config) ConfigPath() string {
	return filepath.Join(".", "config.json")
}

func loadKeyValues(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		if key != "" && os.Getenv(key) == "" {
			_ = os.Setenv(key, val)
		}
	}
	return s.Err()
}
