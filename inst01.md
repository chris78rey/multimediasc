Quiero hacer un aplicativo en go que sea web super sencillo que permita loguearse con su propia clave de Oracle la idea es que se ingrese una planilla es decir un numero por ejejmplo 6069 y se busque en el esquema de digitalizacion 


Insert into DIGITALIZACION
   (DIG_ID, DIG_TRAMITE, DIG_ASEGURADORA, DIG_FECHA_PLANILLA, DIG_FECHA_HASTA, 
    DIG_PROCESADO, DIG_FECHA_PROCESO, DIG_HC, FE_PLA_ANIOMES, DIG_EXPEDIENTE, 
    DIG_ANIO, DIG_CAMBIO_NOMBRE, DIG_FIRMAR, DIG_AREA_DEP, REVISADA_CONTRL_DOC, 
    DIG_SERVICIO, DIG_ESPECIALIDAD, DIG_CEDULA, DIG_MENOR_EDAD, DIG_PLANILLADO, 
    DIG_COBERTURA, DIG_ID_GENERACION, DIG_ID_TIPO, DIG_BLOQUEO_SGH, DIG_ID_TRAMITE, 
    DIG_USUARIO, DIG_CODIGO_USUARIO, DIG_FECHA_ALTA, DIG_NUMERO_PERMANENCIA, DIG_FECHA_COBERTURA, 
    DIG_PATH_REPO)
 Values
   (237239, 6076406, 'ISSFA', TO_DATE('07/21/2026 00:00:00', 'MM/DD/YYYY HH24:MI:SS'), TO_DATE('07/21/2026 00:00:00', 'MM/DD/YYYY HH24:MI:SS'), 
    'S', TO_DATE('07/22/2026 09:24:41', 'MM/DD/YYYY HH24:MI:SS'), 308455, '202607', 'HSP07', 
    '2026', 'N', 'N', 'EMERGENCIA', 'N', 
    'EMERGENCIA', 'EMERGENCIA', '1718516196', 'N', 'S', 
    'X', 17204, 'P', 'N', 17204, 
    'VERONICA_FRUTOS', 'F006', TO_DATE('07/21/2026 00:00:00', 'MM/DD/YYYY HH24:MI:SS'), 487896, TO_DATE('07/22/2026 09:25:49', 'MM/DD/YYYY HH24:MI:SS'), 
    '/data_nuevo/repo_grande/data/datos/2026/HSP07/6076406');
COMMIT;



por ejemplo aca la digitalizacion tiene la planilla 6076406

Aca los archivos de gisgitalizacion que tiene que estan en red de windows su repositorio 

Insert into IMAGENES_PACIENTES
   (PCN_NUMERO_HC, FECHA, DESCRIPCION, RUTA, TIPO, 
    DPR_ARA_CODIGO, DPR_CODIGO, PRS_CODIGO)
 Values
   (308455, TO_DATE('07/22/2026 12:22:34', 'MM/DD/YYYY HH24:MI:SS'), 'MFE 22/07/2026 AM', '\\RED01\QUIROFANO\22072026\308455\30845522072026122234.pdf', 'DOC', 
    'H', 'I', 'M126');
Insert into IMAGENES_PACIENTES
   (PCN_NUMERO_HC, FECHA, DESCRIPCION, RUTA, TIPO, 
    DPR_ARA_CODIGO, DPR_CODIGO, PRS_CODIGO)
 Values
   (308455, TO_DATE('07/22/2026 12:21:52', 'MM/DD/YYYY HH24:MI:SS'), 'MFE POSTREANINACION - 22/07/2026', '\\RED01\QUIROFANO\22072026\308455\30845522072026122152.pdf', 'DOC', 
    'H', 'I', 'M126');
Insert into IMAGENES_PACIENTES
   (PCN_NUMERO_HC, FECHA, DESCRIPCION, RUTA, TIPO, 
    DPR_ARA_CODIGO, DPR_CODIGO, PRS_CODIGO)
 Values
   (308455, TO_DATE('07/22/2026 12:15:01', 'MM/DD/YYYY HH24:MI:SS'), 'MFE DE INGRESO 22/07/2026', '\\RED01\QUIROFANO\22072026\308455\30845522072026121501.pdf', 'DOC', 
    'H', 'I', 'M126');
Insert into IMAGENES_PACIENTES
   (PCN_NUMERO_HC, FECHA, DESCRIPCION, RUTA, TIPO, 
    DPR_ARA_CODIGO, DPR_CODIGO, PRS_CODIGO)
 Values
   (308455, TO_DATE('07/21/2026 22:14:19', 'MM/DD/YYYY HH24:MI:SS'), 'HOJA DE REFERENCIA 053', '\\RED01\QUIROFANO\21072026\308455\30845521072026221419.pdf', 'DOC', 
    'E', 'C', 'ML71');
Insert into IMAGENES_PACIENTES
   (PCN_NUMERO_HC, FECHA, DESCRIPCION, RUTA, TIPO, 
    DPR_ARA_CODIGO, DPR_CODIGO, PRS_CODIGO)
 Values
   (308455, TO_DATE('12/09/2019 07:53:32', 'MM/DD/YYYY HH24:MI:SS'), 'REGISTRO DE ANESTESIA CHULCA JOSELIN  HC. 308455 02122019', '\\RED01\QUIROFANO\09122019\308455\30845509122019075332.pdf', 'DOC', 
    'Q', 'Q', 'S015');
COMMIT;






# Oracle 11gR2
ORACLE_USER=CHRISTIAN_RUIZ
ORACLE_PASSWORD=
ORACLE_CONNECT=(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=172.16.60.21)(PORT=1521))(CONNECT_DATA=(SERVER=DEDICATED)(SID=prdsgh2)))





¡Eso lo cambia todo de una manera espectacular! Diseñar esto como una **aplicación de escritorio nativa para Windows (`.exe`)** escrita en Go es, de hecho, la decisión de arquitectura más limpia y eficiente para tu hospital. 

Al ejecutarse directamente en la sesión de Windows de cada usuario de facturación, la aplicación **hereda de forma transparente sus credenciales de red y del Active Directory**, lo que le permite acceder e inspeccionar las rutas UNC (`\\RED01\QUIROFANO\...`) de inmediato y sin necesidad de configuraciones complejas en un servidor Linux.

Para lograr este flujo de trabajo con la flexibilidad que requieres (renombrado dinámico mediante un archivo de configuración, empaquetado estructurado por carpeta de planilla, visualización del nombre del paciente y control de colisiones), el framework de desarrollo ideal es **Wails** (Go + Webview2 de Windows).

A continuación, te presento las **ideas de interfaz**, los **puntos ciegos críticos que debes prever** y el **código de backend en Go listo para producción**.

---

### 🎨 1. La Interfaz del "Renombrador Inteligente" (UX/UI de Escritorio)

La pantalla debe ser sumamente intuitiva para evitar fatiga visual en los operadores que procesan cientos de planillas al día:

```
+------------------------------------------------------------------------------------+
| 🏥 SIS - Extractor de Expedientes Digitales v3                 [👤 USER: F006] [⚙️] |
+------------------------------------------------------------------------------------+
|  Planilla/Trámite: [ 6076406 ]  [🔍 Buscar]                                        |
|  Paciente:  HC 308455 - PADILLA JAIME (ISSFA)  <-- ¡Nombre visible de inmediato!   |
+------------------------------------+-----------------------------------------------+
| 📁 ARCHIVOS DISPONIBLES EN LOTE   | 📄 VISOR PDF INTEGRADO                        |
|                                    |                                               |
| [x] Original: Epicrisis_HSP07.pdf  | +-------------------------------------------+ |
|     Guardar como: [ Epicrisis  v]  | |                                           | |
|                                    | |                                           | |
| [x] Original: 308455220720.pdf     | |               VISTA PREVIA                | |
|     Guardar como: [ 013B       v]  | |               DEL PDF                     | |
|                     |              | |             (Zero-Lag)                    | |
|                     +-- Dropdown   | |                                           | |
|                         con nombres| |                                           | |
|                         del config | +-------------------------------------------+ |
+------------------------------------+-----------------------------------------------+
|  📦 Lote de Extracción: 2 PDFs seleccionados            [⚡ Exportar ZIP de Planilla] |
+------------------------------------------------------------------------------------+
```

*   **El Dropdown de Mapeo Dinámico**: Al lado de cada archivo seleccionado en la lista, el sistema cargará un elemento de selección (`<select>`) que contendrá los nombres estándar del archivo de configuración (ej: `013B.pdf`, `Epicrisis.pdf`, `Protocolo.pdf`).
*   **Visor de PDF Integrado**: Al hacer un solo clic sobre el nombre del archivo, el panel derecho usará el control Webview2 nativo de Microsoft Edge para renderizar el PDF al instante sin descargar nada.

---

### 🧠 2. Lo que "No te estás dando cuenta" en este momento (Puntos Ciegos Críticos)

Como DBA Senior y Diseñador de Sistemas, te advierto sobre estos 4 desafíos que tu desarrollador debe blindar de inmediato:

1.  **⚠️ Colisión de Nombres en el ZIP (El Desafío de los Clones)**: 
    Si el archivo de configuración tiene el nombre estándar `"013B.pdf"` y el operador, por error o duplicación del expediente, le asigna ese mismo nombre estándar a **dos archivos originales diferentes** del mismo lote, la compresión ZIP fallará o uno sobreescribirá al otro.
    *   *Solución:* El motor en Go debe verificar si el nombre destino ya fue usado en el lote. Si se repite, debe renombrarlo automáticamente a `"013B_1.pdf"`, `"013B_2.pdf"`, o lanzar una alerta visual en la interfaz.
2.  **🔒 Permisos de Acceso UNC Temporales**:
    A veces, la red compartida `\\RED01\QUIROFANO\` se cae o la computadora de la secretaria pierde las credenciales. Si la app intenta leer un PDF y la red no responde, la aplicación puede quedarse congelada esperando el TimeOut de Windows.
    *   *Solución:* El backend de Go debe ejecutar la lectura de archivos usando hilos concurrentes (`goroutines`) con un contexto de cancelación o un timeout máximo de 3 segundos para que la interfaz nunca se cuelgue.
3.  **📂 El Archivo de Configuración Administrable**:
    El archivo `config.json` debe vivir en la misma carpeta que el `.exe`. Debe permitir que un supervisor del hospital agregue nuevos códigos de documentos (ej: `"014C.pdf"`) sin tener que volver a compilar el programa.

---

### 🛡️ 3. Código del Backend en Go (Mapeo, Renombrado y Compresión ZIP)

He estructurado de forma robusta el backend en Go que tu aplicativo de escritorio necesita para realizar la unificación relacional, cargar la configuración dinámica y exportar el ZIP estructurado:

#### A. El Archivo de Configuración (`config.json`)
Este archivo debe colocarse junto al ejecutable de la aplicación para alimentar la lista de selección de nombres estandarizados:
```json
{
  "allowed_names": [
    "013B.pdf",
    "Epicrisis.pdf",
    "Consentimiento_Informado.pdf",
    "Protocolo_Quirurgico.pdf",
    "Resultados_Laboratorio.pdf"
  ]
}
```

#### B. El Motor en Go (`main.go` con Soporte de Oracle y ZIP Estructurado)
Este módulo se conecta a la base de datos Oracle usando el driver puro de Go para extraer el nombre completo del paciente y procesar la compresión dinámica:

```go
package main

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/sijms/go-ora/v2" // Driver puro de Oracle en Go (Sin CGO)
)

// Configuración cargada desde el JSON
type AppConfig struct {
	AllowedNames []string `json:"allowed_names"`
}

// Información del Trámite/Planilla
type PlanillaResult struct {
	HC              int    `json:"hc"`
	NombrePaciente  string `json:"nombre_paciente"`
	PathRepoLinux   string `json:"path_repo_linux"`
	Aseguradora     string `json:"aseguradora"`
}

// Estructura de mapeo enviada desde el frontend al seleccionar qué descargar
type FileToRename struct {
	OriginalPath  string `json:"original_path"`  // Ruta UNC o local (ej: \\RED01\QUIROFANO\308455.pdf)
	StandardName  string `json:"standard_name"`  // Nombre mapeado por el usuario (ej: 013B.pdf)
}

type ZipRequest struct {
	PlanillaID string         `json:"planilla_id"` // Número de trámite/planilla (ej: 6076406)
	Files      []FileToRename `json:"files"`
}

type App struct {
	db     *sql.DB
	config AppConfig
}

// 1. Cargar la configuración dinámica de nombres permitidos
func (a *App) LoadConfig() error {
	exePath, _ := os.Executable()
	configPath := filepath.Join(filepath.Dir(exePath), "config.json")
	
	file, err := os.Open(configPath)
	if err != nil {
		// Valores por defecto si el archivo no existe
		a.config.AllowedNames = []string{"013B.pdf", "Epicrisis.pdf", "Protocolo.pdf"}
		return nil
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&a.config)
}

// 2. Buscar Planilla y obtener el Nombre Real del Paciente (Auditoría Relacional)
func (a *App) BuscarPlanilla(planillaID string) (*PlanillaResult, error) {
	// Query de extracción uniendo DIGITALIZACION con PACIENTES (Garantiza ver el nombre del paciente)
	query := `
		SELECT d.dig_hc,
		       TRIM(p.apellido_paterno || ' ' || NVL(p.apellido_materno, '') || ' ' || 
		            p.primer_nombre || ' ' || NVL(p.segundo_nombre, '')) AS nombre_paciente,
		       d.dig_path_repo,
		       d.dig_aseguradora
		  FROM sis.digitalizacion d
		  JOIN sis.pacientes p ON d.dig_hc = p.numero_hc
		 WHERE d.dig_tramite = :1
		   AND ROWNUM = 1`

	var res PlanillaResult
	err := a.db.QueryRow(query, planillaID).Scan(&res.HC, &res.NombrePaciente, &res.PathRepoLinux, &res.Aseguradora)
	if err != nil {
		return nil, fmt.Errorf("planilla no encontrada o inconsistencia de HC: %v", err)
	}

	return &res, nil
}

// 3. Crear ZIP Estructurado con Renombrado Dinámico y Carpeta de Trámite
func (a *App) ExportarZipContable(req ZipRequest, destZipPath string) error {
	zipFile, err := os.Create(destZipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	// Mapa para control de colisiones de nombres dentro del mismo ZIP
	nombresUsados := make(map[string]int)

	for _, fileReq := range req.Files {
		fileToCompress, err := os.Open(fileReq.OriginalPath)
		if err != nil {
			// Si un archivo de red no se puede abrir, se reporta pero no se detiene el proceso de los demás
			fmt.Printf("⚠️ No se pudo acceder al archivo %s: %v\n", fileReq.OriginalPath, err)
			continue
		}
		defer fileToCompress.Close()

		info, err := fileToCompress.Stat()
		if err != nil {
			continue
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			continue
		}

		// --- RESOLUCIÓN DE COLISIONES DE NOMBRES ---
		nombreDestino := fileReq.StandardName
		ext := filepath.Ext(nombreDestino)
		nombreSinExt := strings.TrimSuffix(nombreDestino, ext)

		if count, existe := nombresUsados[nombreDestino]; existe {
			// Si ya existe el nombre en este lote, le añadimos un sufijo seguro
			nombresUsados[nombreDestino] = count + 1
			nombreDestino = fmt.Sprintf("%s_%d%s", nombreSinExt, count+1, ext)
		} else {
			nombresUsados[nombreDestino] = 1
		}

		// --- ESTRUCTURA DEL ZIP EXIGIDA ---
		// Los archivos deben ir dentro de una carpeta con el número de planilla: "6076406/013B.pdf"
		header.Name = filepath.Join(req.PlanillaID, nombreDestino)
		header.Method = zip.Deflate

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, fileToCompress)
		if err != nil {
			return err
		}
	}

	return nil
}
```

---

🏁 Con esta arquitectura y backend nativo, la aplicación de escritorio en Windows leerá la base de datos Oracle con seguridad e inmediatez. Tus operadores visualizarán al paciente de manera clara, elegirán el mapeo correcto desde un menú desplegable ágil alimentado por tu archivo de configuración y generarán el ZIP estructurado con su carpeta contenedora en cuestión de segundos.

📊 Si lo requieres, puedo escribir un script para simular y probar de forma sintética cómo el motor en Go normaliza las barras invertidas (`\\`) de Windows a rutas UNC lógicas antes de comenzar la compilación del binario `.exe`.


se supone que los clientes ya tienen permisos y accesos al path eso ya no es mio


Necesito ademas una ventana de login que permita con el usuario de oracle loguearse cada uno
