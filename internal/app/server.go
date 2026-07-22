package app

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"multimediasc/internal/oracle"
)

//go:embed templates/*
var templatesFS embed.FS

type Server struct {
	cfg    Config
	tmpl   *template.Template
	mux    *http.ServeMux
	client *http.Client
}

type pageData struct {
	OracleConnect string
	AllowedNames  []string
	Error         string
	Notice        string
	Result        *oracle.PlanillaDetalle
	Planilla      string
	Selected      map[string]string
}

type zipDocument struct {
	Index     int    `json:"index"`
	Selected  bool   `json:"selected"`
	ExportAs  string `json:"export_as"`
	ValidName bool   `json:"valid_name"`
}

func NewServer(cfg Config) *Server {
	tmpl := template.Must(template.ParseFS(templatesFS, "templates/*.html"))
	s := &Server{
		cfg:  cfg,
		tmpl: tmpl,
		mux:  http.NewServeMux(),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	s.routes()
	return s
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/api/search", s.handleSearch)
	s.mux.HandleFunc("/api/preview", s.handlePreview)
	s.mux.HandleFunc("/api/export", s.handleExport)
	s.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("ok")) })
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	data := pageData{
		OracleConnect: s.cfg.OracleConnect,
		AllowedNames:  s.cfg.AllowedNames,
		Selected:      map[string]string{},
	}
	_ = s.tmpl.ExecuteTemplate(w, "index.html", data)
}

func (s *Server) openRepo(r *http.Request, user, pass string) (*oracle.Repository, error) {
	connect := r.FormValue("connect")
	if connect == "" {
		connect = s.cfg.OracleConnect
	}
	if connect == "" {
		return nil, fmt.Errorf("falta ORACLE_CONNECT")
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	return oracle.Open(ctx, user, pass, oracle.OpenConfig{
		ConnectString: connect,
		MaxOpenConns:  s.cfg.DefaultMaxOpen,
		MaxIdleConns:  s.cfg.DefaultMaxIdle,
		ConnMaxLife:   10 * time.Minute,
	})
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	planilla := strings.TrimSpace(r.FormValue("planilla"))
	user := strings.TrimSpace(r.FormValue("user"))
	pass := r.FormValue("password")

	data := pageData{
		OracleConnect: r.FormValue("connect"),
		AllowedNames:  s.cfg.AllowedNames,
		Planilla:      planilla,
		Selected:      map[string]string{},
	}
	if err := validatePlanilla(planilla); err != nil {
		data.Error = err.Error()
		_ = s.tmpl.ExecuteTemplate(w, "index.html", data)
		return
	}

	repo, err := s.openRepo(r, user, pass)
	if err != nil {
		data.Error = friendlyErr(err)
		_ = s.tmpl.ExecuteTemplate(w, "index.html", data)
		return
	}
	defer repo.Close()

	tramite, _ := strconv.ParseInt(planilla, 10, 64)
	det, err := repo.ObtenerDetallePlanilla(r.Context(), tramite)
	if err != nil {
		data.Error = friendlyErr(err)
		_ = s.tmpl.ExecuteTemplate(w, "index.html", data)
		return
	}
	data.Result = det
	data.Notice = fmt.Sprintf("Se cargaron %d documentos", len(det.Documentos))
	_ = s.tmpl.ExecuteTemplate(w, "index.html", data)
}

func (s *Server) handlePreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "missing path", http.StatusBadRequest)
		return
	}
	f, err := os.Open(path)
	if err != nil {
		http.Error(w, friendlyErr(err), http.StatusNotFound)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline")
	_, _ = io.Copy(w, f)
}

func (s *Server) handleExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Planilla   string            `json:"planilla"`
		Documents  []zipDocument     `json:"documents"`
		PathLookup map[string]string `json:"path_lookup"`
		Dest       string            `json:"dest"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "json inválido", http.StatusBadRequest)
		return
	}
	if err := validatePlanilla(req.Planilla); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Dest == "" {
		http.Error(w, "destino requerido", http.StatusBadRequest)
		return
	}
	if len(req.Documents) == 0 {
		http.Error(w, "sin documentos", http.StatusBadRequest)
		return
	}
	if err := exportZip(req.Planilla, req.Documents, req.PathLookup, req.Dest, s.cfg.AllowDuplicateFix); err != nil {
		http.Error(w, friendlyErr(err), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func validatePlanilla(s string) error {
	if s == "" {
		return fmt.Errorf("planilla vacía")
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return fmt.Errorf("la planilla debe contener solo números")
		}
	}
	if len(s) < 3 {
		return fmt.Errorf("la planilla es demasiado corta")
	}
	return nil
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

func errorsIs(err, target error) bool {
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
		name string
		path string
	}
	var items []pair
	seen := map[string]int{}
	for _, d := range docs {
		if !d.Selected {
			continue
		}
		path := lookup[strconv.Itoa(d.Index)]
		if path == "" {
			continue
		}
		name := d.ExportAs
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
		items = append(items, pair{name: name, path: path})
	}
	if len(items) == 0 {
		return fmt.Errorf("no hay documentos seleccionados")
	}
	sort.Slice(items, func(i, j int) bool { return items[i].name < items[j].name })

	buf := &bytes.Buffer{}
	zw := newZipWriter(buf)
	for _, it := range items {
		if err := zw.addFile(filepath.Join(planilla, it.name), it.path); err != nil {
			return err
		}
	}
	if err := zw.close(); err != nil {
		return err
	}
	_, err := io.Copy(w, buf)
	return err
}

func validWindowsFilename(name string) bool {
	if name == "" || !strings.HasSuffix(strings.ToLower(name), ".pdf") {
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
