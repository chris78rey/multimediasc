package oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/sijms/go-ora/v2"
	_ "github.com/sijms/go-ora/v2"
)

const (
	QueryBuscarPlanilla = `
SELECT
  d.DIG_ID,
  d.DIG_TRAMITE,
  d.DIG_ASEGURADORA,
  d.DIG_FECHA_PLANILLA,
  d.DIG_FECHA_HASTA,
  d.DIG_PROCESADO,
  d.DIG_FECHA_PROCESO,
  d.DIG_HC,
  d.FE_PLA_ANIOMES,
  d.DIG_EXPEDIENTE,
  d.DIG_ANIO,
  d.DIG_CAMBIO_NOMBRE,
  d.DIG_FIRMAR,
  d.DIG_AREA_DEP,
  d.REVISADA_CONTRL_DOC,
  d.DIG_SERVICIO,
  d.DIG_ESPECIALIDAD,
  d.DIG_CEDULA,
  d.DIG_MENOR_EDAD,
  d.DIG_PLANILLADO,
  d.DIG_COBERTURA,
  d.DIG_ID_GENERACION,
  d.DIG_ID_TIPO,
  d.DIG_BLOQUEO_SGH,
  d.DIG_ID_TRAMITE,
  d.DIG_USUARIO,
  d.DIG_CODIGO_USUARIO,
  d.DIG_FECHA_ALTA,
  d.DIG_NUMERO_PERMANENCIA,
  d.DIG_FECHA_COBERTURA,
  d.DIG_PATH_REPO
FROM DIGITALIZACION d
WHERE d.DIG_TRAMITE = :1`

	QueryBuscarPaciente = `
SELECT
  p.NUMERO_HC,
  p.APELLIDO_PATERNO,
  p.APELLIDO_MATERNO,
  p.PRIMER_NOMBRE,
  p.SEGUNDO_NOMBRE,
  p.CEDULA
FROM PACIENTES p
WHERE p.NUMERO_HC = :1`

	QueryListarDocumentos = `
SELECT
  i.PCN_NUMERO_HC,
  i.FECHA,
  i.DESCRIPCION,
  i.RUTA,
  i.TIPO,
  i.DPR_ARA_CODIGO,
  i.DPR_CODIGO,
  i.PRS_CODIGO
FROM IMAGENES_PACIENTES i
WHERE i.PCN_NUMERO_HC = :1
ORDER BY i.FECHA, i.DESCRIPCION`
)

type OpenConfig struct {
	ConnectString string
	MaxOpenConns  int
	MaxIdleConns  int
	ConnMaxLife   time.Duration
}

type Repository struct {
	db *sql.DB
}

func Open(ctx context.Context, user, password string, cfg OpenConfig) (*Repository, error) {
	if user == "" || password == "" {
		return nil, fmt.Errorf("oracle: missing credentials")
	}
	dsns, err := buildDSNCandidates(user, password, cfg.ConnectString)
	if err != nil {
		return nil, err
	}
	var lastErr error
	for _, dsn := range dsns {
		db, err := sql.Open("oracle", dsn)
		if err != nil {
			lastErr = ClassifyOracleError(err)
			continue
		}
		if cfg.MaxOpenConns > 0 {
			db.SetMaxOpenConns(cfg.MaxOpenConns)
		}
		if cfg.MaxIdleConns > 0 {
			db.SetMaxIdleConns(cfg.MaxIdleConns)
		}
		if cfg.ConnMaxLife > 0 {
			db.SetConnMaxLifetime(cfg.ConnMaxLife)
		}
		if err := pingContext(ctx, db); err != nil {
			lastErr = ClassifyOracleError(err)
			_ = db.Close()
			continue
		}
		return &Repository{db: db}, nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("oracle: unable to connect with provided descriptors")
	}
	return nil, lastErr
}

func buildDSNCandidates(user, password, connect string) ([]string, error) {
	connect = strings.TrimSpace(connect)
	if connect == "" {
		return nil, fmt.Errorf("oracle: missing connect string")
	}
	if strings.HasPrefix(connect, "oracle://") || strings.HasPrefix(connect, "user=") || strings.Contains(connect, "DESCRIPTION=") || strings.HasPrefix(connect, "(") {
		return []string{connect}, nil
	}
	host, port, service, sid, err := parseEasyConnect(connect)
	if err != nil {
		return nil, err
	}
	candidates := []string{}
	names := uniqueNames(service, sid, trimDigits(service), trimDigits(sid))
	for _, name := range names {
		if name == "" {
			continue
		}
		candidates = append(candidates, go_ora.BuildUrl(host, port, name, user, password, nil))
		candidates = append(candidates, go_ora.BuildUrl(host, port, "", user, password, map[string]string{"SID": name}))
	}
	if len(candidates) == 0 {
		candidates = append(candidates, go_ora.BuildUrl(host, port, "", user, password, nil))
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(candidates))
	for _, dsn := range candidates {
		if _, ok := seen[dsn]; ok {
			continue
		}
		seen[dsn] = struct{}{}
		out = append(out, dsn)
	}
	return out, nil
}

func uniqueNames(names ...string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(names))
	for _, n := range names {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		key := strings.ToUpper(n)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, n)
	}
	return out
}

func trimDigits(s string) string {
	s = strings.TrimSpace(s)
	end := len(s)
	for end > 0 && s[end-1] >= '0' && s[end-1] <= '9' {
		end--
	}
	if end == len(s) {
		return s
	}
	return s[:end]
}

func parseEasyConnect(connect string) (host string, port int, service string, sid string, err error) {
	if strings.Contains(connect, "/") {
		parts := strings.SplitN(connect, "/", 2)
		hostPort := strings.TrimSpace(parts[0])
		target := strings.TrimSpace(parts[1])
		h, p, perr := splitHostPortLoose(hostPort)
		if perr != nil {
			return "", 0, "", "", perr
		}
		pi, perr := strconv.Atoi(p)
		if perr != nil {
			return "", 0, "", "", fmt.Errorf("oracle: invalid port in connect string %q", connect)
		}
		if strings.EqualFold(target, "sid") || strings.HasPrefix(strings.ToUpper(target), "SID=") {
			return h, pi, "", strings.TrimPrefix(strings.TrimPrefix(target, "SID="), "sid="), nil
		}
		return h, pi, target, "", nil
	}
	if strings.Contains(connect, ":") {
		h, p, perr := splitHostPortLoose(connect)
		if perr == nil {
			pi, perr := strconv.Atoi(p)
			if perr != nil {
				return "", 0, "", "", fmt.Errorf("oracle: invalid port in connect string %q", connect)
			}
			return h, pi, "", "", nil
		}
	}
	return "", 0, "", "", fmt.Errorf("oracle: unsupported connect string %q", connect)
}

func splitHostPortLoose(s string) (string, string, error) {
	if h, p, err := net.SplitHostPort(s); err == nil {
		return h, p, nil
	}
	idx := strings.LastIndex(s, ":")
	if idx <= 0 || idx == len(s)-1 {
		return "", "", fmt.Errorf("oracle: invalid connect string %q", s)
	}
	return s[:idx], s[idx+1:], nil
}

func pingContext(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}

func (r *Repository) Close() error {
	if r == nil || r.db == nil {
		return nil
	}
	return r.db.Close()
}

func (r *Repository) BuscarPlanilla(ctx context.Context, tramite int64) ([]DigitalizacionRow, error) {
	rows, err := r.db.QueryContext(ctx, QueryBuscarPlanilla, tramite)
	if err != nil {
		return nil, ClassifyOracleError(err)
	}
	defer rows.Close()

	var out []DigitalizacionRow
	for rows.Next() {
		row, scanErr := scanDigitalizacion(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return nil, ClassifyOracleError(err)
	}
	return out, nil
}

func (r *Repository) BuscarPaciente(ctx context.Context, numeroHC int64) (*PacienteRow, error) {
	row := r.db.QueryRowContext(ctx, QueryBuscarPaciente, numeroHC)
	var p PacienteRow
	if err := row.Scan(&p.NumeroHC, &p.ApellidoPaterno, &p.ApellidoMaterno, &p.PrimerNombre, &p.SegundoNombre, &p.Cedula); err != nil {
		return nil, ClassifyOracleError(err)
	}
	return &p, nil
}

func (r *Repository) ListarDocumentos(ctx context.Context, numeroHC int64) ([]ImagenPacienteRow, error) {
	rows, err := r.db.QueryContext(ctx, QueryListarDocumentos, numeroHC)
	if err != nil {
		return nil, ClassifyOracleError(err)
	}
	defer rows.Close()

	var out []ImagenPacienteRow
	for rows.Next() {
		var d ImagenPacienteRow
		if err := rows.Scan(&d.PCNNumeroHC, &d.Fecha, &d.Descripcion, &d.Ruta, &d.Tipo, &d.DprAraCodigo, &d.DprCodigo, &d.PrsCodigo); err != nil {
			return nil, ClassifyOracleError(err)
		}
		out = append(out, d)
	}
	if err := rows.Err(); err != nil {
		return nil, ClassifyOracleError(err)
	}
	return out, nil
}

func (r *Repository) ObtenerDetallePlanilla(ctx context.Context, tramite int64) (*PlanillaDetalle, error) {
	planillas, err := r.BuscarPlanilla(ctx, tramite)
	if err != nil {
		return nil, err
	}
	if len(planillas) == 0 {
		return nil, sql.ErrNoRows
	}

	// Si hay duplicados, la UI puede mostrarlos todos.
	paciente, err := r.BuscarPaciente(ctx, planillas[0].DigHC)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	documentos, err := r.ListarDocumentos(ctx, planillas[0].DigHC)
	if err != nil {
		return nil, err
	}
	return &PlanillaDetalle{
		Planilla:   planillas[0],
		Paciente:   paciente,
		Documentos: documentos,
	}, nil
}

func scanDigitalizacion(rows *sql.Rows) (DigitalizacionRow, error) {
	var d DigitalizacionRow
	if err := rows.Scan(
		&d.DigID,
		&d.DigTramite,
		&d.DigAseguradora,
		&d.DigFechaPlanilla,
		&d.DigFechaHasta,
		&d.DigProcesado,
		&d.DigFechaProceso,
		&d.DigHC,
		&d.FePlaAniomes,
		&d.DigExpediente,
		&d.DigAnio,
		&d.DigCambioNombre,
		&d.DigFirmar,
		&d.DigAreaDep,
		&d.RevisadaContrlDoc,
		&d.DigServicio,
		&d.DigEspecialidad,
		&d.DigCedula,
		&d.DigMenorEdad,
		&d.DigPlanillado,
		&d.DigCobertura,
		&d.DigIdGeneracion,
		&d.DigIdTipo,
		&d.DigBloqueoSgh,
		&d.DigIdTramite,
		&d.DigUsuario,
		&d.DigCodigoUsuario,
		&d.DigFechaAlta,
		&d.DigNumeroPermanencia,
		&d.DigFechaCobertura,
		&d.DigPathRepo,
	); err != nil {
		return DigitalizacionRow{}, ClassifyOracleError(err)
	}
	return d, nil
}
