# Graph Report - multimediasc  (2026-07-24)

## Corpus Check
- 94 files · ~103,603 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 869 nodes · 1262 edges · 88 communities (83 shown, 5 thin omitted)
- Extraction: 88% EXTRACTED · 12% INFERRED · 0% AMBIGUOUS · INFERRED: 150 edges (avg confidence: 0.82)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0e628dd1`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]
- [[_COMMUNITY_Community 17|Community 17]]
- [[_COMMUNITY_Community 18|Community 18]]
- [[_COMMUNITY_Community 19|Community 19]]
- [[_COMMUNITY_Community 20|Community 20]]
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
- [[_COMMUNITY_Community 72|Community 72]]
- [[_COMMUNITY_Community 73|Community 73]]
- [[_COMMUNITY_Community 74|Community 74]]
- [[_COMMUNITY_Community 75|Community 75]]
- [[_COMMUNITY_Community 76|Community 76]]
- [[_COMMUNITY_Community 77|Community 77]]
- [[_COMMUNITY_Community 78|Community 78]]
- [[_COMMUNITY_Community 79|Community 79]]
- [[_COMMUNITY_Community 82|Community 82]]
- [[_COMMUNITY_Community 83|Community 83]]
- [[_COMMUNITY_Community 84|Community 84]]
- [[_COMMUNITY_Community 85|Community 85]]
- [[_COMMUNITY_Community 86|Community 86]]
- [[_COMMUNITY_Community 87|Community 87]]
- [[_COMMUNITY_Community 88|Community 88]]

## God Nodes (most connected - your core abstractions)
1. `trimSpace()` - 43 edges
2. `BatchCore` - 40 edges
3. `App` - 31 edges
4. `desktopState` - 30 edges
5. `BatchSnapshot` - 23 edges
6. `$()` - 20 edges
7. `RunDesktop()` - 20 edges
8. `Repository` - 18 edges
9. `Open()` - 16 edges
10. `Requirements Document inst02.md` - 16 edges

## Surprising Connections (you probably didn't know these)
- `QueryListarDocumentos` --references--> `IMAGENES_PACIENTES Table`  [INFERRED]
  internal/oracle/repository.go → inst01.md
- `planillaSummary()` --implements--> `Resumen de planilla`  [INFERRED]
  internal/app/desktop.go → prompt2_ui_zip.md
- `inspectFile()` --implements--> `Archivos físicos`  [INFERRED]
  internal/app/fileobs.go → prompt2_ui_zip.md
- `classifyFileError()` --implements--> `Error Handling`  [INFERRED]
  internal/app/fileobs.go → inst02.md
- `ClassifyOracleError()` --implements--> `Error Handling`  [INFERRED]
  internal/oracle/errors.go → inst02.md

## Import Cycles
- None detected.

## Communities (88 total, 5 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.33
Nodes (4): API endpoint: /api/export, API endpoint: /api/preview, API endpoint: /api/search, index.html template

### Community 1 - "Community 1"
Cohesion: 0.10
Nodes (34): parsePlanillaInputs(), Config JSON (Dynamic Renaming Rules), DIGITALIZACION Table, Go Backend (Database & ZIP Logic), Patient HC 308455 (PADILLA JAIME), IMAGENES_PACIENTES Table, Login Window (Oracle User Authentication), Oracle Database (prdsgh2) (+26 more)

### Community 2 - "Community 2"
Cohesion: 0.11
Nodes (24): documentStatusLine(), documentSummary(), ensureZipExtension(), firstNonEmpty(), formatTime(), planillaPatientName(), planillaSummary(), RunDesktop() (+16 more)

### Community 4 - "Community 4"
Cohesion: 0.12
Nodes (30): 1. Encabezado de sesión, 2. Panel de búsqueda, 3. Resumen de planilla, 4. Tabla de documentos, 5. Menú de asignación de nombre, 6. Vista previa PDF, 7. Pie de exportación, Componentes UI (+22 more)

### Community 5 - "Community 5"
Cohesion: 0.38
Nodes (5): ClearHiddenAttributes(), Duration, fail(), filepathAbs(), main()

### Community 6 - "Community 6"
Cohesion: 0.25
Nodes (7): Current runtime references, Keep as backend boundaries, Migration phases, New Wails-facing structure, Replace the UI boundary, Risks to watch, Wails Migration Plan

### Community 9 - "Community 9"
Cohesion: 0.11
Nodes (16): Config, getenv(), LoadEnvConfig(), loadKeyValues(), normalizeAllowedNames(), parseNamesText(), parseRangesText(), friendlyErr() (+8 more)

### Community 10 - "Community 10"
Cohesion: 0.22
Nodes (12): basenameAnySeparator(), countSelected(), createZip(), exportZip(), normalizeExportName(), validWindowsFilename(), zipDocument, zipWriter (+4 more)

### Community 11 - "Community 11"
Cohesion: 0.25
Nodes (7): Env var: MULTIMEDIASC_ALLOW_DUPLICATE_FIX, Command: go run ./cmd/multimediasc, About, Building, Live Development, README, Env var: ORACLE_CONNECT

### Community 12 - "Community 12"
Cohesion: 0.06
Nodes (23): BatchCore, BatchSnapshot, DefaultAutoOutDir(), DefaultZipPath(), documentSummary(), firstNonEmpty(), formatTime(), NewBatchCore() (+15 more)

### Community 14 - "Community 14"
Cohesion: 0.08
Nodes (25): Communities (19 total, 3 thin omitted), Community 0 - "Community 0", Community 10 - "Community 10", Community 11 - "Community 11", Community 14 - "Community 14", Community 15 - "Community 15", Community 16 - "Community 16", Community 1 - "Community 1" (+17 more)

### Community 15 - "Community 15"
Cohesion: 0.05
Nodes (3): EventsOn(), EventsOnce(), EventsOnMultiple()

### Community 16 - "Community 16"
Cohesion: 0.29
Nodes (6): 🎨 1. La Interfaz del "Renombrador Inteligente" (UX/UI de Escritorio), 🧠 2. Lo que "No te estás dando cuenta" en este momento (Puntos Ciegos Críticos), 🛡️ 3. Código del Backend en Go (Mapeo, Renombrado y Compresión ZIP), A. El Archivo de Configuración (`config.json`), B. El Motor en Go (`main.go` con Soporte de Oracle y ZIP Estructurado), Oracle 11gR2

### Community 17 - "Community 17"
Cohesion: 0.06
Nodes (29): appendReason(), classifyFileError(), enrichDocuments(), fileKindFromExt(), inspectFile(), isHidden(), isHidden(), fileObservation (+21 more)

### Community 19 - "Community 19"
Cohesion: 0.15
Nodes (12): Communities (2 total, 0 thin omitted), Community 0 - "Community 0", Community 1 - "Community 1", Community Hubs (Navigation), Corpus Check, God Nodes (most connected - your core abstractions), Graph Report - .  (2026-07-22), Import Cycles (+4 more)

### Community 20 - "Community 20"
Cohesion: 0.15
Nodes (12): Communities (2 total, 0 thin omitted), Community 0 - "Community 0", Community 1 - "Community 1", Community Hubs (Navigation), Corpus Check, God Nodes (most connected - your core abstractions), Graph Report - .  (2026-07-22), Import Cycles (+4 more)

### Community 22 - "Community 22"
Cohesion: 0.40
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
Cohesion: 0.50
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

### Community 72 - "Community 72"
Cohesion: 0.18
Nodes (10): frontend:build, frontend:dev:serverUrl, frontend:dev:watcher, frontend:install, author, email, name, name (+2 more)

### Community 73 - "Community 73"
Cohesion: 0.20
Nodes (9): devDependencies, vite, name, private, scripts, build, dev, preview (+1 more)

### Community 74 - "Community 74"
Cohesion: 0.19
Nodes (15): $(), applyTheme(), call(), canonicalName(), escapeHtml(), pickVisibleDoc(), render(), renderNameSelect() (+7 more)

### Community 75 - "Community 75"
Cohesion: 0.25
Nodes (7): author, email, name, name, outputfilename, $schema, wailsjsdir

### Community 76 - "Community 76"
Cohesion: 0.40
Nodes (3): Greet(), nameElement, resultElement

### Community 78 - "Community 78"
Cohesion: 0.40
Nodes (4): About, Building, Live Development, README

### Community 79 - "Community 79"
Cohesion: 0.12
Nodes (15): author, bugs, url, description, homepage, keywords, license, main (+7 more)

### Community 82 - "Community 82"
Cohesion: 0.12
Nodes (15): author, bugs, url, description, homepage, keywords, license, main (+7 more)

### Community 84 - "Community 84"
Cohesion: 0.08
Nodes (34): mustJSON(), DB, DigitalizacionRow, Time, Context, ImagenPacienteRow, PlanillaDetalle, NullString (+26 more)

### Community 85 - "Community 85"
Cohesion: 0.18
Nodes (3): BatchSnapshot, DocumentoView, PlanillaView

### Community 86 - "Community 86"
Cohesion: 0.40
Nodes (4): EnvironmentInfo, Position, Screen, Size

### Community 87 - "Community 87"
Cohesion: 0.40
Nodes (4): EnvironmentInfo, Position, Screen, Size

### Community 88 - "Community 88"
Cohesion: 0.67
Nodes (3): EventsOn(), EventsOnce(), EventsOnMultiple()

## Knowledge Gaps
- **316 isolated node(s):** `name`, `private`, `version`, `dev`, `build` (+311 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `trimSpace()` connect `Community 9` to `Community 1`, `Community 2`, `Community 10`, `Community 12`, `Community 17`, `Community 84`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **Why does `Open()` connect `Community 84` to `Community 1`, `Community 9`, `Community 10`, `Community 12`, `Community 17`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Why does `BatchCore` connect `Community 12` to `Community 2`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Are the 41 inferred relationships involving `trimSpace()` (e.g. with `.ExportNameForIndex()` and `.FilteredPlanillaIndexes()`) actually correct?**
  _`trimSpace()` has 41 INFERRED edges - model-reasoned connections that need verification._
- **What connects `name`, `private`, `version` to the rest of the system?**
  _316 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Community 1` be split into smaller, more focused modules?**
  _Cohesion score 0.09915966386554621 - nodes in this community are weakly interconnected._
- **Should `Community 2` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._