package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	ListenAddr        string   `json:"listen_addr"`
	OracleConnect     string   `json:"oracle_connect"`
	DefaultMaxOpen    int      `json:"default_max_open"`
	DefaultMaxIdle    int      `json:"default_max_idle"`
	MaxBatchPlanillas int      `json:"max_batch_planillas"`
	AllowedNames      []string `json:"allowed_names"`
	PlanillaRanges    []string `json:"planilla_ranges"`
	AllowDuplicateFix bool     `json:"allow_duplicate_fix"`
}

func LoadEnvConfig() Config {
	cfg := Config{
		ListenAddr:        getenv("MULTIMEDIASC_LISTEN", "127.0.0.1:8080"),
		OracleConnect:     getenv("ORACLE_CONNECT", "172.16.60.21:1521/prdsgh2"),
		DefaultMaxOpen:    3,
		DefaultMaxIdle:    1,
		MaxBatchPlanillas: 25,
		AllowedNames:      []string{"013B", "Epicrisis", "Consentimiento_Informado", "Protocolo_Quirurgico", "Resultados_Laboratorio", "Imagen_Estudio"},
		PlanillaRanges:    nil,
		AllowDuplicateFix: false,
	}

	_ = loadKeyValues("oracle.local.env")
	_ = loadKeyValues("oracle.local.env.example")

	if b, err := os.ReadFile("config.json"); err == nil {
		_ = json.Unmarshal(b, &cfg)
		if len(cfg.AllowedNames) == 0 {
			cfg.AllowedNames = []string{"013B", "Epicrisis", "Consentimiento_Informado", "Protocolo_Quirurgico", "Resultados_Laboratorio", "Imagen_Estudio"}
		}
		if cfg.MaxBatchPlanillas <= 0 {
			cfg.MaxBatchPlanillas = 25
		}
		cfg.AllowedNames = normalizeAllowedNames(cfg.AllowedNames)
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

func (c Config) Save() error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(c.ConfigPath(), b, 0o644); err != nil {
		return fmt.Errorf("no se pudo guardar config.json: %w", err)
	}
	return nil
}

func (c Config) AllowedRangesText() string {
	return strings.Join(c.PlanillaRanges, "\n")
}

func (c Config) AllowedNamesText() string {
	return strings.Join(c.AllowedNames, "\n")
}

func parseRangesText(text string) []string {
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func parseNamesText(text string) []string {
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))
	seen := map[string]bool{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.TrimSuffix(line, filepath.Ext(line))
		if line == "" || seen[strings.ToLower(line)] {
			continue
		}
		seen[strings.ToLower(line)] = true
		out = append(out, line)
	}
	return out
}

func normalizeAllowedNames(names []string) []string {
	out := make([]string, 0, len(names))
	seen := map[string]bool{}
	for _, n := range names {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		n = strings.TrimSuffix(n, filepath.Ext(n))
		if n == "" || seen[strings.ToLower(n)] {
			continue
		}
		seen[strings.ToLower(n)] = true
		out = append(out, n)
	}
	if len(out) == 0 {
		return []string{"013B", "Epicrisis", "Consentimiento_Informado", "Protocolo_Quirurgico", "Resultados_Laboratorio", "Imagen_Estudio"}
	}
	return out
}

func (c Config) MaxBatchPlanillasOrDefault() int {
	if c.MaxBatchPlanillas <= 0 {
		return 25
	}
	return c.MaxBatchPlanillas
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
