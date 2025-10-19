# Repository Guidelines

## Project Structure & Module Organization
3X-UI is a Go module with `main.go` bootstrapping the HTTP control panel. Runtime logic resides in `web/`, which contains controllers, middleware, sessions, and localized strings; static assets sit in `web/assets` and templates in `web/html`. Xray integration lives in `xray/`, subscription APIs in `sub/`, and background jobs in `web/service`. Configuration helpers are under `config/`, persistent stores in `database/`, and shared helpers in `util/` and `logger/`. Container and install entrypoints (`Dockerfile`, `docker-compose.yml`, `install.sh`) mirror production defaults—keep them updated when ports, volumes, or binaries change.

## Build, Test, and Development Commands
- `go build ./...` compiles every package and surfaces compile-time regressions.
- `go run ./main.go` launches the panel using local config within `config/` and the bundled SQLite artifacts under `database/`.
- `go test ./...` executes the Go suite; run it before each PR even though current coverage is small.
- `docker-compose up --build` reproduces the containerized deployment; update `docker-compose.yml` when listener ports, volumes, or service names move.

## Coding Style & Naming Conventions
Format Go code with `gofmt` (or `goimports`) before committing—CI expects zero diff. Follow Go naming: exported symbols in `UpperCamelCase`, internal helpers in `lowerCamelCase`, and filenames lowercase with optional underscores (`sub/subService.go`). Keep handlers cohesive inside packages instead of monolithic files, and prefer structured logging via the helpers in `logger/`.

## Testing Guidelines
Place tests beside implementation files using the `_test.go` suffix and table-driven cases. Favor integration tests in `web/service` to validate routing behavior, mocking Xray dependencies through abstractions in `xray/`. When adding migrations or models, add coverage under `database/` to confirm schema assumptions. Target meaningful coverage gains on new work; document gaps in the PR when tests are impractical.

## Commit & Pull Request Guidelines
Write commit subjects in the imperative mood (e.g., “Guard xray template handling”), mirroring existing history. Keep changes focused; separate refactors from features. Pull requests should include a concise summary, linked issues, testing notes (`go test ./...`, docker run results), and UI evidence when screens change. Highlight configuration or migration impacts so reviewers can verify upgrade paths.

## Security & Configuration Tips
Avoid committing secrets or customer data; rely on environment variables or Docker secrets. Update example configs under `config/` when defaults shift, and document new ports or permissions. If you touch install scripts, validate both `install.sh` and `DockerEntrypoint.sh` so operators avoid regressions during upgrades.
