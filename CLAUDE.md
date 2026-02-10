# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

A zsh configuration framework ("spell book") that organizes shell aliases/functions into modular files under `src/spells/`, plus Go-based CLI utilities. A bundler (`zsb_bundle`) compiles all `.zsh` files into a single `result.zsh` that gets sourced from `.zshrc`.

## Build Commands

All make targets run from the repo root. Go workspace is in `go-work/` (Go 1.25.0).

```sh
make build                  # Build a Go command (interactive fzf selection)
make TARGET=<name> build    # Build specific command (e.g., TARGET=zsb_bundle)
make build-all              # Build all Go commands
make dev                    # Live dev mode for Go (requires entr)
make zsh-dev                # Live reload on .zsh file changes (requires entr)
```

Binaries go to `go-work/bin/`. After building zsb_bundle: `./go-work/bin/zsb_bundle` regenerates `result.zsh`.

## Testing

Go tests use standard `go test` with table-driven patterns:

```sh
cd go-work && go test ./cmd/<name>/...    # Test a specific command
cd go-work && go test ./...               # Test everything
```

## Architecture

**Two layers:**

1. **Zsh layer** (`src/`): Modular shell scripts organized by category in `src/spells/` (git, docker, files, tmux, etc.). `src/utils/` has shared zsh helpers. `src/automatic-calls/` runs on shell init. `src/globalVariables.zsh` defines shared constants (colors, git branch patterns, file type regexes).

2. **Go layer** (`go-work/`): CLI tools as separate modules under `go-work/cmd/`. Shared libraries in `go-work/lib/` (args, utils, git, open). Uses Go workspace (`go.work`) to manage multi-module setup.

**Key Go commands:** `zsb_bundle` (the bundler), `gitworktree` (git worktree TUI using Bubble Tea), `countdown`, `airbnb_calculator`, `repeatstr`, `get_repo_name`, `zsb_clipcopy`, `zsb_open`, `zsb_charm_tmux_renametab`, `zsb_charm_tmux_urlopen`.

**Bundler flow:** `zsb_bundle` walks `src/`, concatenates all `.zsh` files respecting ordering, and outputs `result.zsh`. To debug without bundling, source `debug.zsh` instead.

## Conventions

- Go modules follow one-command-per-directory pattern in `go-work/cmd/`
- Shared Go code goes in `go-work/lib/`
- New shell spells go in `src/spells/<category>/` as individual `.zsh` files
- Temporary/local shell config goes in `src/temp/` (gitignored)
- Environment config: copy `.env.example` to `.env`
- All ZSB globals/variables are prefixed with `ZSB_`
