package app

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"multimediasc/internal/oracle"
)

// BatchCore concentra la lógica de negocio sin depender de Fyne.
// La UI actual y una futura UI Wails pueden usar este núcleo.
type BatchCore struct {
	cfg          Config
	repo         *oracle.Repository
	items        []batchItem
	activeIndex  int
	activeDocIdx int
	docChecks    map[int64][]bool
	docNames     map[int64][]string
}

type batchItem struct {
	Planilla *oracle.PlanillaDetalle
	Error    error
}

type PlanillaView struct {
	Index        int             `json:"index"`
	Tramite      int64           `json:"tramite"`
	HC           int64           `json:"hc"`
	Paciente     string          `json:"paciente"`
	Cedula       string          `json:"cedula"`
	TotalDocs    int             `json:"total_docs"`
	SelectedDocs int             `json:"selected_docs"`
	Error        string          `json:"error"`
	Documentos   []DocumentoView `json:"documentos"`
}

type DocumentoView struct {
	Index        int    `json:"index"`
	Planilla     int64  `json:"planilla"`
	Descripcion  string `json:"descripcion"`
	Fecha        string `json:"fecha"`
	Tipo         string `json:"tipo"`
	Estado       string `json:"estado"`
	Motivo       string `json:"motivo"`
	Nombre       string `json:"nombre"`
	Ruta         string `json:"ruta"`
	Seleccionado bool   `json:"seleccionado"`
	Activo       bool   `json:"activo"`
}

type BatchSnapshot struct {
	Status          string         `json:"status"`
	LoggedIn        bool           `json:"logged_in"`
	PreviewInfo     string         `json:"preview_info"`
	DocInfo         string         `json:"doc_info"`
	ExportInfo      string         `json:"export_info"`
	BatchInfo       string         `json:"batch_info"`
	RangeInfo       string         `json:"range_info"`
	ZipDest         string         `json:"zip_dest"`
	OutDir          string         `json:"out_dir"`
	ActivePlanilla  int            `json:"active_planilla"`
	ActiveDocumento int            `json:"active_documento"`
	Planillas       []PlanillaView `json:"planillas"`
	SelectedDocs    int            `json:"selected_docs"`
	AllowedNames    []string       `json:"allowed_names"`
}

func NewBatchCore(cfg Config) *BatchCore {
	return &BatchCore{
		cfg:       cfg,
		docChecks: map[int64][]bool{},
		docNames:  map[int64][]string{},
	}
}

func (b *BatchCore) Config() Config { return b.cfg }

func (b *BatchCore) SetAllowedNames(names []string) {
	b.cfg.AllowedNames = normalizeAllowedNames(names)
}

func (b *BatchCore) Close() error {
	if b.repo == nil {
		return nil
	}
	err := b.repo.Close()
	b.repo = nil
	return err
}

func (b *BatchCore) Login(ctx context.Context, user, pass string) error {
	if b.repo != nil {
		return nil
	}
	if strings.TrimSpace(user) == "" || pass == "" {
		return errors.New("ingresa usuario y contraseña para abrir la sesión")
	}
	repo, err := oracle.Open(ctx, user, pass, oracle.OpenConfig{
		ConnectString: b.cfg.OracleConnect,
		MaxOpenConns:  b.cfg.DefaultMaxOpen,
		MaxIdleConns:  b.cfg.DefaultMaxIdle,
		ConnMaxLife:   10 * time.Minute,
	})
	if err != nil {
		return err
	}
	b.repo = repo
	return nil
}

func (b *BatchCore) Logout() {
	_ = b.Close()
	b.items = nil
	b.activeDocIdx = 0
	b.docChecks = map[int64][]bool{}
	b.docNames = map[int64][]string{}
	b.activeIndex = 0
}

func (b *BatchCore) Search(ctx context.Context, inputs []int64) ([]batchItem, error) {
	if b.repo == nil {
		return nil, errors.New("sesión Oracle cerrada")
	}
	if len(inputs) == 0 {
		return nil, errors.New("ingresa una o más planillas")
	}
	if max := b.cfg.MaxBatchPlanillasOrDefault(); max > 0 && len(inputs) > max {
		return nil, fmt.Errorf("el lote supera el máximo permitido (%d)", max)
	}

	b.items = nil
	b.activeIndex = 0
	b.activeDocIdx = 0
	b.docChecks = map[int64][]bool{}
	b.docNames = map[int64][]string{}

	ctx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()
	for _, tramite := range inputs {
		det, err := b.repo.ObtenerDetallePlanilla(ctx, tramite)
		if err != nil {
			b.items = append(b.items, batchItem{Error: err})
			continue
		}
		enrichDocuments(det.Documentos)
		b.items = append(b.items, batchItem{Planilla: det})
		b.docChecks[det.Planilla.DigTramite] = make([]bool, len(det.Documentos))
		b.docNames[det.Planilla.DigTramite] = make([]string, len(det.Documentos))
	}
	return b.items, nil
}

func (b *BatchCore) Items() []batchItem { return b.items }

func (b *BatchCore) Planillas() []PlanillaView {
	out := make([]PlanillaView, 0, len(b.items))
	for i, item := range b.items {
		view := PlanillaView{Index: i}
		if item.Error != nil {
			view.Error = friendlyErr(item.Error)
			out = append(out, view)
			continue
		}
		det := item.Planilla
		view.Tramite = det.Planilla.DigTramite
		view.HC = det.Planilla.DigHC
		view.Paciente = planillaPatientName(det)
		if det.Paciente != nil {
			view.Cedula = det.Paciente.Cedula
		}
		view.TotalDocs = len(det.Documentos)
		if checks := b.docChecks[det.Planilla.DigTramite]; len(checks) > 0 {
			for _, v := range checks {
				if v {
					view.SelectedDocs++
				}
			}
		}
		view.Documentos = b.documentosForItem(i, det)
		out = append(out, view)
	}
	return out
}

func (b *BatchCore) documentosForItem(planillaIndex int, det *oracle.PlanillaDetalle) []DocumentoView {
	out := make([]DocumentoView, 0, len(det.Documentos))
	if det == nil {
		return out
	}
	tramite := det.Planilla.DigTramite
	checks := b.docChecks[tramite]
	for i, doc := range det.Documentos {
		out = append(out, DocumentoView{
			Index:        i,
			Planilla:     tramite,
			Descripcion:  doc.Descripcion,
			Fecha:        formatTime(doc.Fecha),
			Tipo:         firstNonEmpty(doc.Tipo, doc.Kind),
			Estado:       firstNonEmpty(doc.State, "desconocido"),
			Motivo:       firstNonEmpty(doc.Reason, "sin observaciones"),
			Nombre:       b.ExportNameForIndex(tramite, i, doc),
			Ruta:         doc.Ruta,
			Seleccionado: i < len(checks) && checks[i],
			Activo:       b.activeIndex == planillaIndex && i == b.activeDocIdx,
		})
	}
	return out
}

func (b *BatchCore) ActiveIndex() int { return b.activeIndex }

func (b *BatchCore) ActiveDetalle() *oracle.PlanillaDetalle {
	if b.activeIndex < 0 || b.activeIndex >= len(b.items) {
		return nil
	}
	return b.items[b.activeIndex].Planilla
}

func (b *BatchCore) SetActiveIndex(i int) {
	if i < 0 || i >= len(b.items) {
		return
	}
	b.activeIndex = i
	if det := b.ActiveDetalle(); det != nil && len(det.Documentos) > 0 && b.activeDocIdx >= len(det.Documentos) {
		b.activeDocIdx = 0
	}
}

func (b *BatchCore) ActiveDocIndex() int { return b.activeDocIdx }

func (b *BatchCore) SetActiveDocIndex(i int) {
	det := b.ActiveDetalle()
	if det == nil || len(det.Documentos) == 0 {
		b.activeDocIdx = 0
		return
	}
	if i < 0 {
		i = 0
	}
	if i >= len(det.Documentos) {
		i = len(det.Documentos) - 1
	}
	b.activeDocIdx = i
}

func (b *BatchCore) VisibleDocIndex() int {
	det := b.ActiveDetalle()
	if det == nil || len(det.Documentos) == 0 {
		return -1
	}
	if idx := b.nextPendingDocIndex(b.activeDocIdx, 1); idx >= 0 {
		return idx
	}
	if idx := b.nextPendingDocIndex(b.activeDocIdx, -1); idx >= 0 {
		return idx
	}
	return -1
}

func (b *BatchCore) SetSelected(planillaIndex, docIndex int, selected bool) {
	if planillaIndex < 0 || planillaIndex >= len(b.items) {
		return
	}
	item := b.items[planillaIndex]
	if item.Planilla == nil || docIndex < 0 || docIndex >= len(item.Planilla.Documentos) {
		return
	}
	tramite := item.Planilla.Planilla.DigTramite
	checks := b.docChecks[tramite]
	if len(checks) <= docIndex {
		tmp := make([]bool, len(item.Planilla.Documentos))
		copy(tmp, checks)
		checks = tmp
	}
	checks[docIndex] = selected
	b.docChecks[tramite] = checks
}

func (b *BatchCore) AdvanceToNextPending() int {
	det := b.ActiveDetalle()
	if det == nil || len(det.Documentos) == 0 {
		return b.activeDocIdx
	}
	next := b.nextPendingDocIndex(b.activeDocIdx, 1)
	if next < 0 {
		next = b.nextPendingDocIndex(b.activeDocIdx, -1)
	}
	if next >= 0 {
		b.activeDocIdx = next
	}
	return b.activeDocIdx
}

func (b *BatchCore) SetDocumentName(planillaIndex, docIndex int, name string) error {
	if planillaIndex < 0 || planillaIndex >= len(b.items) {
		return fmt.Errorf("planilla fuera de rango")
	}
	item := b.items[planillaIndex]
	if item.Planilla == nil || docIndex < 0 || docIndex >= len(item.Planilla.Documentos) {
		return fmt.Errorf("documento fuera de rango")
	}
	tramite := item.Planilla.Planilla.DigTramite
	name = strings.TrimSpace(name)
	if name != "" && b.nameUsedByOtherDoc(tramite, docIndex, name) {
		return fmt.Errorf("el nombre %q ya está usado en esta planilla", name)
	}
	names := b.docNames[tramite]
	if len(names) <= docIndex {
		tmp := make([]string, len(item.Planilla.Documentos))
		copy(tmp, names)
		names = tmp
	}
	names[docIndex] = name
	b.docNames[tramite] = names
	return nil
}

func (b *BatchCore) ClearSelection() {
	for tramite := range b.docChecks {
		checks := b.docChecks[tramite]
		for i := range checks {
			checks[i] = false
		}
		b.docChecks[tramite] = checks
	}
}

func (b *BatchCore) RemoveSelection(planillaIndex, docIndex int) {
	b.SetSelected(planillaIndex, docIndex, false)
}

func (b *BatchCore) JumpToPlanilla(planillaIndex int) {
	b.SetActiveIndex(planillaIndex)
	if det := b.ActiveDetalle(); det != nil {
		b.activeDocIdx = 0
		if len(det.Documentos) == 0 {
			b.activeDocIdx = 0
		}
	}
}

func (b *BatchCore) DocSummary(planillaIndex, docIndex int) string {
	if planillaIndex < 0 || planillaIndex >= len(b.items) {
		return ""
	}
	item := b.items[planillaIndex]
	if item.Planilla == nil || docIndex < 0 || docIndex >= len(item.Planilla.Documentos) {
		return ""
	}
	return documentSummary(item.Planilla.Documentos[docIndex])
}

func (b *BatchCore) ZipSummary(dest string) string { return b.ZipStructureSummary(dest) }

func (b *BatchCore) Snapshot(status, previewInfo, docInfo, exportInfo, batchInfo, rangeInfo, zipDest, outDir string) BatchSnapshot {
	return BatchSnapshot{
		Status:          status,
		LoggedIn:        b.repo != nil,
		PreviewInfo:     previewInfo,
		DocInfo:         docInfo,
		ExportInfo:      exportInfo,
		BatchInfo:       batchInfo,
		RangeInfo:       rangeInfo,
		ZipDest:         zipDest,
		OutDir:          outDir,
		ActivePlanilla:  b.activeIndex,
		ActiveDocumento: b.activeDocIdx,
		Planillas:       b.Planillas(),
		SelectedDocs:    b.SelectedDocCount(),
		AllowedNames:    append([]string(nil), b.cfg.AllowedNames...),
	}
}

func (b *BatchCore) MovePlanilla(delta int, filtered []int) int {
	if len(filtered) == 0 {
		return b.activeIndex
	}
	currentPos := 0
	for i, idx := range filtered {
		if idx == b.activeIndex {
			currentPos = i
			break
		}
	}
	nextPos := currentPos + delta
	if nextPos < 0 {
		nextPos = 0
	}
	if nextPos >= len(filtered) {
		nextPos = len(filtered) - 1
	}
	b.SetActiveIndex(filtered[nextPos])
	return b.activeIndex
}

func (b *BatchCore) MoveDocumento(delta int) int {
	det := b.ActiveDetalle()
	if det == nil || len(det.Documentos) == 0 {
		return b.activeDocIdx
	}
	next := b.nextPendingDocIndex(b.activeDocIdx, delta)
	if next >= 0 {
		b.activeDocIdx = next
	}
	return b.activeDocIdx
}

func (b *BatchCore) nextPendingDocIndex(from, delta int) int {
	det := b.ActiveDetalle()
	if det == nil || len(det.Documentos) == 0 {
		return -1
	}
	checks := b.docChecks[det.Planilla.DigTramite]
	if delta == 0 {
		delta = 1
	}
	start := from + delta
	if start < 0 {
		start = 0
	}
	if start >= len(det.Documentos) {
		start = len(det.Documentos) - 1
	}
	step := 1
	if delta < 0 {
		step = -1
	}
	for i := start; i >= 0 && i < len(det.Documentos); i += step {
		if i >= len(checks) || !checks[i] {
			return i
		}
	}
	return -1
}

func (b *BatchCore) FilteredPlanillaIndexes(query string) []int {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		indexes := make([]int, 0, len(b.items))
		for i := range b.items {
			indexes = append(indexes, i)
		}
		return indexes
	}
	var indexes []int
	for i, item := range b.items {
		if item.Error != nil {
			if strings.Contains(strings.ToLower(friendlyErr(item.Error)), q) {
				indexes = append(indexes, i)
			}
			continue
		}
		det := item.Planilla
		if strings.Contains(strconv.FormatInt(det.Planilla.DigTramite, 10), q) ||
			strings.Contains(strings.ToLower(planillaPatientName(det)), q) ||
			strings.Contains(strings.ToLower(det.Planilla.DigCedula), q) ||
			strings.Contains(strings.ToLower(det.Planilla.DigServicio), q) ||
			strings.Contains(strings.ToLower(det.Planilla.DigEspecialidad), q) {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func (b *BatchCore) SelectedDocCount() int {
	total := 0
	for _, item := range b.items {
		if item.Planilla == nil {
			continue
		}
		for _, v := range b.docChecks[item.Planilla.Planilla.DigTramite] {
			if v {
				total++
			}
		}
	}
	return total
}

func (b *BatchCore) ExportPlanillaLabel() string {
	var ids []string
	for _, item := range b.items {
		if item.Planilla == nil {
			continue
		}
		ids = append(ids, strconv.FormatInt(item.Planilla.Planilla.DigTramite, 10))
	}
	switch len(ids) {
	case 0:
		return ""
	case 1:
		return ids[0]
	default:
		sort.Strings(ids)
		return "pl_" + time.Now().Format("200601021504")
	}
}

func (b *BatchCore) ExportNameForIndex(tramite int64, index int, doc oracle.ImagenPacienteRow) string {
	name := ""
	if names := b.docNames[tramite]; index >= 0 && index < len(names) {
		name = strings.TrimSpace(names[index])
	}
	if name == "" {
		name = basenameAnySeparator(doc.Ruta)
	}
	if filepath.Ext(name) == "" {
		if ext := filepath.Ext(doc.Ruta); ext != "" {
			name += ext
		}
	}
	return name
}

func (b *BatchCore) AllowedNames() []string {
	return append([]string(nil), b.cfg.AllowedNames...)
}

func (b *BatchCore) nameUsedByOtherDoc(tramite int64, docIndex int, name string) bool {
	name = strings.ToLower(strings.TrimSpace(strings.TrimSuffix(name, filepath.Ext(name))))
	if name == "" {
		return false
	}
	names := b.docNames[tramite]
	for i, other := range names {
		if i == docIndex {
			continue
		}
		other = strings.ToLower(strings.TrimSpace(strings.TrimSuffix(other, filepath.Ext(other))))
		if other != "" && other == name {
			return true
		}
	}
	return false
}

func (b *BatchCore) ZipStructureSummary(dest string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Destino: %s\n", dest)
	fmt.Fprintf(&sb, "Planillas: %d\n", len(b.items))
	fmt.Fprintf(&sb, "Seleccionados: %d\n\n", b.SelectedDocCount())
	for _, item := range b.items {
		if item.Planilla == nil {
			continue
		}
		tramite := item.Planilla.Planilla.DigTramite
		fmt.Fprintf(&sb, "Carpeta %d/\n", tramite)
		checks := b.docChecks[tramite]
		for i, d := range item.Planilla.Documentos {
			if i >= len(checks) || !checks[i] {
				continue
			}
			fmt.Fprintf(&sb, "  - %d/%s\n", tramite, b.ExportNameForIndex(tramite, i, d))
		}
	}
	return sb.String()
}

func (b *BatchCore) ExportToPath(dest string) error {
	if len(b.items) == 0 {
		return errors.New("busca planillas primero")
	}
	if err := validateZipDestination(dest, b.items); err != nil {
		return err
	}
	var docs []zipDocument
	lookup := map[string]string{}
	for _, item := range b.items {
		if item.Planilla == nil {
			continue
		}
		tramite := strconv.FormatInt(item.Planilla.Planilla.DigTramite, 10)
		checks := b.docChecks[item.Planilla.Planilla.DigTramite]
		names := b.docNames[item.Planilla.Planilla.DigTramite]
		for i, d := range item.Planilla.Documentos {
			lookup[strconv.Itoa(len(docs))] = d.Ruta
			exportAs := ""
			if i < len(names) {
				exportAs = names[i]
			}
			selected := i < len(checks) && checks[i]
			docs = append(docs, zipDocument{
				Index:     len(docs),
				Selected:  selected,
				ExportAs:  exportAs,
				ValidName: true,
				Planilla:  tramite,
				Ruta:      d.Ruta,
			})
		}
	}
	return exportZip("", docs, lookup, dest, b.cfg.AllowDuplicateFix)
}

func planillaPatientName(det *oracle.PlanillaDetalle) string {
	if det == nil || det.Paciente == nil {
		return ""
	}
	return det.Paciente.NombreCompleto()
}

func documentSummary(doc oracle.ImagenPacienteRow) string {
	return fmt.Sprintf(
		"%s\n%s\nTipo: %s  Ext: %s  Estado: %s\nMotivo: %s",
		doc.Descripcion,
		doc.Ruta,
		firstNonEmpty(doc.Tipo, doc.Kind),
		firstNonEmpty(doc.Ext, "-"),
		firstNonEmpty(doc.State, "desconocido"),
		firstNonEmpty(doc.Reason, "sin observaciones"),
	)
}

func planillaSummary(det *oracle.PlanillaDetalle) string {
	if det == nil {
		return ""
	}
	paciente := ""
	if det.Paciente != nil {
		paciente = det.Paciente.NombreCompleto()
	}
	return fmt.Sprintf(
		"Planilla: %d\nHC: %d\nPaciente: %s",
		det.Planilla.DigTramite,
		det.Planilla.DigHC,
		paciente,
	)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04")
}

func validateZipDestination(dest string, items []batchItem) error {
	dest = strings.TrimSpace(dest)
	if dest == "" {
		return fmt.Errorf("destino ZIP vacío")
	}
	if len(items) == 0 {
		return fmt.Errorf("busca planillas primero")
	}
	base := strings.ToLower(filepath.Base(filepath.Dir(dest)))
	for _, item := range items {
		if item.Planilla == nil {
			continue
		}
		if base == strings.ToLower(strconv.FormatInt(item.Planilla.Planilla.DigTramite, 10)) {
			return fmt.Errorf("no se puede guardar dentro de la carpeta %s", filepath.Base(filepath.Dir(dest)))
		}
	}
	return nil
}

func PlanillaSummary(det *oracle.PlanillaDetalle) string {
	return planillaSummary(det)
}

func ParseNamesText(text string) []string {
	return parseNamesText(text)
}

func ParseRangesText(text string) []string {
	return parseRangesText(text)
}

func NormalizeAllowedNames(names []string) []string {
	return normalizeAllowedNames(names)
}

func DefaultZipPath(outDir, planillaLabel string, selectedCount int) string {
	outDir = strings.TrimSpace(outDir)
	if outDir == "" {
		outDir = DefaultAutoOutDir()
	}
	name := defaultZipFilename(planillaLabel, selectedCount)
	return filepath.Join(outDir, name)
}

func DefaultAutoOutDir() string {
	return filepath.Join(".", defaultExportFolderName())
}
