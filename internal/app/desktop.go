package app

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"multimediasc/internal/oracle"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type desktopState struct {
	cfg Config

	repo  *oracle.Repository
	items []batchItem

	activeIndex int
	docChecks   map[int64][]bool
	docNames    map[int64][]string

	previewPath string
	status      *widget.Label
	previewInfo *widget.Label
	exportInfo  *widget.Label
	batchInfo   *widget.Label
	docsBox     *fyne.Container
	planillaBox *fyne.Container
	planillaUI  *widget.List

	planillasText *widget.Entry
	batchFilter   *widget.Entry
	rangeText     *widget.Entry
	rangeInfo     *widget.Label
	namesText     *widget.Entry
	zipDest       *widget.Entry
	exportReady   string
	window        fyne.Window
}

type batchItem struct {
	Planilla *oracle.PlanillaDetalle
	Error    error
}

func RunDesktop() error {
	cfg := LoadEnvConfig()
	a := fyneapp.NewWithID("com.multimediasc.desktop")
	w := a.NewWindow("MultimediaSC")
	w.Resize(fyne.NewSize(1360, 860))

	state := &desktopState{
		cfg:       cfg,
		docChecks: map[int64][]bool{},
		docNames:  map[int64][]string{},
		window:    w,
	}

	loginUser := widget.NewEntry()
	loginPass := widget.NewPasswordEntry()
	if v := firstNonEmpty(os.Getenv("ORACLE_USER")); v != "" {
		loginUser.SetText(v)
	}
	if v := firstNonEmpty(os.Getenv("ORACLE_PASSWORD")); v != "" {
		loginPass.SetText(v)
	}
	searchBtn := widget.NewButton("Ingresar y buscar", func() {
		state.search(loginUser.Text, loginPass.Text)
	})
	logoutBtn := widget.NewButton("Cerrar sesión", func() {
		state.logout()
	})

	state.status = widget.NewLabel("Listo")
	state.previewInfo = widget.NewLabel("Busca planillas para ver documentos.")
	state.exportInfo = widget.NewLabel("Prepara la exportación para ver el resumen aquí.")
	state.batchInfo = widget.NewLabel("Planillas cargadas: 0")
	state.docsBox = container.NewVBox(widget.NewLabel("Sin resultados todavía"))
	state.rangeInfo = widget.NewLabel("Rangos activos: " + strings.Join(cfg.PlanillaRanges, ", "))
	state.planillasText = widget.NewMultiLineEntry()
	state.planillasText.SetPlaceHolder("Una planilla, un rango o una lista separada por comas o saltos de línea")
	state.batchFilter = widget.NewEntry()
	state.batchFilter.SetPlaceHolder("Filtrar planillas por número, paciente o error")
	state.batchFilter.OnChanged = func(string) {
		if state.planillaUI != nil {
			state.planillaUI.Refresh()
		}
	}

	state.rangeText = widget.NewMultiLineEntry()
	state.rangeText.SetPlaceHolder("Nombres base permitidos, uno por línea, sin extensión")
	state.rangeText.OnChanged = func(string) {
		state.cfg.AllowedNames = normalizeAllowedNames(parseNamesText(state.rangeText.Text))
		_ = state.cfg.Save()
		state.refreshDocs()
	}
	state.rangeText.SetText(cfg.AllowedNamesText())

	state.zipDest = widget.NewEntry()
	state.zipDest.SetPlaceHolder("Destino ZIP, por ejemplo C:\\Temp\\pl_202607231530.zip")

	searchPlanillas := widget.NewButton("Buscar planillas", func() {
		state.search(loginUser.Text, loginPass.Text)
	})
	reloadDocs := widget.NewButton("Recargar documentos", func() {
		state.refreshDocs()
	})
	prevPlanilla := widget.NewButton("Planilla anterior", func() {
		state.movePlanilla(-1)
	})
	nextPlanilla := widget.NewButton("Planilla siguiente", func() {
		state.movePlanilla(1)
	})

	leftPanel := container.NewVBox(
		widget.NewLabelWithStyle("Acceso Oracle", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewForm(
			widget.NewFormItem("Usuario", loginUser),
			widget.NewFormItem("Contraseña", loginPass),
		),
		widget.NewLabelWithStyle("Búsqueda de planillas", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		state.planillasText,
		widget.NewLabel("Escribe una planilla, un rango tipo 6076406-6076422 o varios valores separados por coma."),
		state.batchFilter,
		container.NewHBox(searchBtn, searchPlanillas, logoutBtn),
		container.NewHBox(prevPlanilla, nextPlanilla),
		widget.NewLabelWithStyle("Rangos / catálogos", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		state.rangeInfo,
		state.batchInfo,
		widget.NewLabel("Nombres base permitidos"),
		state.rangeText,
		widget.NewLabel("Se agrega la extensión original al exportar."),
	)

	rightPanel := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Documentos", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			state.previewInfo,
			container.NewAppTabs(
				container.NewTabItem("Revisar", container.NewVBox(
					reloadDocs,
					state.docsBox,
				)),
				container.NewTabItem("Exportar", container.NewVBox(
					widget.NewForm(widget.NewFormItem("Destino ZIP", state.zipDest)),
					state.exportInfo,
					container.NewHBox(
						widget.NewButton("Preparar ZIP", func() {
							state.prepareExport()
						}),
						widget.NewButton("Confirmar ZIP", func() {
							state.exportConfirmed()
						}),
					),
				)),
				container.NewTabItem("Resumen", container.NewVBox(
					widget.NewLabel("Selecciona una planilla de la lista para revisar sus documentos."),
					state.status,
				)),
			),
		),
		nil, nil, nil,
		container.NewVSplit(
			container.NewVScroll(container.NewVBox(
				widget.NewLabelWithStyle("Planillas", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				state.planillaList(),
			)),
			container.NewVScroll(container.NewVBox()),
		),
	)

	split := container.NewHSplit(container.NewVScroll(leftPanel), rightPanel)
	split.Offset = 0.34
	w.SetContent(split)
	w.ShowAndRun()
	return nil
}

func (s *desktopState) search(user, pass string) {
	if s.repo == nil {
		if strings.TrimSpace(user) == "" || pass == "" {
			s.status.SetText("ingresa usuario y contraseña para abrir la sesión")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		repo, err := oracle.Open(ctx, user, pass, oracle.OpenConfig{
			ConnectString: s.cfg.OracleConnect,
			MaxOpenConns:  s.cfg.DefaultMaxOpen,
			MaxIdleConns:  s.cfg.DefaultMaxIdle,
			ConnMaxLife:   10 * time.Minute,
		})
		if err != nil {
			s.status.SetText(friendlyErr(err))
			return
		}
		s.repo = repo
		s.status.SetText("sesión Oracle abierta")
	}

	inputs := parsePlanillaInputs(s.planillasText.Text)
	if len(inputs) == 0 {
		s.status.SetText("ingresa una o más planillas")
		return
	}
	if max := s.cfg.MaxBatchPlanillasOrDefault(); max > 0 && len(inputs) > max {
		s.status.SetText(fmt.Sprintf("el lote supera el máximo permitido (%d)", max))
		return
	}

	s.items = nil
	s.docChecks = map[int64][]bool{}
	s.docNames = map[int64][]string{}
	s.activeIndex = 0

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	for _, tramite := range inputs {
		det, err := s.repo.ObtenerDetallePlanilla(ctx, tramite)
		if err != nil {
			s.items = append(s.items, batchItem{Error: err})
			continue
		}
		enrichDocuments(det.Documentos)
		s.items = append(s.items, batchItem{Planilla: det})
		s.docChecks[det.Planilla.DigTramite] = make([]bool, len(det.Documentos))
		s.docNames[det.Planilla.DigTramite] = make([]string, len(det.Documentos))
		for i := range det.Documentos {
			if len(s.cfg.AllowedNames) > 0 {
				s.docNames[det.Planilla.DigTramite][i] = s.cfg.AllowedNames[0]
			}
		}
	}

	if len(s.items) == 0 {
		s.status.SetText("no se encontraron planillas")
		s.batchInfo.SetText("Planillas cargadas: 0")
		s.refreshPlanillaList()
		s.refreshDocs()
		return
	}
	s.refreshPlanillaList()
	s.refreshDocs()
	s.batchInfo.SetText(fmt.Sprintf("Planillas cargadas: %d | Activa: %d", len(s.items), s.activeIndex+1))
	s.status.SetText(fmt.Sprintf("cargadas %d planillas", len(s.items)))
}

func (s *desktopState) logout() {
	if s.repo != nil {
		_ = s.repo.Close()
		s.repo = nil
	}
	s.items = nil
	s.docChecks = map[int64][]bool{}
	s.docNames = map[int64][]string{}
	s.activeIndex = 0
	s.previewPath = ""
	s.previewInfo.SetText("sesión cerrada")
	s.status.SetText("sesión Oracle cerrada")
	s.batchInfo.SetText("Planillas cargadas: 0")
	s.refreshPlanillaList()
	s.refreshDocs()
}

func (s *desktopState) activeDetalle() *oracle.PlanillaDetalle {
	if s.activeIndex < 0 || s.activeIndex >= len(s.items) {
		return nil
	}
	return s.items[s.activeIndex].Planilla
}

func (s *desktopState) refreshPlanillaList() {
	if s.planillaUI == nil {
		return
	}
	s.planillaUI.Refresh()
	s.batchInfo.SetText(fmt.Sprintf("Planillas cargadas: %d | Activa: %d", len(s.items), s.activeIndex+1))
}

func (s *desktopState) planillaList() *widget.List {
	s.planillaUI = widget.NewList(
		func() int { return len(s.filteredPlanillaIndexes()) },
		func() fyne.CanvasObject {
			return widget.NewLabel("planilla")
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			indexes := s.filteredPlanillaIndexes()
			if i < 0 || i >= len(indexes) {
				label.SetText("")
				return
			}
			item := s.items[indexes[i]]
			if item.Error != nil {
				label.SetText(fmt.Sprintf("Error: %s", friendlyErr(item.Error)))
				return
			}
			det := item.Planilla
			selected := 0
			if checks := s.docChecks[det.Planilla.DigTramite]; len(checks) > 0 {
				for _, v := range checks {
					if v {
						selected++
					}
				}
			}
			label.SetText(fmt.Sprintf("%d - %s (%d docs, %d sel.)", det.Planilla.DigTramite, planillaPatientName(det), len(det.Documentos), selected))
			if i == s.activeIndex {
				label.TextStyle = fyne.TextStyle{Bold: true}
			}
		},
	)
	s.planillaUI.OnSelected = func(id widget.ListItemID) {
		indexes := s.filteredPlanillaIndexes()
		if id < 0 || id >= len(indexes) {
			return
		}
		s.activeIndex = indexes[id]
		s.batchInfo.SetText(fmt.Sprintf("Planillas cargadas: %d | Activa: %d", len(s.items), s.activeIndex+1))
		s.refreshDocs()
	}
	return s.planillaUI
}

func (s *desktopState) movePlanilla(delta int) {
	indexes := s.filteredPlanillaIndexes()
	if len(indexes) == 0 {
		return
	}
	currentPos := 0
	for i, idx := range indexes {
		if idx == s.activeIndex {
			currentPos = i
			break
		}
	}
	nextPos := currentPos + delta
	if nextPos < 0 {
		nextPos = 0
	}
	if nextPos >= len(indexes) {
		nextPos = len(indexes) - 1
	}
	next := indexes[nextPos]
	if next == s.activeIndex {
		return
	}
	s.activeIndex = next
	if s.planillaUI != nil {
		s.planillaUI.Select(s.activeIndex)
	}
	s.batchInfo.SetText(fmt.Sprintf("Planillas cargadas: %d | Activa: %d", len(s.items), s.activeIndex+1))
	s.refreshDocs()
}

func (s *desktopState) filteredPlanillaIndexes() []int {
	if s == nil {
		return nil
	}
	q := strings.ToLower(strings.TrimSpace(s.batchFilter.Text))
	if q == "" {
		indexes := make([]int, 0, len(s.items))
		for i := range s.items {
			indexes = append(indexes, i)
		}
		return indexes
	}
	var indexes []int
	for i, item := range s.items {
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

func (s *desktopState) refreshDocs() {
	det := s.activeDetalle()
	if det == nil {
		s.docsBox.Objects = []fyne.CanvasObject{widget.NewLabel("Sin resultados todavía")}
		s.docsBox.Refresh()
		return
	}
	tramite := det.Planilla.DigTramite
	checks := s.docChecks[tramite]
	names := s.docNames[tramite]
	var rows []fyne.CanvasObject
	for i, d := range det.Documentos {
		doc := d
		i := i
		checked := i < len(checks) && checks[i]
		nameText := ""
		if i < len(names) {
			nameText = names[i]
		}
		cb := widget.NewCheck("", func(v bool) {
			if i < len(s.docChecks[tramite]) {
				s.docChecks[tramite][i] = v
			}
		})
		cb.SetChecked(checked)
		nameSel := widget.NewSelect(s.cfg.AllowedNames, func(v string) {
			if i < len(s.docNames[tramite]) {
				s.docNames[tramite][i] = v
			}
		})
		if nameText != "" {
			nameSel.SetSelected(nameText)
		} else if len(s.cfg.AllowedNames) > 0 {
			nameSel.SetSelected(s.cfg.AllowedNames[0])
		}
		toggleBtn := widget.NewButton("", nil)
		if checked {
			toggleBtn.SetText("Incluido")
		} else {
			toggleBtn.SetText("Excluido")
		}
		toggleBtn.OnTapped = func() {
			if i < len(s.docChecks[tramite]) {
				s.docChecks[tramite][i] = !s.docChecks[tramite][i]
				if s.docChecks[tramite][i] {
					toggleBtn.SetText("Incluido")
				} else {
					toggleBtn.SetText("Excluido")
				}
			}
		}
		openBtn := widget.NewButton("Ver", func() {
			s.openDocument(doc)
		})
		actions := container.NewHBox(toggleBtn, openBtn, nameSel)
		row := container.NewBorder(nil, nil, cb, nil, container.NewVBox(
			widget.NewLabel(formatTime(doc.Fecha)),
			widget.NewLabel(doc.Descripcion),
			widget.NewLabel(documentStatusLine(doc)),
			actions,
		))
		rows = append(rows, row)
	}
	s.docsBox.Objects = rows
	s.docsBox.Refresh()
	s.previewInfo.SetText(planillaSummary(det))
}

func (s *desktopState) openDocument(doc oracle.ImagenPacienteRow) {
	s.previewPath = doc.Ruta
	s.previewInfo.SetText(documentSummary(doc))
	if s.window == nil {
		return
	}
	if _, err := os.Stat(doc.Ruta); err != nil {
		s.status.SetText(friendlyErr(err))
		return
	}
	u := &url.URL{Scheme: "file", Path: filepath.ToSlash(doc.Ruta)}
	if err := fyne.CurrentApp().OpenURL(u); err != nil {
		s.status.SetText(friendlyErr(err))
		return
	}
	s.status.SetText("Abriendo archivo: " + filepath.Base(doc.Ruta))
}

func (s *desktopState) prepareExport() {
	if len(s.items) == 0 {
		s.exportInfo.SetText("Busca planillas primero.")
		return
	}
	def := defaultZipFilename(s.exportPlanillaLabel(), s.selectedDocCount())
	if s.zipDest != nil && strings.TrimSpace(s.zipDest.Text) == "" {
		s.zipDest.SetText(def)
	}
	dest := strings.TrimSpace(s.zipDest.Text)
	if dest == "" {
		dest = def
		s.zipDest.SetText(dest)
	}
	dest = ensureZipExtension(dest)
	if err := s.validateZipDestination(dest); err != nil {
		s.exportInfo.SetText(friendlyErr(err))
		s.exportReady = ""
		return
	}
	s.exportReady = dest
	s.exportInfo.SetText(s.zipStructureSummary(dest))
}

func (s *desktopState) exportConfirmed() {
	if s.exportReady == "" {
		s.prepareExport()
		if s.exportReady == "" {
			return
		}
	}
	if err := s.exportToPath(s.exportReady); err != nil {
		s.status.SetText(friendlyErr(err))
		s.exportInfo.SetText(friendlyErr(err))
		return
	}
	s.status.SetText("ZIP generado: " + s.exportReady)
	s.exportInfo.SetText("ZIP generado correctamente: " + s.exportReady)
}

func (s *desktopState) exportToPath(dest string) error {
	if len(s.items) == 0 {
		return errors.New("busca planillas primero")
	}
	if err := s.validateZipDestination(dest); err != nil {
		return err
	}
	var docs []zipDocument
	lookup := map[string]string{}
	for _, item := range s.items {
		if item.Planilla == nil {
			continue
		}
		tramite := strconv.FormatInt(item.Planilla.Planilla.DigTramite, 10)
		checks := s.docChecks[item.Planilla.Planilla.DigTramite]
		names := s.docNames[item.Planilla.Planilla.DigTramite]
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
	if err := exportZip("", docs, lookup, dest, s.cfg.AllowDuplicateFix); err != nil {
		return err
	}
	s.status.SetText("ZIP generado: " + dest)
	return nil
}

func (s *desktopState) validateZipDestination(dest string) error {
	dest = strings.TrimSpace(dest)
	if dest == "" {
		return fmt.Errorf("destino ZIP vacío")
	}
	if len(s.items) == 0 {
		return fmt.Errorf("busca planillas primero")
	}
	base := strings.ToLower(filepath.Base(filepath.Dir(dest)))
	for _, item := range s.items {
		if item.Planilla == nil {
			continue
		}
		if base == strings.ToLower(strconv.FormatInt(item.Planilla.Planilla.DigTramite, 10)) {
			return fmt.Errorf("no se puede guardar dentro de la carpeta %s", filepath.Base(filepath.Dir(dest)))
		}
	}
	return nil
}

func ensureZipExtension(dest string) string {
	dest = strings.TrimSpace(dest)
	if dest == "" {
		return dest
	}
	if strings.EqualFold(filepath.Ext(dest), ".zip") {
		return dest
	}
	return dest + ".zip"
}

func (s *desktopState) selectedDocCount() int {
	total := 0
	for _, item := range s.items {
		if item.Planilla == nil {
			continue
		}
		for _, v := range s.docChecks[item.Planilla.Planilla.DigTramite] {
			if v {
				total++
			}
		}
	}
	return total
}

func (s *desktopState) exportPlanillaLabel() string {
	var ids []string
	for _, item := range s.items {
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

func (s *desktopState) exportNameForIndex(tramite int64, index int, doc oracle.ImagenPacienteRow) string {
	name := ""
	if names := s.docNames[tramite]; index >= 0 && index < len(names) {
		name = strings.TrimSpace(names[index])
	}
	if name == "" {
		name = filepath.Base(doc.Ruta)
	}
	if filepath.Ext(name) == "" {
		if ext := filepath.Ext(doc.Ruta); ext != "" {
			name += ext
		}
	}
	return name
}

func (s *desktopState) zipStructureSummary(dest string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Destino: %s\n", dest)
	fmt.Fprintf(&b, "Planillas: %d\n", len(s.items))
	fmt.Fprintf(&b, "Seleccionados: %d\n\n", s.selectedDocCount())
	for _, item := range s.items {
		if item.Planilla == nil {
			continue
		}
		tramite := item.Planilla.Planilla.DigTramite
		fmt.Fprintf(&b, "Carpeta %d/\n", tramite)
		checks := s.docChecks[tramite]
		for i, d := range item.Planilla.Documentos {
			if i >= len(checks) || !checks[i] {
				continue
			}
			fmt.Fprintf(&b, "  - %d/%s\n", tramite, s.exportNameForIndex(tramite, i, d))
		}
	}
	return b.String()
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

func planillaSummary(det *oracle.PlanillaDetalle) string {
	if det == nil {
		return ""
	}
	paciente := ""
	if det.Paciente != nil {
		paciente = det.Paciente.NombreCompleto()
	}
	return fmt.Sprintf(
		"Planilla: %d\nHC: %d\nPaciente: %s\nCédula: %s\nAseguradora: %s\nExpediente: %s\nServicio: %s\nEspecialidad: %s\nPermanencia: %d\nRepositorio: %s",
		det.Planilla.DigTramite,
		det.Planilla.DigHC,
		paciente,
		func() string {
			if det.Paciente != nil {
				return det.Paciente.Cedula
			}
			return ""
		}(),
		det.Planilla.DigAseguradora,
		det.Planilla.DigExpediente,
		det.Planilla.DigServicio,
		det.Planilla.DigEspecialidad,
		det.Planilla.DigNumeroPermanencia,
		det.Planilla.DigPathRepo,
	)
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

func documentStatusLine(doc oracle.ImagenPacienteRow) string {
	parts := []string{}
	if v := firstNonEmpty(doc.Kind, doc.Tipo); v != "" {
		parts = append(parts, v)
	}
	if v := doc.Ext; v != "" {
		parts = append(parts, v)
	}
	if v := doc.State; v != "" {
		parts = append(parts, v)
	}
	if v := doc.Reason; v != "" {
		parts = append(parts, v)
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " | ")
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
