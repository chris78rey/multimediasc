package app

import (
	"errors"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"multimediasc/internal/oracle"
)

func enrichDocuments(docs []oracle.ImagenPacienteRow) {
	for i := range docs {
		obs := inspectFile(docs[i].Ruta)
		docs[i].Ext = obs.Ext
		docs[i].Kind = obs.Kind
		docs[i].State = obs.State
		docs[i].Reason = obs.Reason
		docs[i].Hidden = obs.Hidden
		docs[i].Accessible = obs.Accessible
		docs[i].Exists = obs.Exists
	}
}

type fileObservation struct {
	Ext        string
	Kind       string
	State      string
	Reason     string
	Hidden     bool
	Accessible bool
	Exists     bool
}

func inspectFile(path string) fileObservation {
	obs := fileObservation{
		Ext:  strings.ToLower(filepath.Ext(path)),
		Kind: fileKindFromExt(filepath.Ext(path)),
	}

	st, err := os.Stat(path)
	if err != nil {
		obs.State = "inaccesible"
		obs.Reason = classifyFileError(err)
		return obs
	}
	obs.Exists = true
	obs.Hidden = isHidden(path, st)

	f, err := os.Open(path)
	if err != nil {
		obs.State = "inaccesible"
		obs.Reason = classifyFileError(err)
		if obs.Hidden {
			obs.Reason = appendReason(obs.Reason, "oculto")
		}
		return obs
	}
	_ = f.Close()

	obs.Accessible = true
	obs.State = "ok"
	if obs.Hidden {
		obs.State = "oculto"
		obs.Reason = "atributo Hidden en el filesystem"
	}
	return obs
}

func classifyFileError(err error) string {
	if err == nil {
		return ""
	}
	var pe *os.PathError
	if errors.As(err, &pe) {
		switch {
		case errors.Is(pe.Err, os.ErrNotExist):
			return "no existe"
		case errors.Is(pe.Err, os.ErrPermission):
			return "sin permisos"
		}
	}
	if errors.Is(err, os.ErrPermission) {
		return "sin permisos"
	}
	if errors.Is(err, os.ErrNotExist) {
		return "no existe"
	}
	return err.Error()
}

func fileKindFromExt(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".pdf":
		return "pdf"
	case ".doc", ".docx":
		return "word"
	case ".xls", ".xlsx":
		return "excel"
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".tif", ".tiff", ".webp":
		return "imagen"
	default:
		if t := mime.TypeByExtension(ext); t != "" {
			return strings.ToLower(strings.SplitN(t, "/", 2)[0])
		}
		return "archivo"
	}
}

func appendReason(base, extra string) string {
	if base == "" {
		return extra
	}
	if extra == "" {
		return base
	}
	return base + "; " + extra
}
