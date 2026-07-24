package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	coreapp "multimediasc/internal/app"
	stdruntime "runtime"
)

type App struct {
	ctx        context.Context
	core       *coreapp.BatchCore
	status     string
	preview    string
	docInfo    string
	exportInfo string
	batchInfo  string
	rangeInfo  string
	zipDest    string
	outDir     string
	searchText string
}

func NewApp() *App {
	cfg := coreapp.LoadEnvConfig()
	return &App{
		core:       coreapp.NewBatchCore(cfg),
		status:     "Listo",
		preview:    "Busca una o más planillas para revisar documentos.",
		docInfo:    "Documento activo: ninguno",
		exportInfo: "Prepara la exportación para ver el resumen aquí.",
		batchInfo:  "Planillas cargadas: 0",
		rangeInfo:  "Rangos activos: " + strings.Join(cfg.PlanillaRanges, ", "),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	wailsruntime.WindowMaximise(a.ctx)
}

func (a *App) State() coreapp.BatchSnapshot {
	return a.core.Snapshot(a.status, a.preview, a.docInfo, a.exportInfo, a.batchInfo, a.rangeInfo, a.zipDest, a.outDir)
}

func (a *App) SearchText() string {
	return a.searchText
}

func (a *App) Login(user, pass string) coreapp.BatchSnapshot {
	if err := a.core.Login(context.Background(), user, pass); err != nil {
		a.status = friendlyErr(err)
		return a.State()
	}
	a.status = "sesión iniciada"
	return a.State()
}

func (a *App) Logout() coreapp.BatchSnapshot {
	a.core.Logout()
	a.status = "sesión cerrada"
	a.preview = "sesión cerrada"
	a.docInfo = "Documento activo: ninguno"
	a.exportInfo = "Prepara la exportación para ver el resumen aquí."
	a.batchInfo = "Planillas cargadas: 0"
	return a.State()
}

func (a *App) SetSearchText(text string) coreapp.BatchSnapshot {
	a.searchText = strings.TrimSpace(text)
	return a.State()
}

func (a *App) Search() coreapp.BatchSnapshot {
	inputs := parsePlanillaInputs(a.searchText)
	items, err := a.core.Search(context.Background(), inputs)
	if err != nil {
		a.status = friendlyErr(err)
		return a.State()
	}
	a.status = fmt.Sprintf("cargadas %d planillas", len(items))
	a.batchInfo = fmt.Sprintf("Planillas cargadas: %d | Activa: %d", len(items), 1)
	if len(items) == 0 {
		a.preview = "Busca una o más planillas para revisar documentos."
	}
	return a.State()
}

func (a *App) SelectOutputDirectory() coreapp.BatchSnapshot {
	dir, err := wailsruntime.OpenDirectoryDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title: "Selecciona dónde guardar los archivos",
	})
	if err != nil || dir == "" {
		return a.State()
	}
	a.outDir = strings.TrimSpace(dir)
	a.zipDest = ""
	a.exportInfo = "Carpeta seleccionada: " + a.outDir
	return a.State()
}

func (a *App) UpdateZipDestination() coreapp.BatchSnapshot {
	if strings.TrimSpace(a.outDir) == "" {
		a.outDir = coreapp.DefaultAutoOutDir()
	}
	dest := coreapp.DefaultZipPath(a.outDir, a.core.ExportPlanillaLabel(), a.core.SelectedDocCount())
	a.zipDest = dest
	return a.State()
}

func (a *App) SelectPlanilla(index int) coreapp.BatchSnapshot {
	a.core.SetActiveIndex(index)
	a.refreshPlanillaState()
	return a.State()
}

func (a *App) MovePlanilla(delta int) coreapp.BatchSnapshot {
	idx := a.core.MovePlanilla(delta, a.core.FilteredPlanillaIndexes(""))
	a.core.SetActiveIndex(idx)
	a.refreshPlanillaState()
	return a.State()
}

func (a *App) MoveDocumento(delta int) coreapp.BatchSnapshot {
	idx := a.core.MoveDocumento(delta)
	a.docInfo = fmt.Sprintf("Documento activo: %d", idx+1)
	a.preview = a.core.DocSummary(a.core.ActiveIndex(), idx)
	return a.State()
}

func (a *App) ToggleDocumento(planillaIndex, docIndex int, selected bool) coreapp.BatchSnapshot {
	a.core.SetSelected(planillaIndex, docIndex, selected)
	if selected {
		a.core.SetActiveIndex(planillaIndex)
		a.core.AdvanceToNextPending()
	}
	a.refreshExportPreview()
	a.refreshPlanillaState()
	return a.State()
}

func (a *App) ClearSelection() coreapp.BatchSnapshot {
	a.core.ClearSelection()
	a.refreshExportPreview()
	return a.State()
}

func (a *App) RemoveSelection(planillaIndex, docIndex int) coreapp.BatchSnapshot {
	a.core.RemoveSelection(planillaIndex, docIndex)
	a.core.SetActiveIndex(planillaIndex)
	if docIndex >= 0 {
		a.core.SetActiveDocIndex(docIndex)
	}
	a.refreshExportPreview()
	a.refreshPlanillaState()
	return a.State()
}

func (a *App) JumpToPlanilla(planillaIndex int) coreapp.BatchSnapshot {
	a.core.JumpToPlanilla(planillaIndex)
	a.refreshPlanillaState()
	return a.State()
}

func (a *App) RenameDocumento(planillaIndex, docIndex int, name string) coreapp.BatchSnapshot {
	if err := a.core.SetDocumentName(planillaIndex, docIndex, name); err != nil {
		a.status = friendlyErr(err)
		return a.State()
	}
	a.status = "nombre actualizado"
	a.refreshExportPreview()
	return a.State()
}

func (a *App) OpenDocument(planillaIndex, docIndex int) coreapp.BatchSnapshot {
	docs := a.core.Planillas()
	if planillaIndex < 0 || planillaIndex >= len(docs) {
		a.status = "planilla fuera de rango"
		return a.State()
	}
	item := docs[planillaIndex]
	if docIndex < 0 || docIndex >= len(item.Documentos) {
		a.status = "documento fuera de rango"
		return a.State()
	}
	ruta := strings.TrimSpace(item.Documentos[docIndex].Ruta)
	if ruta == "" {
		a.status = "el documento no tiene ruta"
		return a.State()
	}
	if err := openNative(ruta); err != nil {
		a.status = friendlyErr(err)
		return a.State()
	}
	a.status = "Abriendo archivo: " + ruta
	a.core.SetActiveIndex(planillaIndex)
	a.core.SetActiveDocIndex(docIndex)
	a.refreshPlanillaState()
	a.preview = a.core.DocSummary(planillaIndex, docIndex)
	return a.State()
}

func (a *App) SetOutDir(path string) coreapp.BatchSnapshot {
	a.outDir = strings.TrimSpace(path)
	a.zipDest = ""
	return a.State()
}

func (a *App) SetZipDest(path string) coreapp.BatchSnapshot {
	a.zipDest = strings.TrimSpace(path)
	return a.State()
}

func (a *App) PrepareExport() coreapp.BatchSnapshot {
	if a.zipDest == "" {
		a.UpdateZipDestination()
	}
	a.exportInfo = a.core.ZipSummary(a.zipDest)
	return a.State()
}

func (a *App) Export() coreapp.BatchSnapshot {
	if a.zipDest == "" {
		a.UpdateZipDestination()
	}
	if err := a.core.ExportToPath(a.zipDest); err != nil {
		a.status = friendlyErr(err)
		a.exportInfo = friendlyErr(err)
		return a.State()
	}
	a.status = "archivo listo: " + a.zipDest
	a.exportInfo = "ZIP generado correctamente: " + a.zipDest
	return a.State()
}

func (a *App) SetAllowedNames(text string) coreapp.BatchSnapshot {
	cfg := a.core.Config()
	names := coreapp.NormalizeAllowedNames(coreapp.ParseNamesText(text))
	a.core.SetAllowedNames(names)
	cfg.AllowedNames = names
	_ = cfg.Save()
	a.status = "catálogo de nombres actualizado"
	return a.State()
}

func (a *App) SetRangeText(text string) coreapp.BatchSnapshot {
	cfg := a.core.Config()
	cfg.PlanillaRanges = coreapp.ParseRangesText(text)
	_ = cfg.Save()
	a.rangeInfo = "Rangos activos: " + strings.Join(cfg.PlanillaRanges, ", ")
	return a.State()
}

func (a *App) SetStatus(text string) coreapp.BatchSnapshot {
	a.status = text
	return a.State()
}

func (a *App) refreshPlanillaState() {
	a.batchInfo = fmt.Sprintf("Planillas cargadas: %d | Activa: %d", len(a.core.Planillas()), a.core.ActiveIndex()+1)
	if det := a.core.ActiveDetalle(); det != nil {
		a.preview = coreapp.PlanillaSummary(det)
		if len(det.Documentos) == 0 {
			a.docInfo = "Documento activo: ninguno"
			return
		}
		a.docInfo = fmt.Sprintf("Documento activo: %d de %d", a.core.ActiveDocIndex()+1, len(det.Documentos))
	}
}

func (a *App) refreshExportPreview() {
	if a.zipDest == "" {
		a.UpdateZipDestination()
	}
	if a.zipDest != "" {
		a.exportInfo = a.core.ZipSummary(a.zipDest)
	}
}

func parsePlanillaInputs(text string) []int64 {
	var out []int64
	seen := map[int64]struct{}{}
	fields := strings.FieldsFunc(text, func(r rune) bool {
		switch r {
		case '\n', '\r', '\t', ' ', ',', ';':
			return true
		}
		return false
	})
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		if strings.Contains(field, "-") {
			parts := strings.SplitN(field, "-", 2)
			if len(parts) != 2 {
				continue
			}
			min, err1 := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
			max, err2 := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
			if err1 != nil || err2 != nil {
				continue
			}
			if min > max {
				min, max = max, min
			}
			for n := min; n <= max; n++ {
				if _, ok := seen[n]; ok {
					continue
				}
				seen[n] = struct{}{}
				out = append(out, n)
			}
			continue
		}
		n, err := strconv.ParseInt(field, 10, 64)
		if err != nil {
			continue
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		out = append(out, n)
	}
	return out
}

func friendlyErr(err error) string {
	if err == nil {
		return ""
	}
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	switch {
	case strings.Contains(msg, "sesión oracle cerrada"):
		return "abre la sesión para continuar"
	case strings.Contains(msg, "ingresa usuario y contraseña"):
		return "ingresa usuario y contraseña"
	case strings.Contains(msg, "oracle") && strings.Contains(msg, "bloqueada"):
		return "la cuenta está bloqueada"
	case strings.Contains(msg, "oracle") && strings.Contains(msg, "vencida"):
		return "la contraseña está vencida"
	case strings.Contains(msg, "nombre duplicado en la planilla"):
		return "el mismo nombre ya está usado en esta planilla"
	case strings.Contains(msg, "no se puede guardar dentro de la carpeta"):
		return "elige una carpeta distinta para guardar el archivo"
	case strings.Contains(msg, "destino zip vacío"):
		return "elige una carpeta de destino"
	case strings.Contains(msg, "busca planillas primero"):
		return "primero busca una o más planillas"
	case strings.Contains(msg, "documento fuera de rango"), strings.Contains(msg, "planilla fuera de rango"):
		return "la selección ya no está disponible"
	case strings.Contains(msg, "no hay documentos seleccionados"):
		return "no hay documentos marcados"
	case strings.Contains(msg, "ingresa una o más planillas"):
		return "ingresa al menos una planilla"
	case strings.Contains(msg, "el lote supera el máximo permitido"):
		return "el lote supera el máximo permitido"
	default:
		if strings.Contains(msg, "oracle") {
			return "no se pudo abrir la sesión de datos"
		}
		return err.Error()
	}
}

func mustJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func openNative(path string) error {
	switch stdruntime.GOOS {
	case "windows":
		return exec.Command("cmd", "/c", "start", "", path).Start()
	case "darwin":
		return exec.Command("open", path).Start()
	default:
		return exec.Command("xdg-open", path).Start()
	}
}
