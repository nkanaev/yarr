# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

- **Build for host OS/architecture**: `make host` (outputs to `out/yarr`)
- **Run development server**: `make serve` (runs with local.db)
- **Run tests**: `make test` (runs all Go tests with SQLite tags)
- **Cross-compile CLI versions**:
  - Linux: `make linux_amd64`, `make linux_arm64`, `make linux_armv7`
  - Windows: `make windows_amd64`, `make windows_arm64`
  - macOS: `make darwin_amd64`, `make darwin_arm64`
- **Build GUI versions** (must run on target OS):
  - `make darwin_arm64_gui` (creates yarr.app)
  - `make darwin_amd64_gui` (creates yarr.app)
  - `make windows_amd64_gui` (creates yarr.exe)
  - `make windows_arm64_gui` (creates yarr.exe)

### Prerequisites
- Go >= 1.23
- C Compiler (GCC/Clang)
- Zig >= 0.14.0 (for cross-compilation)
- binutils (for Windows GUI builds)

## Project Architecture

This is a Go-based RSS feed reader with an embedded SQLite database and web interface.

### Core Components

**Main Application** (`cmd/yarr/main.go`)
- Entry point that parses command-line flags and environment variables
- Handles authentication file parsing
- Integrates all major components

**Server Layer** (`src/server/`)
- HTTP server with routing and middleware
- Fever API compatibility for third-party clients
- Authentication and TLS support
- Static asset serving from embedded files

**Storage Layer** (`src/storage/`)
- SQLite database with WAL mode and foreign key constraints
- Models for feeds, items, folders, and settings
- Database migration system
- Built with tags: `sqlite_foreign_keys sqlite_json`

**Feed Processing** (`src/parser/` and `src/worker/`)
- Multi-format feed parser (RSS, Atom, JSON Feed, RDF)
- Concurrent feed crawler with configurable worker pool (NUM_WORKERS = 4)
- Content scraping and readability extraction
- HTML sanitization and URL resolution

**Platform Integration** (`src/platform/`)
- Cross-platform GUI support via system tray
- OS-specific file opening and console handling
- Conditional compilation with build tags (`gui` tag for GUI builds)

**Content Processing** (`src/content/`)
- HTML parsing and manipulation utilities
- Readability extraction for full-text content
- Content sanitization and security filtering
- Media handling and URL resolution

### Key Build Tags and Flags

The project uses specific Go build tags and linker flags:
- Tags: `sqlite_foreign_keys sqlite_json` (always), `gui` (for GUI builds)
- Linker flags inject version info: `-X 'main.Version=$(VERSION)' -X 'main.GitHash=$(GITHASH)'`
- CGO is required (`CGO_ENABLED=1`) for SQLite integration

### Database

Uses SQLite with optimized settings:
- WAL journal mode for better concurrency
- NORMAL sync mode for performance
- 5-second busy timeout
- Shared cache enabled

### Testing

Tests are spread across modules and can be run with `make test`. Key test files include feed parsing, storage operations, and routing logic.