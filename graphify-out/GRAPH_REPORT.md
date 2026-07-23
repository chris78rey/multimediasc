# Graph Report - multimediasc  (2026-07-23)

## Corpus Check
- 18 files · ~8,662 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 218 nodes · 363 edges · 19 communities (16 shown, 3 thin omitted)
- Extraction: 83% EXTRACTED · 17% INFERRED · 0% AMBIGUOUS · INFERRED: 62 edges (avg confidence: 0.79)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a0fe8e9d`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]
- [[_COMMUNITY_Community 17|Community 17]]

## God Nodes (most connected - your core abstractions)
1. `Server` - 22 edges
2. `Repository` - 17 edges
3. `Requirements Document inst02.md` - 16 edges
4. `Open()` - 15 edges
5. `UI + ZIP Export Feature` - 14 edges
6. `trimSpace()` - 13 edges
7. `Requerimientos para entregar al desarrollador` - 13 edges
8. `Validaciones` - 13 edges
9. `desktopState` - 12 edges
10. `Prompt 2: UI + ZIP` - 11 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `RunDesktop()`  [INFERRED]
  cmd/multimediasc/main.go → internal/app/desktop.go
- `main()` --calls--> `getenv()`  [INFERRED]
  cmd/oraclecheck/main.go → internal/app/config.go
- `main()` --calls--> `Open()`  [INFERRED]
  cmd/oraclecheck/main.go → internal/oracle/repository.go
- `firstNonEmpty()` --calls--> `trimSpace()`  [INFERRED]
  cmd/oraclecheck/main.go → internal/oracle/models.go
- `loadEnvFile()` --calls--> `trimSpace()`  [INFERRED]
  cmd/oraclecheck/main.go → internal/oracle/models.go

## Import Cycles
- None detected.

## Communities (19 total, 3 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.14
Nodes (19): pageData, Server, createZip(), errorsIs(), exportZip(), friendlyErr(), NewServer(), validatePlanilla() (+11 more)

### Community 1 - "Community 1"
Cohesion: 0.11
Nodes (29): Config, getenv(), LoadEnvConfig(), loadKeyValues(), isHidden(), Context, DB, DigitalizacionRow (+21 more)

### Community 2 - "Community 2"
Cohesion: 0.15
Nodes (16): documentStatusLine(), documentSummary(), firstNonEmpty(), formatTime(), planillaSummary(), RunDesktop(), desktopState, Container (+8 more)

### Community 3 - "Community 3"
Cohesion: 0.24
Nodes (17): Desktop Application, Audit Log, Configuration File, DIGITALIZACION Table, Requirements Document inst02.md, Duplicate Name Avoidance, Error Handling, IMAGENES_PACIENTES Table (+9 more)

### Community 4 - "Community 4"
Cohesion: 0.18
Nodes (23): Archivos físicos, Contrato visual sugerido, Encabezado de sesión, Entrada de planilla, UI + ZIP Export Feature, Flujo de exportación local, Fuera de alcance, Manejo de errores (+15 more)

### Community 5 - "Community 5"
Cohesion: 0.20
Nodes (14): Config JSON (Dynamic Renaming Rules), DIGITALIZACION Table, Go Backend (Database & ZIP Logic), Patient HC 308455 (PADILLA JAIME), IMAGENES_PACIENTES Table, Login Window (Oracle User Authentication), Oracle Database (prdsgh2), PACIENTES Table (+6 more)

### Community 6 - "Community 6"
Cohesion: 0.39
Nodes (7): Time, DigitalizacionRow, ImagenPacienteRow, joinNonEmpty(), joinWithSpace(), PacienteRow, PlanillaDetalle

### Community 8 - "Community 8"
Cohesion: 0.47
Nodes (4): zipWriter, newZipWriter(), File, Writer

### Community 9 - "Community 9"
Cohesion: 0.40
Nodes (4): API endpoint: /api/export, API endpoint: /api/preview, API endpoint: /api/search, index.html template

### Community 10 - "Community 10"
Cohesion: 0.14
Nodes (13): 10. Manejo de errores, 11. Auditoría, 12. Alcance que no corresponde al desarrollador, 1. Inicio de sesión, 2. Conexión con Oracle, 3. Búsqueda de la planilla, 4. Localización de documentos, 5. Regla para relacionar documentos (+5 more)

### Community 11 - "Community 11"
Cohesion: 0.29
Nodes (6): Env var: MULTIMEDIASC_ALLOW_DUPLICATE_FIX, Command: go run ./cmd/multimediasc, Ejecutar, MultimediaSC, Env var: ORACLE_CONNECT, Variables opcionales

### Community 14 - "Community 14"
Cohesion: 0.43
Nodes (7): appendReason(), classifyFileError(), enrichDocuments(), fileKindFromExt(), inspectFile(), fileObservation, ImagenPacienteRow

### Community 15 - "Community 15"
Cohesion: 0.25
Nodes (8): 1. Encabezado de sesión, 2. Panel de búsqueda, 3. Resumen de planilla, 4. Tabla de documentos, 5. Menú de asignación de nombre, 6. Vista previa PDF, 7. Pie de exportación, Componentes UI

### Community 16 - "Community 16"
Cohesion: 0.29
Nodes (6): 🎨 1. La Interfaz del "Renombrador Inteligente" (UX/UI de Escritorio), 🧠 2. Lo que "No te estás dando cuenta" en este momento (Puntos Ciegos Críticos), 🛡️ 3. Código del Backend en Go (Mapeo, Renombrado y Compresión ZIP), A. El Archivo de Configuración (`config.json`), B. El Motor en Go (`main.go` con Soporte de Oracle y ZIP Estructurado), Oracle 11gR2

## Knowledge Gaps
- **61 isolated node(s):** `multimediasc`, `Config`, `Repository`, `Label`, `Container` (+56 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Open()` connect `Community 1` to `Community 0`, `Community 8`, `Community 2`, `Community 14`?**
  _High betweenness centrality (0.106) - this node is a cross-community bridge._
- **Why does `trimSpace()` connect `Community 1` to `Community 0`, `Community 2`, `Community 6`?**
  _High betweenness centrality (0.077) - this node is a cross-community bridge._
- **Are the 9 inferred relationships involving `Open()` (e.g. with `loadKeyValues()` and `.search()`) actually correct?**
  _`Open()` has 9 INFERRED edges - model-reasoned connections that need verification._
- **What connects `multimediasc`, `Config`, `Repository` to the rest of the system?**
  _61 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Community 0` be split into smaller, more focused modules?**
  _Cohesion score 0.1396011396011396 - nodes in this community are weakly interconnected._
- **Should `Community 1` be split into smaller, more focused modules?**
  _Cohesion score 0.11411411411411411 - nodes in this community are weakly interconnected._
- **Should `Community 10` be split into smaller, more focused modules?**
  _Cohesion score 0.14285714285714285 - nodes in this community are weakly interconnected._