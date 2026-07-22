package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"multimediasc/internal/oracle"
)

func main() {
	_ = loadEnvFile("oracle.local.env")
	_ = loadEnvFile("oracle.local.env.example")

	user := firstNonEmpty(os.Getenv("ORACLE_USER"), os.Getenv("USER"))
	pass := os.Getenv("ORACLE_PASSWORD")
	connect := firstNonEmpty(os.Getenv("ORACLE_CONNECT"), "172.16.60.21:1521/prdsgh2")
	tramite := firstNonEmpty(os.Getenv("ORACLE_TRAMITE"), "6076406")

	if user == "" || pass == "" {
		fmt.Println("missing ORACLE_USER or ORACLE_PASSWORD")
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	repo, err := oracle.Open(ctx, user, pass, oracle.OpenConfig{
		ConnectString: connect,
		MaxOpenConns:  2,
		MaxIdleConns:  1,
		ConnMaxLife:   5 * time.Minute,
	})
	if err != nil {
		fmt.Printf("open error: %v\n", err)
		os.Exit(1)
	}
	defer repo.Close()

	var planilla int64
	_, _ = fmt.Sscan(tramite, &planilla)
	det, err := repo.ObtenerDetallePlanilla(ctx, planilla)
	if err != nil {
		fmt.Printf("query error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("ok planilla=%d hc=%d docs=%d\n", det.Planilla.DigTramite, det.Planilla.DigHC, len(det.Documentos))
	if det.Paciente != nil {
		fmt.Printf("paciente=%s cedula=%s\n", det.Paciente.NombreCompleto(), det.Paciente.Cedula)
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func loadEnvFile(path string) error {
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
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, `"'`)
		if key != "" {
			_ = os.Setenv(key, val)
		}
	}
	return s.Err()
}
