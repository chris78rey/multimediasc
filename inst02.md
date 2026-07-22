## Requerimientos para entregar al desarrollador

Se requiere desarrollar una aplicación de escritorio para Windows, construida en Go y distribuida como archivo `.exe`. La aplicación permitirá consultar una planilla en Oracle, localizar los documentos PDF relacionados y generar un archivo ZIP con nombres estandarizados.

### 1. Inicio de sesión

El desarrollador deberá implementar:

* Una ventana de inicio de sesión con:

  * Usuario de Oracle.
  * Contraseña de Oracle.
  * Botón «Ingresar».
* La validación deberá realizarse intentando una conexión real a Oracle 11gR2.
* Cada usuario deberá conectarse con sus propias credenciales de Oracle.
* La contraseña:

  * No deberá guardarse en archivos.
  * No deberá aparecer en registros o mensajes de error.
  * Deberá permanecer únicamente en memoria mientras dure la sesión.
* Se deberá incluir una opción para cerrar sesión.
* La aplicación deberá informar claramente si:

  * El usuario o la contraseña son incorrectos.
  * La cuenta está bloqueada.
  * La contraseña está vencida.
  * La base de datos no está disponible.
  * El usuario no tiene permisos de consulta.

### 2. Conexión con Oracle

* El servidor, puerto, SID y esquema propietario deberán administrarse mediante configuración.
* No deberán colocarse usuarios ni contraseñas dentro del ejecutable.
* Todas las consultas deberán ejecutarse usando el usuario que inició sesión.
* La aplicación será exclusivamente de consulta; no deberá insertar, actualizar ni eliminar información de Oracle.
* Las consultas deberán protegerse contra inyección SQL.

### 3. Búsqueda de la planilla

La pantalla principal deberá contener:

* Campo para ingresar el número de planilla.
* Botón «Buscar».
* Validación para aceptar únicamente números válidos.
* La búsqueda deberá realizarse en `DIGITALIZACION`, utilizando `DIG_TRAMITE` como número de planilla. Por ejemplo: `6076406`.
* Si existen varios registros asociados, la aplicación deberá mostrarlos o aplicar una regla definida por el área funcional, pero no deberá escoger uno arbitrariamente.

Después de encontrar la planilla se deberán presentar, como mínimo:

* Número de planilla.
* Historia clínica.
* Nombre completo del paciente.
* Cédula.
* Aseguradora.
* Expediente.
* Fecha de la planilla.
* Servicio y especialidad.
* Número de permanencia.
* Ruta del repositorio de digitalización.

### 4. Localización de documentos

Con la historia clínica obtenida, la aplicación deberá consultar `IMAGENES_PACIENTES` y mostrar:

* Fecha del documento.
* Descripción.
* Tipo.
* Ruta UNC del archivo.
* Códigos de área, departamento y profesional, cuando correspondan.
* Estado del archivo: disponible, inaccesible, inexistente o no válido.

La lista deberá permitir:

* Seleccionar uno o varios documentos.
* Seleccionar o desmarcar todos.
* Ordenar y filtrar por fecha o descripción.
* Abrir una vista previa del PDF.
* Identificar claramente documentos antiguos o posiblemente ajenos a la planilla.

### 5. Regla para relacionar documentos

Este es el punto más importante que deberá definirse antes del desarrollo.

No es suficiente consultar `IMAGENES_PACIENTES` solamente por historia clínica, porque podrían aparecer documentos históricos que no pertenecen a la planilla actual. En el ejemplo existe incluso un documento de 2019.

El desarrollador deberá implementar una regla aprobada que considere, según corresponda:

* Historia clínica.
* Fecha de planilla.
* Fecha de alta.
* Rango de permanencia.
* Servicio o especialidad.
* Número de permanencia.
* Descripción o tipo documental.

Si todavía no existe una regla exacta, la aplicación podrá mostrar todos los documentos del paciente, pero deberá advertir cuáles se encuentran fuera del periodo de la planilla y dejar la selección bajo responsabilidad del operador.

### 6. Vista previa

* La aplicación deberá visualizar archivos PDF sin modificar el original.
* La vista previa deberá funcionar directamente desde la ruta de red.
* Un archivo inexistente, dañado o inaccesible no deberá bloquear la aplicación.
* La interfaz deberá mostrar un mensaje entendible cuando no sea posible abrirlo.

### 7. Renombrado estandarizado

Cada documento seleccionado deberá tener un menú desplegable para escoger el nombre con el que será exportado, por ejemplo:

* `013B.pdf`
* `Epicrisis.pdf`
* `Consentimiento_Informado.pdf`
* `Protocolo_Quirurgico.pdf`

La lista de nombres deberá:

* Administrarse mediante un archivo de configuración externo.
* Poder actualizarse sin recompilar la aplicación.
* Impedir nombres inválidos o caracteres no permitidos por Windows.
* Mantener la extensión `.pdf`.
* No cambiar el nombre del archivo original.

### 8. Nombres duplicados

La aplicación deberá detectar cuando dos documentos tengan asignado el mismo nombre.

Se recomienda que:

* Muestre una advertencia antes de exportar.
* Permita corregir la selección.
* Como alternativa configurable, agregue sufijos como `_1`, `_2`, etc.
* Nunca sobrescriba silenciosamente un documento.

### 9. Generación del ZIP

La aplicación deberá:

* Permitir seleccionar dónde guardar el archivo.
* Generar un ZIP identificado con el número de planilla.
* Crear dentro del ZIP una carpeta con ese mismo número.
* Incluir solamente los documentos seleccionados.
* Incorporarlos con los nombres estandarizados.
* Mostrar el progreso de la operación.
* Permitir cancelar el proceso.
* Informar qué documentos fueron incluidos y cuáles presentaron errores.
* Evitar presentar como exitoso un ZIP incompleto sin advertencia.
* No modificar, mover ni eliminar los documentos originales.

Ejemplo de estructura esperada:

```text
6076406.zip
└── 6076406
    ├── 013B.pdf
    ├── Epicrisis.pdf
    └── Protocolo_Quirurgico.pdf
```

### 10. Manejo de errores

Deberán controlarse, como mínimo:

* Planilla inexistente.
* Planilla duplicada.
* Paciente no encontrado.
* Planilla sin documentos.
* Pérdida de conexión con Oracle.
* Ruta de red no disponible.
* Archivo eliminado, dañado o sin acceso.
* Nombre de exportación duplicado.
* Falta de espacio en disco.
* Error durante la creación del ZIP.

Las operaciones sobre rutas de red deberán tener tiempo de espera y no podrán congelar la interfaz.

### 11. Auditoría

Se deberá registrar:

* Usuario Oracle que ejecutó la operación.
* Fecha y hora.
* Número de planilla.
* Documentos seleccionados.
* Nombres asignados.
* Resultado de la exportación.
* Errores producidos.

Los registros no deberán contener contraseñas ni información clínica innecesaria.

### 12. Alcance que no corresponde al desarrollador

La aplicación no deberá encargarse de:

* Crear usuarios de Oracle.
* Cambiar contraseñas.
* Otorgar permisos de base de datos.
* Crear accesos a carpetas compartidas.
* Administrar permisos de Active Directory o Windows.
* Corregir rutas incorrectas registradas en Oracle.
* Modificar documentos originales.
* Actualizar información de `DIGITALIZACION` o `IMAGENES_PACIENTES`.

Los equipos clientes ya cuentan con permisos sobre las rutas UNC. La aplicación solamente deberá detectar y reportar cuando un archivo no pueda abrirse.

Antes de iniciar el desarrollo deberán quedar confirmados dos aspectos: la tabla y los campos exactos para obtener el nombre del paciente, y la regla funcional que determina cuáles documentos de `IMAGENES_PACIENTES` pertenecen realmente a cada planilla.
