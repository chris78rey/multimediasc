package app

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"multimediasc/internal/oracle"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type desktopState struct {
	cfg         Config
	repo        *oracle.Repository
	planilla    *oracle.PlanillaDetalle
	docChecks   []bool
	docNames    []string
	previewPath string
	status      *widget.Label
	previewInfo *widget.Label
	docsBox     *fyne.Container
	rangeText   *widget.Entry
	rangeInfo   *widget.Label
	window      fyne.Window
}

func RunDesktop() error {
	cfg := LoadEnvConfig()
	a := fyneapp.NewWithID("com.multimediasc.desktop")
	w := a.NewWindow("MultimediaSC")
	w.Resize(fyne.NewSize(1320, 840))

	state := &desktopState{cfg: cfg}
	state.window = w

	loginUser := widget.NewEntry()
	loginPass := widget.NewPasswordEntry()
	planilla := widget.NewEntry()
	planilla.SetPlaceHolder("Solo números")
	if v := firstNonEmpty(os.Getenv("ORACLE_USER")); v != "" {
		loginUser.SetText(v)
	}
	if v := firstNonEmpty(os.Getenv("ORACLE_PASSWORD")); v != "" {
		loginPass.SetText(v)
	}
	silentLogin := widget.NewCheck("Usar credenciales locales", func(bool) {})
	silentLogin.SetChecked(true)

	searchBtn := widget.NewButton("Ingresar y buscar", func() {
		state.search(w, loginUser.Text, loginPass.Text, planilla.Text)
	})
	logoutBtn := widget.NewButton("Cerrar sesión", func() {
		state.logout()
	})

	state.status = widget.NewLabel("Listo")
	state.previewInfo = widget.NewLabel("Busca una planilla para ver documentos.")
	state.docsBox = container.NewVBox(widget.NewLabel("Sin resultados todavía"))
	state.rangeInfo = widget.NewLabel("Rangos activos: " + strings.Join(cfg.PlanillaRanges, ", "))
	state.rangeText = widget.NewMultiLineEntry()
	state.rangeText.SetPlaceHolder("Una planilla por línea o un rango por línea, por ejemplo 6076406-6076422")
	state.rangeText.OnChanged = func(string) {
		state.cfg.PlanillaRanges = parseRangesText(state.rangeText.Text)
		state.rangeInfo.SetText("Rangos activos: " + strings.Join(state.cfg.PlanillaRanges, ", "))
		if len(state.cfg.PlanillaRanges) == 1 && isSinglePlanilla(state.cfg.PlanillaRanges[0]) && strings.TrimSpace(planilla.Text) == "" {
			planilla.SetText(strings.TrimSpace(state.cfg.PlanillaRanges[0]))
		}
		_ = state.cfg.Save()
	}
	state.rangeText.SetText(cfg.AllowedRangesText())

	leftPanel := container.NewVBox(
		widget.NewLabelWithStyle("Acceso Oracle", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewForm(
			widget.NewFormItem("Usuario", loginUser),
			widget.NewFormItem("Contraseña", loginPass),
			widget.NewFormItem("Planilla", planilla),
		),
		widget.NewLabelWithStyle("Rangos permitidos", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		state.rangeText,
		state.rangeInfo,
		silentLogin,
		searchBtn,
		logoutBtn,
		state.status,
	)

	rightPanel := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Documentos", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			state.previewInfo,
			widget.NewButton("Exportar ZIP", func() {
				state.exportWithDialog(w)
			}),
		),
		nil, nil, nil,
		container.NewVScroll(state.docsBox),
	)

	split := container.NewHSplit(container.NewVScroll(leftPanel), rightPanel)
	split.Offset = 0.28
	w.SetContent(split)
	w.ShowAndRun()
	return nil
}

func (s *desktopState) search(w fyne.Window, user, pass, planillaText string) {
	planillaText = strings.TrimSpace(planillaText)
	if planillaText == "" && len(s.cfg.PlanillaRanges) == 1 && isSinglePlanilla(s.cfg.PlanillaRanges[0]) {
		planillaText = strings.TrimSpace(s.cfg.PlanillaRanges[0])
	}
	if err := validatePlanilla(planillaText, s.cfg.PlanillaRanges); err != nil {
		s.status.SetText(err.Error())
		return
	}
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

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	tramite, _ := strconv.ParseInt(planillaText, 10, 64)
	det, err := s.repo.ObtenerDetallePlanilla(ctx, tramite)
	if err != nil {
		s.status.SetText(friendlyErr(err))
		return
	}
	enrichDocuments(det.Documentos)
	s.planilla = det
	s.docChecks = make([]bool, len(det.Documentos))
	s.docNames = make([]string, len(det.Documentos))
	for i := range det.Documentos {
		s.docChecks[i] = true
		if len(s.cfg.AllowedNames) > 0 {
			s.docNames[i] = s.cfg.AllowedNames[0]
		}
	}
	s.previewPath = ""
	s.status.SetText(fmt.Sprintf("Cargada planilla %s con %d documentos", planillaText, len(det.Documentos)))
	s.previewInfo.SetText(planillaSummary(det))
	s.refreshDocs()
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func (s *desktopState) logout() {
	if s.repo != nil {
		_ = s.repo.Close()
		s.repo = nil
	}
	s.planilla = nil
	s.docChecks = nil
	s.docNames = nil
	s.previewPath = ""
	s.previewInfo.SetText("sesión cerrada")
	s.status.SetText("sesión Oracle cerrada")
	if s.docsBox != nil {
		s.docsBox.Objects = []fyne.CanvasObject{widget.NewLabel("Sin resultados todavía")}
		s.docsBox.Refresh()
	}
}

func (s *desktopState) refreshDocs() {
	if s.planilla == nil {
		s.docsBox.Objects = []fyne.CanvasObject{widget.NewLabel("Sin documentos")}
		s.docsBox.Refresh()
		return
	}
	rows := make([]fyne.CanvasObject, 0, len(s.planilla.Documentos))
	for i, d := range s.planilla.Documentos {
		i := i
		doc := d
		cb := widget.NewCheck("", func(v bool) { s.docChecks[i] = v })
		cb.SetChecked(true)
		nameSel := widget.NewSelect(s.cfg.AllowedNames, func(v string) { s.docNames[i] = v })
		if len(s.cfg.AllowedNames) > 0 {
			nameSel.SetSelected(s.cfg.AllowedNames[0])
		}
		openBtn := widget.NewButton("Ver", func() {
			s.openDocument(doc)
		})
		row := container.NewBorder(nil, nil, cb, openBtn, container.NewVBox(
			widget.NewLabel(formatTime(doc.Fecha)),
			widget.NewLabel(doc.Descripcion),
			widget.NewLabel(documentStatusLine(doc)),
			nameSel,
		))
		rows = append(rows, row)
	}
	s.docsBox.Objects = rows
	s.docsBox.Refresh()
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

func (s *desktopState) exportWithDialog(w fyne.Window) {
	if s.planilla == nil {
		dialog.ShowInformation("Exportar ZIP", "Busca una planilla primero.", w)
		return
	}
	def := fmt.Sprintf("%s.zip", strconv.FormatInt(s.planilla.Planilla.DigTramite, 10))
	dialog.ShowFileSave(func(uc fyne.URIWriteCloser, err error) {
		if err != nil || uc == nil {
			return
		}
		dest := ensureZipExtension(uc.URI().Path())
		_ = uc.Close()
		if err := s.exportToPath(dest); err != nil {
			s.status.SetText(friendlyErr(err))
			dialog.ShowError(err, w)
		}
	}, w)
	s.status.SetText("Elige dónde guardar " + def)
}

func (s *desktopState) exportToPath(dest string) error {
	if s.planilla == nil {
		return errors.New("busca una planilla primero")
	}
	lookup := map[string]string{}
	var docs []zipDocument
	for i, d := range s.planilla.Documentos {
		lookup[strconv.Itoa(i)] = d.Ruta
		docs = append(docs, zipDocument{
			Index:     i,
			Selected:  s.docChecks[i],
			ExportAs:  s.docNames[i],
			ValidName: true,
		})
	}
	if err := exportZip(strconv.FormatInt(s.planilla.Planilla.DigTramite, 10), docs, lookup, dest, s.cfg.AllowDuplicateFix); err != nil {
		return err
	}
	s.status.SetText("ZIP generado: " + dest)
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

func isSinglePlanilla(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" || strings.Contains(s, "-") {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04")
}
