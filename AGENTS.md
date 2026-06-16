# AGENTS.md

## Project Overview

carapace-bridge is a Go CLI tool and library that bridges shell completion from foreign completion frameworks into the [carapace](https://github.com/carapace-sh/carapace) ecosystem. It provides uniform `carapace.Action` wrappers for ~20 completion backends (cobra, argcomplete, click, yargs, bash, zsh, fish, powershell, etc.) so that any shell in the carapace family can consume completions from tools built with any of those frameworks.

## Essential Commands

```sh
go build -v ./...                          # build all packages
go test -v -coverprofile=profile.cov ./...  # run tests with coverage
gofmt -d -s .                              # check formatting (must produce no output)
staticcheck ./...                           # lint (CI installs latest before running)
```

The CI also runs `go test -v -coverprofile=profile.cov ./...` and sends coverage to Coveralls. There are **no `_test.go` files** in the repository currently — tests are effectively an integration smoke test via CI build.

## Architecture

### Go Workspace Layout

This project uses a **Go workspace** (`go.work`) with two modules:

| Module | Path | Purpose |
|--------|------|---------|
| Root | `.` (`github.com/carapace-sh/carapace-bridge`) | Library packages — `pkg/actions/bridge`, `pkg/bridges`, `pkg/choices`, `pkg/env` |
| CLI | `./cmd` (`github.com/carapace-sh/carapace-bridge/cmd`) | Binary entry point — `cmd/carapace-bridge/main.go` |

The `cmd/go.mod` uses a `replace` directive to point at the parent: `replace github.com/carapace-sh/carapace-bridge => ../`

### Two Distinct "bridge" Concepts

The codebase has two different meanings of "bridge" that are easy to confuse:

1. **`pkg/actions/bridge/`** — **Completion action bridges**: Functions like `ActionCobra`, `ActionBash`, `ActionYargs` that invoke a target command's completion mechanism and translate its output into a `carapace.Action`. These are the **public API** that consumers import.
2. **`pkg/bridges/`** — **Shell discovery**: Functions like `Bash()`, `Zsh()`, `Fish()` that enumerate which commands have completions registered in a given shell. These are used by `ActionBridge` (the meta-bridge) and `ActionBridges` (the completer for bridge names) to know what's available.

### Data Flow for Completion

1. Consumer calls `bridge.ActionCobra("kubectl", "get")` (or similar)
2. `actionCommand()` wraps it: if no command is provided, it creates an ad-hoc cobra command to prompt for one
3. The bridge action constructs env vars / CLI args specific to the target framework's completion protocol
4. `carapace.ActionExecCommand()` spawns the target command
5. Output callback parses the framework-specific format (tab-separated, colon-separated, directive integer, etc.)
6. Returns a `carapace.Action` (typically `ActionValuesDescribed` or `ActionFiles` as fallback)

### Key Patterns

- **`actionCommand()` adapter** (`carapace.go:54`): Most bridge actions use this wrapper. When `command` is empty, it creates a standalone cobra command to let the user pick a command first; when provided, it passes the command through to the inner action directly. This is the standard entry point for new bridge actions.

- **Embedded shell scripts**: Bridge implementations for bash, zsh, and fish embed shell snippets via `//go:embed` (`bash.sh`, `zsh.sh` in the capture-completion third-party package) that are written to temp files and executed in the target shell.

- **Fallback to `ActionFiles()`**: When a bridge gets empty output from the target command, most bridges fall back to `carapace.ActionFiles()`. This is deliberate — file completion is better than nothing.

- **`NoSpace([]rune("/=@:.,")...)`**: Shell bridges (bash, zsh, fish) add NoSpace for common separator characters since native completers often produce values ending with these.

### Choice System (`pkg/choices/`)

Choices are persistent per-command preferences stored as files in `$XDG_CONFIG_HOME/carapace/choices/<command>`. Format: `command/variant@group` (e.g., `kubectl/cobra@bridge`). The `ActionBridge` meta-bridge checks choices first before falling back to `CARAPACE_BRIDGES` env var.

### Environment Variables

- **`CARAPACE_BRIDGES`**: Comma-separated ordered list of implicit shell bridges to try (e.g., `bash,zsh,fish`). Used by `ActionBridge` to resolve which shell bridge to use for a given command.
- Each bridge action sets framework-specific env vars (e.g., `_ARGCOMPLETE*` for argcomplete, `COMP_LINE`/`COMP_POINT` for bash, `_<CMD>_COMPLETE` for click).

### Bridge Discovery Cache (`pkg/bridges/cache.go`)

Shell discovery (listing which commands have completions in bash/zsh/fish) is expensive, so results are cached as JSON in `$XDG_CACHE_HOME/carapace/bridges-<shell>.json` with a 24-hour TTL.

## Adding a New Bridge Action

To add a new completion framework bridge:

1. Create a new file in `pkg/actions/bridge/<name>.go`
2. Define `func Action<Name>(command ...string) carapace.Action` following the existing pattern:
   - Wrap with `actionCommand(command...)(func(command ...string) carapace.Action { ... })`
   - Inside, use `carapace.ActionCallback` → `carapace.ActionExecCommand` → parse output → return `carapace.Action`
3. Register the action in the `bridgeActions` map in `pkg/actions/bridge/bridge.go`
4. Add a subcommand in `cmd/carapace-bridge/cmd/root.go` via `addSubCommand()`
5. If the bridge is detectable (can probe a command to see if it uses this framework), add it to the `candidates` slice in `pkg/actions/bridge/detect.go`
6. If the bridge has a shell discovery counterpart (enumerate commands using it), add a function in `pkg/bridges/`

## Conventions

- **Import alias**: `shlex "github.com/carapace-sh/carapace-shlex"` — always aliased since the package name is `shlex` not `carapace-shlex`
- **Third-party code**: `third_party/` contains vendored external code (e.g., zsh-capture-completion) imported directly, NOT via go modules
- **Version injection**: `main.go` uses `var commit, date string` and `var version = "develop"` — set via ldflags at release time by GoReleaser
- **Cobra `Standalone()`**: All bridge subcommands call `carapace.Gen(cmd).Standalone()` to make them work as standalone completion providers
- **`DisableFlagParsing: true`**: Bridge subcommands disable cobra flag parsing since flags are forwarded to the bridged command

## Gotchas

- **No unit tests**: There are zero `_test.go` files. Testing is effectively done via the CI build step and Docker-based integration testing. The `detect.go` `check()` function creates temp dirs and runs actual commands, so it requires the target executables to be present.
- **Windows limitations**: Shell bridges (bash, zsh, fish) return empty results on `windows` (`runtime.GOOS == "windows"` guard in `bridges/bash.go`, `bridges/zsh.go`, etc.). `ActionArgcompleteV1` is deprecated specifically because it uses fd 8/9 which is unsupported on PowerShell/Windows.
- **Docker integration tests**: `.docker/` contains Docker Compose services that install tools for each framework (argcomplete, click, cobra, yargs, etc.) for integration testing. These are defined in `compose.yaml` and per-framework YAML files.
- **`$PATH` typo**: In `pkg/bridges/bash.go:102`, `os.Getenv("$PATH")` includes the `$` prefix — this appears to be a bug but may work in certain environments; be cautious when editing that area.
- **`urfavecli_v1` naming**: Despite the name, `urfavecli_v1` actually bridges **urfave/cli v3** (not v1), and `urfavecli` bridges v2. See the descriptions in `bridge.go`.
- **go.work is required**: Since `go.work` references `./cmd`, building the CLI module requires running from the workspace root, not from `cmd/` directly.

## Release

GoReleaser handles releases on tag push. Config in `.goreleaser.yml`:
- Binary: `carapace-bridge`
- Entry point: `./cmd/carapace-bridge`
- Targets: linux, windows, darwin (+ termux/android with CGO)
- Publishes to: Homebrew (cask), Scoop, AUR, nfpm (deb/rpm/apk), Fury.io
