# Git Workflow for plugin-virtual

This repository contains the Slidebolt Virtual Plugin, which provides purely virtual entities for testing and workspace simulation. It produces a standalone binary.

## Dependencies
- **Internal:**
  - `sb-contract`: Core interfaces.
  - `sb-domain`: Shared domain models.
  - `sb-messenger-sdk`: Shared messaging interfaces.
  - `sb-runtime`: Core execution environment.
  - `sb-storage-sdk`: Shared storage interfaces.
  - `sb-testkit`: Testing utilities.
- **External:** 
  - Standard Go library and NATS.
  - `github.com/cucumber/godog`: BDD testing framework.

## Build Process
- **Type:** Go Application (Plugin).
- **Consumption:** Run as a background plugin service for simulation.
- **Artifacts:** Produces a binary named `plugin-virtual`.
- **Command:** `go build -o plugin-virtual ./cmd/plugin-virtual`
- **Validation:** 
  - Validated through unit tests: `go test -v ./...`
  - Validated through BDD tests: `go test -v ./cmd/plugin-virtual`
  - Validated by successful compilation of the binary.

## Pre-requisites & Publishing
As a simulation plugin, `plugin-virtual` must be updated whenever the core domain, messaging, storage, or testkit SDKs are changed.

**Before publishing:**
1. Determine current tag: `git tag | sort -V | tail -n 1`
2. Ensure all local tests pass: `go test -v ./...`
3. Ensure the binary builds: `go build -o plugin-virtual ./cmd/plugin-virtual`

**Publishing Order:**
1. Ensure all internal dependencies are tagged and pushed.
2. Update `plugin-virtual/go.mod` to reference the latest tags.
3. Determine next semantic version for `plugin-virtual` (e.g., `v1.0.4`).
4. Commit and push the changes to `main`.
5. Tag the repository: `git tag v1.0.4`.
6. Push the tag: `git push origin main v1.0.4`.

## Update Workflow & Verification
1. **Modify:** Update virtual entity logic or seed data in `app/` or `internal/`.
2. **Verify Local:**
   - Run `go mod tidy`.
   - Run `go test ./...`.
   - Run `go test ./cmd/plugin-virtual` (BDD features).
   - Run `go build -o plugin-virtual ./cmd/plugin-virtual`.
3. **Commit:** Ensure the commit message clearly describes the changes.
4. **Tag & Push:** (Follow the Publishing Order above).
