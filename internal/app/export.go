package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"multimediasc/internal/oracle"
)

type zipDocument struct {
	Index     int    `json:"index"`
	Selected  bool   `json:"selected"`
	ExportAs  string `json:"export_as"`
	ValidName bool   `json:"valid_name"`
	Planilla  string `json:"planilla"`
	Ruta      string `json:"ruta"`
}

func friendlyErr(err error) string {
	if err == nil {
		return ""
	}
	switch {
	case oracle.IsKind(err, oracle.ErrInvalidCredentials):
		return "usuario o contraseña incorrectos"
	case oracle.IsKind(err, oracle.ErrAccountLocked):
		return "la cuenta de Oracle está bloqueada"
	case oracle.IsKind(err, oracle.ErrPasswordExpired):
		return "la contraseña de Oracle está vencida"
	case oracle.IsKind(err, oracle.ErrNoQueryPermission):
		return "el usuario no tiene permisos de consulta"
	case oracle.IsKind(err, oracle.ErrDBUnavailable):
		return "la base de datos no está disponible"
	case errorsIs(err, context.DeadlineExceeded):
		return "la operación excedió el tiempo de espera"
	case errorsIs(err, os.ErrNotExist):
		return "el archivo no existe"
	default:
		return err.Error()
	}
}

func errorsIs(err error, target error) bool {
	if err == nil || target == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), strings.ToLower(target.Error()))
}

func exportZip(planilla string, docs []zipDocument, lookup map[string]string, dest string, allowDup bool) error {
	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("el archivo destino ya existe")
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	return createZip(planilla, docs, lookup, f, allowDup)
}

func createZip(planilla string, docs []zipDocument, lookup map[string]string, w io.Writer, allowDup bool) error {
	type pair struct {
		planilla string
		name     string
		path     string
	}
	var items []pair
	seen := map[string]int{}
	for _, d := range docs {
		if !d.Selected {
			continue
		}
		path := strings.TrimSpace(d.Ruta)
		if path == "" {
			path = lookup[strconv.Itoa(d.Index)]
		}
		if path == "" {
			continue
		}
		name := normalizeExportName(d.ExportAs, path)
		if !validWindowsFilename(name) {
			return fmt.Errorf("nombre inválido: %s", name)
		}
		if count := seen[strings.ToLower(name)]; count > 0 {
			if !allowDup {
				return fmt.Errorf("nombre duplicado: %s", name)
			}
			base := strings.TrimSuffix(name, filepath.Ext(name))
			ext := filepath.Ext(name)
			name = fmt.Sprintf("%s_%d%s", base, count, ext)
		}
		seen[strings.ToLower(name)]++
		itemPlanilla := strings.TrimSpace(d.Planilla)
		if itemPlanilla == "" {
			itemPlanilla = planilla
		}
		if itemPlanilla == "" {
			return fmt.Errorf("falta planilla para %s", name)
		}
		items = append(items, pair{planilla: itemPlanilla, name: name, path: path})
	}
	if len(items) == 0 {
		return fmt.Errorf("no hay documentos seleccionados")
	}
	buf := &bytes.Buffer{}
	zw := newZipWriter(buf)
	for _, it := range items {
		if err := zw.addFile(filepath.Join(it.planilla, it.name), it.path); err != nil {
			return err
		}
	}
	if err := zw.close(); err != nil {
		return err
	}
	_, err := io.Copy(w, buf)
	return err
}

func countSelected(docs []zipDocument) int {
	total := 0
	for _, d := range docs {
		if d.Selected {
			total++
		}
	}
	return total
}

func defaultZipFilename(planilla string, count int) string {
	if count > 1 {
		return "pl_" + time.Now().Format("200601021504") + ".zip"
	}
	planilla = strings.TrimSpace(planilla)
	if planilla == "" {
		return "pl_" + time.Now().Format("200601021504") + ".zip"
	}
	return planilla + ".zip"
}

func normalizeExportName(name, sourcePath string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = filepath.Base(sourcePath)
	}
	if filepath.Ext(name) == "" {
		if ext := filepath.Ext(sourcePath); ext != "" {
			name += ext
		}
	}
	return name
}

func validWindowsFilename(name string) bool {
	if name == "" {
		return false
	}
	bad := `<>:"/\|?*`
	for _, r := range name {
		if strings.ContainsRune(bad, r) || r < 32 {
			return false
		}
	}
	return true
}
