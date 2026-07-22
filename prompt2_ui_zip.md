# Prompt 2: UI + ZIP

Alcance de esta fase:

- Solo UI de escritorio.
- Solo exportación local a ZIP.
- No modificar consultas Oracle.
- No crear nueva lógica de acceso a datos.
- Consumir únicamente DTOs ya expuestos por el backend.

## Objetivo funcional

La pantalla debe permitir:

- Ingresar una planilla.
- Ver una tabla de documentos.
- Seleccionar uno o varios documentos.
- Asignar un nombre estandarizado por documento.
- Previsualizar PDF.
- Exportar un ZIP con carpeta interna de la planilla.
- Detectar nombres duplicados.
- Reportar archivos inexistentes o inaccesibles.

## Componentes UI

### 1. Encabezado de sesión

- Usuario activo.
- Estado de conexión.
- Botón de cerrar sesión si ya existe en el shell actual.

### 2. Panel de búsqueda

- Campo `planilla`.
- Botón `Buscar`.
- Mensaje de validación de entrada.
- Estado de búsqueda:
  - `idle`
  - `buscando`
  - `sin_resultados`
  - `resultado_cargado`
  - `error`

### 3. Resumen de planilla

Debe mostrar los DTOs ya disponibles del backend, sin inferir nuevos datos:

- Número de planilla.
- Historia clínica.
- Nombre del paciente.
- Cédula.
- Aseguradora.
- Expediente.
- Fecha de la planilla.
- Servicio.
- Especialidad.
- Número de permanencia.
- Ruta base del repositorio.

### 4. Tabla de documentos

Columnas mínimas:

- Selección.
- Fecha.
- Descripción.
- Tipo.
- Ruta UNC.
- Estado del archivo.
- Nombre estandarizado.
- Acción de vista previa.

Comportamientos:

- Selección múltiple por filas.
- Seleccionar todos.
- Desmarcar todos.
- Ordenar por fecha o descripción.
- Filtrar por texto.
- Resaltar documentos fuera de periodo si el backend ya marca esa condición en el DTO.

### 5. Menú de asignación de nombre

Cada fila seleccionable debe tener un selector con nombres estandarizados.

Reglas:

- La lista viene de configuración externa ya cargada por el backend.
- No se permite editar libremente si el backend ya entrega catálogo cerrado.
- Mantener `.pdf`.
- No permitir caracteres inválidos para Windows.

### 6. Vista previa PDF

- Panel lateral o modal con preview.
- Carga directa desde la ruta del documento.
- No modifica el original.
- Si el archivo no abre, mostrar mensaje explícito en la propia UI.

### 7. Pie de exportación

- Conteo de documentos seleccionados.
- Total de errores detectados.
- Botón `Exportar ZIP`.
- Botón `Cancelar`.
- Barra de progreso.
- Mensaje de estado final.

## Flujo de exportación local

1. El usuario busca una planilla.
2. La UI carga el DTO de planilla y el listado de documentos.
3. El usuario selecciona documentos.
4. La UI asigna un nombre estandarizado a cada fila seleccionada.
5. La UI valida duplicados antes de exportar.
6. La UI solicita destino de guardado.
7. El motor local crea el ZIP con esta estructura:

```text
6076406.zip
└── 6076406
    ├── 013B.pdf
    ├── Epicrisis.pdf
    └── Protocolo_Quirurgico.pdf
```

8. El proceso reporta:

- Incluidos.
- Omitidos.
- Inaccesibles.
- Inexistentes.
- Duplicados.

## Validaciones

### Entrada de planilla

- Solo números.
- No vacía.
- Longitud mínima razonable según el formato existente.

### Selección de documentos

- Debe haber al menos un documento seleccionado para exportar.
- No se puede exportar una fila sin nombre asignado.

### Nombres estandarizados

- No vacíos.
- Sin caracteres inválidos para Windows.
- Con extensión `.pdf`.
- Sin duplicados silenciosos.

### Archivos físicos

- Verificar existencia antes de previsualizar.
- Verificar accesibilidad antes de incluir en ZIP.
- No bloquear la UI si la ruta UNC no responde.

### ZIP

- Validar espacio suficiente si el backend expone esta señal.
- No sobrescribir archivos existentes sin confirmación explícita del usuario.

## Manejo de nombres duplicados

Prioridad recomendada:

1. Bloquear exportación y mostrar advertencia.
2. Permitir corrección manual.
3. Solo si el backend trae una política explícita, aplicar sufijos `_1`, `_2`, etc.

No se debe renombrar silenciosamente.

## Manejo de errores

La UI debe diferenciar estos casos:

- Planilla inexistente.
- Planilla duplicada.
- Paciente no encontrado.
- Planilla sin documentos.
- Ruta de red no disponible.
- Archivo eliminado.
- Archivo dañado.
- Archivo sin permisos.
- Nombre duplicado.
- Error de ZIP.
- Falta de espacio.

Mensajes:

- Breves.
- Entendibles.
- Sin exponer credenciales.
- Sin mostrar detalles técnicos innecesarios al operador.

## Progreso

La exportación debe exponer progreso por etapa:

- Preparando.
- Validando selección.
- Abriendo archivos.
- Escribiendo ZIP.
- Cerrando archivo.
- Finalizando.

Si el backend permite avances por documento, la UI debe reflejarlos por fila y en la barra global.

## Contrato visual sugerido

Distribución:

- Cabecera superior.
- Columna izquierda:
  - búsqueda
  - resumen de planilla
  - tabla de documentos
- Columna derecha:
  - vista previa PDF
  - panel de errores y advertencias
- Franja inferior:
  - progreso
  - exportación

## Resultado esperado de esta fase

El entregable de UI + ZIP debe incluir:

- Componentes de pantalla.
- Flujo de exportación.
- Validaciones.
- Manejo de progreso.
- Manejo de errores.

## Fuera de alcance

- Consultas Oracle nuevas.
- Cambios en SQL.
- Cambios en acceso a datos.
- Escritura o modificación de archivos originales.
- Modificación de reglas clínicas o de negocio de backend.
