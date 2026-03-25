# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

A zsh configuration framework ("spell book") that organizes shell aliases/functions into modular files under `src/spells/`, plus Go and Rust CLI utilities. A bundler (`zsb_bundle`) compiles all `.zsh` files into a single `result.zsh` that gets sourced from `.zshrc`.

## Build Commands

All make targets run from the repo root. Go workspace is in `src-go/` (Go 1.26.0).

```sh
make go-build               # Build a Go command (interactive fzf selection)
make go-build <name>        # Build specific command (e.g., make go-build zsb_bundle)
make go-build-all           # Build all Go commands
make go-dev                 # Live dev mode for Go (requires entr)
make zsh-dev                # Live reload on .zsh file changes (requires entr)
make go-test-all            # Run all Go tests
make rust-build             # Build a Rust command (interactive fzf selection)
make rust-build <name>      # Build specific command (e.g., make rust-build wcase)
make rust-build-all         # Build all Rust commands
make rust-dev               # Live dev mode for Rust (requires entr)
make rust-test-all          # Run all Rust tests
```

Binaries go to `src-go/bin/` (Go) and `src-rust/bin/` (Rust). After building zsb_bundle: `./src-go/bin/zsb_bundle` regenerates `result.zsh`.

## Testing

Go tests use standard `go test` with table-driven patterns:

```sh
cd src-go && go test ./cmd/<name>/...    # Test a specific command
cd src-go && go test ./...               # Test everything
```

## Architecture

**Two layers:**

1. **Zsh layer** (`src/`): Modular shell scripts organized by category in `src/spells/` (git, docker, files, tmux, etc.). `src/utils/` has shared zsh helpers. `src/automatic-calls/` runs on shell init. `src/globalVariables.zsh` defines shared constants (colors, git branch patterns, file type regexes).

2. **Go layer** (`src-go/`): CLI tools as separate modules under `src-go/cmd/`. Shared libraries in `src-go/lib/` (args, utils, git, open). Uses Go workspace (`go.work`) to manage multi-module setup.

3. **Rust layer** (`src-rust/`): CLI tools as crates under `src-rust/cmd/`. Uses Cargo workspace. Shared libraries go in `src-rust/lib/`.

**Key Go commands:** `zsb_bundle` (the bundler), `gitworktree` (git worktree TUI using Bubble Tea), `countdown`, `airbnb_calculator`, `repeatstr`, `get_repo_name`, `zsb_clipcopy`, `zsb_open`.

**Key Rust commands:** `wcase` (text case transformation), `pdoro` (pomodoro daemon), `zsb_charm_tmux_renametab` (renames tmux window to current git repo name), `zsb_charm_tmux_urlopen`.

**Bundler flow:** `zsb_bundle` walks `src/`, concatenates all `.zsh` files respecting ordering, and outputs `result.zsh`. To debug without bundling, source `debug.zsh` instead.

## Conventions

- Go modules follow one-command-per-directory pattern in `src-go/cmd/`
- Shared Go code goes in `src-go/lib/`
- Rust crates follow one-command-per-directory pattern in `src-rust/cmd/`
- Shared Rust code goes in `src-rust/lib/`
- New shell spells go in `src/spells/<category>/` as individual `.zsh` files
- Temporary/local shell config goes in `src/temp/` (gitignored)
- Environment config: copy `.env.example` to `.env`
- All ZSB globals/variables are prefixed with `ZSB_`
