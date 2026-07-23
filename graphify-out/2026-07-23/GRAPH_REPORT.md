# Graph Report - .  (2026-07-23)

## Corpus Check
- cluster-only mode — file stats not available

## Summary
- 565 nodes · 752 edges · 68 communities (66 shown, 2 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 92 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `2543ae47`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]
- [[_COMMUNITY_Community 17|Community 17]]
- [[_COMMUNITY_Community 19|Community 19]]
- [[_COMMUNITY_Community 20|Community 20]]
- [[_COMMUNITY_Community 21|Community 21]]
- [[_COMMUNITY_Community 22|Community 22]]
- [[_COMMUNITY_Community 23|Community 23]]
- [[_COMMUNITY_Community 24|Community 24]]
- [[_COMMUNITY_Community 25|Community 25]]
- [[_COMMUNITY_Community 26|Community 26]]
- [[_COMMUNITY_Community 27|Community 27]]
- [[_COMMUNITY_Community 28|Community 28]]
- [[_COMMUNITY_Community 29|Community 29]]
- [[_COMMUNITY_Community 30|Community 30]]
- [[_COMMUNITY_Community 31|Community 31]]
- [[_COMMUNITY_Community 32|Community 32]]
- [[_COMMUNITY_Community 33|Community 33]]
- [[_COMMUNITY_Community 34|Community 34]]
- [[_COMMUNITY_Community 35|Community 35]]
- [[_COMMUNITY_Community 36|Community 36]]
- [[_COMMUNITY_Community 37|Community 37]]
- [[_COMMUNITY_Community 38|Community 38]]
- [[_COMMUNITY_Community 39|Community 39]]
- [[_COMMUNITY_Community 40|Community 40]]
- [[_COMMUNITY_Community 41|Community 41]]
- [[_COMMUNITY_Community 42|Community 42]]
- [[_COMMUNITY_Community 43|Community 43]]
- [[_COMMUNITY_Community 44|Community 44]]
- [[_COMMUNITY_Community 45|Community 45]]
- [[_COMMUNITY_Community 46|Community 46]]
- [[_COMMUNITY_Community 47|Community 47]]
- [[_COMMUNITY_Community 48|Community 48]]
- [[_COMMUNITY_Community 49|Community 49]]
- [[_COMMUNITY_Community 50|Community 50]]
- [[_COMMUNITY_Community 51|Community 51]]
- [[_COMMUNITY_Community 52|Community 52]]
- [[_COMMUNITY_Community 53|Community 53]]
- [[_COMMUNITY_Community 54|Community 54]]
- [[_COMMUNITY_Community 55|Community 55]]
- [[_COMMUNITY_Community 56|Community 56]]
- [[_COMMUNITY_Community 57|Community 57]]
- [[_COMMUNITY_Community 58|Community 58]]
- [[_COMMUNITY_Community 59|Community 59]]
- [[_COMMUNITY_Community 60|Community 60]]
- [[_COMMUNITY_Community 61|Community 61]]
- [[_COMMUNITY_Community 62|Community 62]]
- [[_COMMUNITY_Community 63|Community 63]]
- [[_COMMUNITY_Community 64|Community 64]]
- [[_COMMUNITY_Community 65|Community 65]]
- [[_COMMUNITY_Community 66|Community 66]]
- [[_COMMUNITY_Community 67|Community 67]]
- [[_COMMUNITY_Community 68|Community 68]]
- [[_COMMUNITY_Community 69|Community 69]]
- [[_COMMUNITY_Community 70|Community 70]]
- [[_COMMUNITY_Community 71|Community 71]]

## God Nodes (most connected - your core abstractions)
1. `Files` - 75 edges
2. `Server` - 20 edges
3. `Repository` - 17 edges
4. `Requirements Document inst02.md` - 16 edges
5. `Communities (19 total, 3 thin omitted)` - 15 edges
6. `Validaciones` - 14 edges
7. `UI + ZIP Export Feature` - 14 edges
8. `desktopState` - 13 edges
9. `Requerimientos para entregar al desarrollador` - 13 edges
10. `Flujo de exportación local` - 12 edges

## Surprising Connections (you probably didn't know these)
- `QueryListarDocumentos` --references--> `IMAGENES_PACIENTES Table`  [INFERRED]
  internal/oracle/repository.go → inst01.md
- `planillaSummary()` --implements--> `Resumen de planilla`  [INFERRED]
  internal/app/desktop.go → prompt2_ui_zip.md
- `enrichDocuments()` --implements--> `4. Localización de documentos`  [INFERRED]
  internal/app/fileobs.go → inst02.md
- `inspectFile()` --implements--> `Archivos físicos`  [INFERRED]
  internal/app/fileobs.go → prompt2_ui_zip.md
- `classifyFileError()` --implements--> `Error Handling`  [INFERRED]
  internal/app/fileobs.go → inst02.md

## Import Cycles
- None detected.

## Communities (68 total, 2 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.11
Nodes (25): API endpoint: /api/export, API endpoint: /api/preview, API endpoint: /api/search, Server, createZip(), errorsIs(), exportZip(), friendlyErr() (+17 more)

### Community 1 - "Community 1"
Cohesion: 0.09
Nodes (28): Context, DB, DigitalizacionRow, Duration, PACIENTES Table, 4. Localización de documentos, Repository, Repository (+20 more)

### Community 2 - "Community 2"
Cohesion: 0.12
Nodes (20): Config, getenv(), LoadEnvConfig(), loadKeyValues(), documentStatusLine(), documentSummary(), firstNonEmpty(), formatTime() (+12 more)

### Community 3 - "Community 3"
Cohesion: 0.11
Nodes (31): Config JSON (Dynamic Renaming Rules), DIGITALIZACION Table, Go Backend (Database & ZIP Logic), Patient HC 308455 (PADILLA JAIME), IMAGENES_PACIENTES Table, Login Window (Oracle User Authentication), Oracle Database (prdsgh2), Integrated PDF Viewer (+23 more)

### Community 4 - "Community 4"
Cohesion: 0.21
Nodes (21): Contrato visual sugerido, Encabezado de sesión, Entrada de planilla, UI + ZIP Export Feature, Flujo de exportación local, Fuera de alcance, Manejo de errores, Menú de asignación de nombre (+13 more)

### Community 6 - "Community 6"
Cohesion: 0.16
Nodes (17): pageData, diagnoseWindowsPath(), escapePowerShellSingleQuotes(), fixWindowsHidden(), ImagenPacienteRow, PlanillaDetalle, windowsFileDiagnostic, Time (+9 more)

### Community 8 - "Community 8"
Cohesion: 0.03
Nodes (75): File: cmd/multimediasc/main.go, File: cmd/oraclecheck/main.go, File: go.mod, File: graphify-out/2026-07-22/graph.json, File: graphify-out/2026-07-22/GRAPH_REPORT.md, File: graphify-out/2026-07-22/.graphify_analysis.json, File: graphify-out/2026-07-22/.graphify_semantic_marker, File: graphify-out/2026-07-22/manifest.json (+67 more)

### Community 10 - "Community 10"
Cohesion: 0.15
Nodes (12): 10. Manejo de errores, 11. Auditoría, 12. Alcance que no corresponde al desarrollador, 1. Inicio de sesión, 2. Conexión con Oracle, 3. Búsqueda de la planilla, 5. Regla para relacionar documentos, 6. Vista previa (+4 more)

### Community 11 - "Community 11"
Cohesion: 0.29
Nodes (6): Env var: MULTIMEDIASC_ALLOW_DUPLICATE_FIX, Command: go run ./cmd/multimediasc, Ejecutar, MultimediaSC, Env var: ORACLE_CONNECT, Variables opcionales

### Community 14 - "Community 14"
Cohesion: 0.08
Nodes (25): Communities (19 total, 3 thin omitted), Community 0 - "Community 0", Community 10 - "Community 10", Community 11 - "Community 11", Community 14 - "Community 14", Community 15 - "Community 15", Community 16 - "Community 16", Community 1 - "Community 1" (+17 more)

### Community 15 - "Community 15"
Cohesion: 0.25
Nodes (8): 1. Encabezado de sesión, 2. Panel de búsqueda, 3. Resumen de planilla, 4. Tabla de documentos, 5. Menú de asignación de nombre, 6. Vista previa PDF, 7. Pie de exportación, Componentes UI

### Community 16 - "Community 16"
Cohesion: 0.29
Nodes (6): 🎨 1. La Interfaz del "Renombrador Inteligente" (UX/UI de Escritorio), 🧠 2. Lo que "No te estás dando cuenta" en este momento (Puntos Ciegos Críticos), 🛡️ 3. Código del Backend en Go (Mapeo, Renombrado y Compresión ZIP), A. El Archivo de Configuración (`config.json`), B. El Motor en Go (`main.go` con Soporte de Oracle y ZIP Estructurado), Oracle 11gR2

### Community 17 - "Community 17"
Cohesion: 0.20
Nodes (11): appendReason(), classifyFileError(), enrichDocuments(), fileKindFromExt(), inspectFile(), isHidden(), isHidden(), fileObservation (+3 more)

### Community 19 - "Community 19"
Cohesion: 0.15
Nodes (12): Communities (2 total, 0 thin omitted), Community 0 - "Community 0", Community 1 - "Community 1", Community Hubs (Navigation), Corpus Check, God Nodes (most connected - your core abstractions), Graph Report - .  (2026-07-22), Import Cycles (+4 more)

### Community 20 - "Community 20"
Cohesion: 0.15
Nodes (12): Communities (2 total, 0 thin omitted), Community 0 - "Community 0", Community 1 - "Community 1", Community Hubs (Navigation), Corpus Check, God Nodes (most connected - your core abstractions), Graph Report - .  (2026-07-22), Import Cycles (+4 more)

### Community 21 - "Community 21"
Cohesion: 0.29
Nodes (6): Directory Structure, File Format, File Summary, Notes, Purpose, Usage Guidelines

### Community 22 - "Community 22"
Cohesion: 0.50
Nodes (4): cmd/oraclecheck/main.go, ast_hash, mtime, semantic_hash

### Community 23 - "Community 23"
Cohesion: 0.40
Nodes (4): cmd/multimediasc/main.go, ast_hash, mtime, semantic_hash

### Community 24 - "Community 24"
Cohesion: 0.40
Nodes (4): cmd/multimediasc/main.go, ast_hash, mtime, semantic_hash

### Community 25 - "Community 25"
Cohesion: 0.50
Nodes (4): cmd/multimediasc/main.go, ast_hash, mtime, semantic_hash

### Community 26 - "Community 26"
Cohesion: 0.50
Nodes (4): go.mod, ast_hash, mtime, semantic_hash

### Community 27 - "Community 27"
Cohesion: 0.40
Nodes (4): inst01.md, ast_hash, mtime, semantic_hash

### Community 28 - "Community 28"
Cohesion: 0.50
Nodes (4): inst02.md, ast_hash, mtime, semantic_hash

### Community 29 - "Community 29"
Cohesion: 0.50
Nodes (4): internal/app/config.go, ast_hash, mtime, semantic_hash

### Community 30 - "Community 30"
Cohesion: 0.50
Nodes (4): internal/app/desktop.go, ast_hash, mtime, semantic_hash

### Community 31 - "Community 31"
Cohesion: 0.50
Nodes (4): internal/app/server.go, ast_hash, mtime, semantic_hash

### Community 32 - "Community 32"
Cohesion: 0.50
Nodes (4): internal/app/templates/index.html, ast_hash, mtime, semantic_hash

### Community 33 - "Community 33"
Cohesion: 0.50
Nodes (4): internal/app/zipwriter.go, ast_hash, mtime, semantic_hash

### Community 34 - "Community 34"
Cohesion: 0.50
Nodes (4): internal/oracle/errors.go, ast_hash, mtime, semantic_hash

### Community 35 - "Community 35"
Cohesion: 0.50
Nodes (4): internal/oracle/models.go, ast_hash, mtime, semantic_hash

### Community 36 - "Community 36"
Cohesion: 0.50
Nodes (4): internal/oracle/repository.go, ast_hash, mtime, semantic_hash

### Community 37 - "Community 37"
Cohesion: 0.50
Nodes (4): prompt2_ui_zip.md, ast_hash, mtime, semantic_hash

### Community 38 - "Community 38"
Cohesion: 0.50
Nodes (4): README.md, ast_hash, mtime, semantic_hash

### Community 39 - "Community 39"
Cohesion: 0.50
Nodes (4): cmd/oraclecheck/main.go, ast_hash, mtime, semantic_hash

### Community 40 - "Community 40"
Cohesion: 0.50
Nodes (4): go.mod, ast_hash, mtime, semantic_hash

### Community 41 - "Community 41"
Cohesion: 0.50
Nodes (4): inst01.md, ast_hash, mtime, semantic_hash

### Community 42 - "Community 42"
Cohesion: 0.50
Nodes (4): inst02.md, ast_hash, mtime, semantic_hash

### Community 43 - "Community 43"
Cohesion: 0.50
Nodes (4): internal/app/config.go, ast_hash, mtime, semantic_hash

### Community 44 - "Community 44"
Cohesion: 0.50
Nodes (4): internal/app/desktop.go, ast_hash, mtime, semantic_hash

### Community 45 - "Community 45"
Cohesion: 0.50
Nodes (4): internal/app/server.go, ast_hash, mtime, semantic_hash

### Community 46 - "Community 46"
Cohesion: 0.50
Nodes (4): internal/app/templates/index.html, ast_hash, mtime, semantic_hash

### Community 47 - "Community 47"
Cohesion: 0.50
Nodes (4): internal/app/zipwriter.go, ast_hash, mtime, semantic_hash

### Community 48 - "Community 48"
Cohesion: 0.50
Nodes (4): internal/oracle/errors.go, ast_hash, mtime, semantic_hash

### Community 49 - "Community 49"
Cohesion: 0.50
Nodes (4): internal/oracle/models.go, ast_hash, mtime, semantic_hash

### Community 50 - "Community 50"
Cohesion: 0.50
Nodes (4): internal/oracle/repository.go, ast_hash, mtime, semantic_hash

### Community 51 - "Community 51"
Cohesion: 0.50
Nodes (4): prompt2_ui_zip.md, ast_hash, mtime, semantic_hash

### Community 52 - "Community 52"
Cohesion: 0.50
Nodes (4): README.md, ast_hash, mtime, semantic_hash

### Community 54 - "Community 54"
Cohesion: 0.50
Nodes (4): cmd/oraclecheck/main.go, ast_hash, mtime, semantic_hash

### Community 55 - "Community 55"
Cohesion: 0.50
Nodes (4): go.mod, ast_hash, mtime, semantic_hash

### Community 56 - "Community 56"
Cohesion: 0.50
Nodes (4): inst01.md, ast_hash, mtime, semantic_hash

### Community 57 - "Community 57"
Cohesion: 0.50
Nodes (4): inst02.md, ast_hash, mtime, semantic_hash

### Community 58 - "Community 58"
Cohesion: 0.50
Nodes (4): internal/app/config.go, ast_hash, mtime, semantic_hash

### Community 59 - "Community 59"
Cohesion: 0.50
Nodes (4): internal/app/desktop.go, ast_hash, mtime, semantic_hash

### Community 60 - "Community 60"
Cohesion: 0.50
Nodes (4): internal/app/fileobs.go, ast_hash, mtime, semantic_hash

### Community 61 - "Community 61"
Cohesion: 0.50
Nodes (4): internal/app/fileobs_other.go, ast_hash, mtime, semantic_hash

### Community 62 - "Community 62"
Cohesion: 0.50
Nodes (4): internal/app/fileobs_windows.go, ast_hash, mtime, semantic_hash

### Community 63 - "Community 63"
Cohesion: 0.50
Nodes (4): internal/app/server.go, ast_hash, mtime, semantic_hash

### Community 64 - "Community 64"
Cohesion: 0.50
Nodes (4): internal/app/templates/index.html, ast_hash, mtime, semantic_hash

### Community 65 - "Community 65"
Cohesion: 0.50
Nodes (4): internal/app/zipwriter.go, ast_hash, mtime, semantic_hash

### Community 66 - "Community 66"
Cohesion: 0.50
Nodes (4): internal/oracle/errors.go, ast_hash, mtime, semantic_hash

### Community 67 - "Community 67"
Cohesion: 0.50
Nodes (4): internal/oracle/models.go, ast_hash, mtime, semantic_hash

### Community 68 - "Community 68"
Cohesion: 0.50
Nodes (4): internal/oracle/repository.go, ast_hash, mtime, semantic_hash

### Community 69 - "Community 69"
Cohesion: 0.50
Nodes (4): prompt2_ui_zip.md, ast_hash, mtime, semantic_hash

### Community 70 - "Community 70"
Cohesion: 0.50
Nodes (4): README.md, ast_hash, mtime, semantic_hash

### Community 71 - "Community 71"
Cohesion: 0.50
Nodes (4): scripts/unhide_share.ps1, ast_hash, mtime, semantic_hash

## Knowledge Gaps
- **317 isolated node(s):** `multimediasc`, `mtime`, `ast_hash`, `semantic_hash`, `mtime` (+312 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `enrichDocuments()` connect `Community 17` to `Community 0`, `Community 1`, `Community 2`, `Community 6`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Why does `Server` connect `Community 0` to `Community 6`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Why does `Files` connect `Community 8` to `Community 21`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `Repository` (e.g. with `Repository` and `Repository`) actually correct?**
  _`Repository` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `multimediasc`, `mtime`, `ast_hash` to the rest of the system?**
  _317 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Community 0` be split into smaller, more focused modules?**
  _Cohesion score 0.10793650793650794 - nodes in this community are weakly interconnected._
- **Should `Community 1` be split into smaller, more focused modules?**
  _Cohesion score 0.08961593172119488 - nodes in this community are weakly interconnected._