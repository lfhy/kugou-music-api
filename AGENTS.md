# AGENTS

## Scope

This repository keeps tracked Go source files small and feature-grouped so regeneration, review, and manual fixes stay manageable.

## File Size Rule

- Keep every tracked `.go` file at or below 300 lines.
- Prefer staying in the 220-280 line range so later fixes do not immediately force another split.
- If a file is approaching 300 lines, split it before adding more logic.

## Comment Rule

- Every new or split source file should have a short file-level or section-level comment explaining what belongs there.
- Use comments to explain grouping and non-obvious behavior, not trivial assignments.

## Grouping Rule

- Split by feature boundary, not by arbitrary line count alone.
- Use suffix-based names when a feature grows, for example `lyric_result.go`, `lyric_live_test.go`, `signature_audio_manual.go`, `manual_pending_rank.go`.
- Keep shared helpers in adjacent helper files such as `*_helpers.go`.

## Generated SDK Files

- The source of truth for generated SDK files is `tools/gen`.
- Run generation with:
  - `go run ./tools/gen`
- The generator is expected to write sharded files directly:
  - `sdk/generated_api_*.go`
  - `sdk/generated_compat_*.go`
- If you add enough APIs that a generated shard approaches 300 lines, update the shard boundaries in:
  - `tools/gen/generate_apis.go`
  - `tools/gen/generate_compat.go`
- Do not reintroduce monolithic `sdk/generated_apis.go` or `sdk/generated_compat.go`.

## Root Alias Files

- Root-package aliases are split into:
  - `aliases_types_primary.go`
  - `aliases_types_secondary.go`
  - `aliases_routes_primary.go`
  - `aliases_routes_secondary.go`
  - `aliases_catalog.go`
- When adding new SDK request/response types or route constants that should be available from the module root, update the matching alias shard instead of creating one large alias file again.

## SDK Manual Files

- Keep manual wrappers grouped by behavior:
  - signature-based APIs in `sdk/signature_*.go`
  - pending/manual compatibility wrappers in `sdk/manual_pending_*.go`
  - collection/youth compatibility wrappers in `sdk/manual_compat_*.go`
- Preserve public method names and request/response types when splitting files.

## Examples

- `examples/subsonic_server` and `examples/subsonic_playlists` are tracked examples and must not be ignored by broad `.gitignore` patterns.
- Root-level build artifacts may be ignored, but ignore rules should be anchored, for example `/subsonic_server` instead of `subsonic_server`.

## Validation

After refactors or regeneration, run:

- `gofmt -w ./...` for touched Go files only, or `gofmt -w` on the specific paths you changed.
- `go test ./...`

## Default Platform

- Keep the SDK default platform behavior aligned with the current project expectation: default to `lite` unless there is an explicit reason to change it.
