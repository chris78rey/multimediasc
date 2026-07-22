package oracle

import "time"

// DigitalizacionRow representa una fila de DIGITALIZACION para lectura.
type DigitalizacionRow struct {
	DigID                int64      `json:"dig_id"`
	DigTramite           int64      `json:"dig_tramite"`
	DigAseguradora       string     `json:"dig_aseguradora"`
	DigFechaPlanilla     *time.Time `json:"dig_fecha_planilla,omitempty"`
	DigFechaHasta        *time.Time `json:"dig_fecha_hasta,omitempty"`
	DigProcesado         string     `json:"dig_procesado"`
	DigFechaProceso      *time.Time `json:"dig_fecha_proceso,omitempty"`
	DigHC                int64      `json:"dig_hc"`
	FePlaAniomes         string     `json:"fe_pla_aniomes"`
	DigExpediente        string     `json:"dig_expediente"`
	DigAnio              string     `json:"dig_anio"`
	DigCambioNombre      string     `json:"dig_cambio_nombre"`
	DigFirmar            string     `json:"dig_firmar"`
	DigAreaDep           string     `json:"dig_area_dep"`
	RevisadaContrlDoc    string     `json:"revisada_contrl_doc"`
	DigServicio          string     `json:"dig_servicio"`
	DigEspecialidad      string     `json:"dig_especialidad"`
	DigCedula            string     `json:"dig_cedula"`
	DigMenorEdad         string     `json:"dig_menor_edad"`
	DigPlanillado        string     `json:"dig_planillado"`
	DigCobertura         string     `json:"dig_cobertura"`
	DigIdGeneracion      int64      `json:"dig_id_generacion"`
	DigIdTipo            string     `json:"dig_id_tipo"`
	DigBloqueoSgh        string     `json:"dig_bloqueo_sgh"`
	DigIdTramite         int64      `json:"dig_id_tramite"`
	DigUsuario           string     `json:"dig_usuario"`
	DigCodigoUsuario     string     `json:"dig_codigo_usuario"`
	DigFechaAlta         *time.Time `json:"dig_fecha_alta,omitempty"`
	DigNumeroPermanencia int64      `json:"dig_numero_permanencia"`
	DigFechaCobertura    *time.Time `json:"dig_fecha_cobertura,omitempty"`
	DigPathRepo          string     `json:"dig_path_repo"`
}

// PacienteRow usa solo las columnas permitidas por el requerimiento.
type PacienteRow struct {
	NumeroHC        int64  `json:"numero_hc"`
	ApellidoPaterno string `json:"apellido_paterno"`
	ApellidoMaterno string `json:"apellido_materno"`
	PrimerNombre    string `json:"primer_nombre"`
	SegundoNombre   string `json:"segundo_nombre"`
	Cedula          string `json:"cedula"`
}

func (p PacienteRow) NombreCompleto() string {
	return joinNonEmpty(p.ApellidoPaterno, p.ApellidoMaterno, p.PrimerNombre, p.SegundoNombre)
}

type ImagenPacienteRow struct {
	PCNNumeroHC  int64      `json:"pcn_numero_hc"`
	Fecha        *time.Time `json:"fecha,omitempty"`
	Descripcion  string     `json:"descripcion"`
	Ruta         string     `json:"ruta"`
	Tipo         string     `json:"tipo"`
	DprAraCodigo string     `json:"dpr_ara_codigo"`
	DprCodigo    string     `json:"dpr_codigo"`
	PrsCodigo    string     `json:"prs_codigo"`
}

type PlanillaDetalle struct {
	Planilla   DigitalizacionRow   `json:"planilla"`
	Paciente   *PacienteRow        `json:"paciente,omitempty"`
	Documentos []ImagenPacienteRow `json:"documentos"`
}

func joinNonEmpty(parts ...string) string {
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := trimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return joinWithSpace(out)
}

func trimSpace(s string) string {
	i := 0
	j := len(s)
	for i < j && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
		i++
	}
	for j > i && (s[j-1] == ' ' || s[j-1] == '\t' || s[j-1] == '\n' || s[j-1] == '\r') {
		j--
	}
	return s[i:j]
}

func joinWithSpace(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for _, p := range parts[1:] {
		out += " " + p
	}
	return out
}
