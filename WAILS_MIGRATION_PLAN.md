# Wails Migration Plan

Goal: replace the current Fyne desktop UI with a Wails-based desktop frontend while preserving the Go backend, Oracle read-only behavior, ZIP export, and file observability.

## Keep as backend boundaries

- `internal/oracle/repository.go`
- `internal/oracle/models.go`
- `internal/oracle/errors.go`
- `internal/app/export.go`
- `internal/app/fileobs.go`
- `internal/app/fileobs_windows.go`
- `internal/app/fileobs_other.go`
- `internal/app/config.go`

These files contain the business rules and should remain the source of truth.

## Replace the UI boundary

- `internal/app/desktop.go`

This file currently holds the Fyne UI. It should be retired or reduced to a thin bootstrap while Wails becomes the desktop frontend.

## New Wails-facing structure

Suggested shape:

- `cmd/multimediasc/main.go`
- `internal/app/services/*.go`
- `frontend/`

Responsibilities:

- `main.go`: start Wails app
- `services`: expose Oracle lookup, batch planilla review, ZIP export, and file checks to the frontend
- `frontend`: render the batch review flow in HTML/CSS/JS

## Migration phases

1. Extract any remaining UI-independent logic from `desktop.go` into backend services.
2. Preserve Oracle-only-read access exactly as it is.
3. Keep ZIP naming, folder grouping, and duplicate-name handling unchanged.
4. Build the Wails frontend for:
   - batch planilla search
   - document review
   - document selection
   - output folder selection
   - ZIP confirmation
5. Rebuild the Windows executable and validate the desktop flow.

## Risks to watch

- Do not expose physical paths in the frontend.
- Do not let the frontend bypass Oracle validation.
- Keep document IDs or safe references between UI and backend.
- Preserve the current ZIP contract: one folder per planilla, documents inside by planilla.

## Current runtime references

- Current desktop entrypoint: `cmd/multimediasc/main.go`
- Current desktop UI: `internal/app/desktop.go`
- Current export behavior: `internal/app/export.go`
- Current Oracle data flow: `internal/oracle/repository.go`
