# Graph Report - .  (2026-07-22)

## Corpus Check
- Corpus is ~2,816 words - fits in a single context window. You may not need a graph.

## Summary
- 21 nodes · 19 edges · 2 communities
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]

## God Nodes (most connected - your core abstractions)
1. `Requerimientos para entregar al desarrollador` - 13 edges
2. `Oracle 11gR2` - 4 edges
3. `🛡️ 3. Código del Backend en Go (Mapeo, Renombrado y Compresión ZIP)` - 3 edges
4. `🎨 1. La Interfaz del "Renombrador Inteligente" (UX/UI de Escritorio)` - 1 edges
5. `🧠 2. Lo que "No te estás dando cuenta" en este momento (Puntos Ciegos Críticos)` - 1 edges
6. `A. El Archivo de Configuración (`config.json`)` - 1 edges
7. `B. El Motor en Go (`main.go` con Soporte de Oracle y ZIP Estructurado)` - 1 edges
8. `1. Inicio de sesión` - 1 edges
9. `2. Conexión con Oracle` - 1 edges
10. `3. Búsqueda de la planilla` - 1 edges

## Surprising Connections (you probably didn't know these)
- None detected - all connections are within the same source files.

## Import Cycles
- None detected.

## Communities (2 total, 0 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.14
Nodes (13): 10. Manejo de errores, 11. Auditoría, 12. Alcance que no corresponde al desarrollador, 1. Inicio de sesión, 2. Conexión con Oracle, 3. Búsqueda de la planilla, 4. Localización de documentos, 5. Regla para relacionar documentos (+5 more)

### Community 1 - "Community 1"
Cohesion: 0.29
Nodes (6): 🎨 1. La Interfaz del "Renombrador Inteligente" (UX/UI de Escritorio), 🧠 2. Lo que "No te estás dando cuenta" en este momento (Puntos Ciegos Críticos), 🛡️ 3. Código del Backend en Go (Mapeo, Renombrado y Compresión ZIP), A. El Archivo de Configuración (`config.json`), B. El Motor en Go (`main.go` con Soporte de Oracle y ZIP Estructurado), Oracle 11gR2

## Knowledge Gaps
- **16 isolated node(s):** `🎨 1. La Interfaz del "Renombrador Inteligente" (UX/UI de Escritorio)`, `🧠 2. Lo que "No te estás dando cuenta" en este momento (Puntos Ciegos Críticos)`, `A. El Archivo de Configuración (`config.json`)`, `B. El Motor en Go (`main.go` con Soporte de Oracle y ZIP Estructurado)`, `1. Inicio de sesión` (+11 more)
  These have ≤1 connection - possible missing edges or undocumented components.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What connects `🎨 1. La Interfaz del "Renombrador Inteligente" (UX/UI de Escritorio)`, `🧠 2. Lo que "No te estás dando cuenta" en este momento (Puntos Ciegos Críticos)`, `A. El Archivo de Configuración (`config.json`)` to the rest of the system?**
  _16 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Community 0` be split into smaller, more focused modules?**
  _Cohesion score 0.14285714285714285 - nodes in this community are weakly interconnected._