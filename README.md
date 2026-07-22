# MultimediaSC

App Go de consulta Oracle y exportación local ZIP.

## Ejecutar

```bash
go run ./cmd/multimediasc
```

Luego se abre una ventana de escritorio.

## Variables opcionales

- `ORACLE_CONNECT`: connect string Oracle, por defecto `172.16.60.21:1521/prdsgh2`
- `MULTIMEDIASC_ALLOW_DUPLICATE_FIX`: `true` para sufijos `_1`, `_2`
