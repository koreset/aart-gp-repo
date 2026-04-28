# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AART Group Risk is a monorepo with two main components:
- **`api/`** — Go backend (Gin framework, GORM ORM, multi-database: MySQL/PostgreSQL/SQL Server)
- **`app/`** — Electron desktop app (Vue 3 + Vuetify 3 + Vite + TypeScript)

The application handles group pricing, valuation jobs, bordereaux management, claims, and insurance-related workflows for AD Solutions.

## Common Commands

### API (Go backend, run from `api/`)

```bash
# Development with hot reload
air                          # uses .air.toml config

# Build
go build -o aart_api

# Run tests
go test ./...                # all tests
go test ./services/...       # single package

# Generate Swagger docs
swag init
```

### App (Electron frontend, run from `app/`)

```bash
# Development
yarn dev

# Build (runs format, type-check, then electron-builder)
yarn build                   # current platform
yarn build:mac               # macOS
yarn build:win               # Windows
yarn build:linux             # Linux

# Lint & Format
yarn lint                    # check
yarn lint:fix                # auto-fix
yarn format                  # check formatting
yarn format:fix              # auto-fix formatting

# Tests
yarn test:unit               # vitest run
yarn test:unit:watch         # vitest in watch mode
yarn test                    # e2e (builds first, then playwright)
```

## Architecture

### API

- **Entry point:** `api/main.go` — sets up Gin server, CORS, middleware, WebSocket, and can install/uninstall as a system service
- **Routes:** `api/routes/routes.go` — large route file defining all API endpoints grouped by feature (valuation jobs, group pricing, bordereaux, claims, etc.)
- **Controllers:** `api/controllers/` — request handlers, one file per domain
- **Services:** `api/services/` — business logic layer
- **Models:** `api/models/` — GORM model definitions
- **Migrations:** `api/migrations/<dialect>/` — generated SQL migration files. Schema changes are diff-based: `tools/generate_migration.go` produces a delta against the connected DB; `services.RunMigrationsOnStartup` applies pending files on every boot, recording versions in the `migrations` table. See `api/MIGRATIONS.md` for the full system, dev workflow, and remote-DB usage.
- **Config:** `api/config/config.go` — DB config via env vars (`DB_HOST`, `DB_USER`, `DB_PORT`, `DB_PWD`, `DB_NAME`)
- **Auth:** JWT-based authentication; protected routes use `GetActiveUser()` middleware
- **WebSocket:** endpoint at `/ws` with query-param auth

### App

- **Electron structure:** main process (`src/main/`), preload (`src/preload/`), renderer (`src/renderer/`)
- **Renderer entry:** `src/renderer/main.ts` — initializes Vue + Pinia + Vue Router + Vuetify + i18n + Vuelidate
- **Root components:** `App.vue` (main), `AppLogin.vue`, `AppSetup.vue`, `AppNoInternet.vue`, `LicenseWindow.vue`
- **Screens:** `src/renderer/screens/` — feature screens including `group_pricing/` (largest, ~27 subcomponents)
- **Stores:** `src/renderer/stores/` — Pinia stores (app, group_pricing, conversations, notifications, etc.)
- **Composables:** `src/renderer/composables/` — reusable composition functions (~20 domains)
- **Services:** `src/renderer/services/` — API client layer (axios-based)
- **i18n:** `src/renderer/locales/` — 13 language directories
- **Path alias:** `@/` maps to `src/` in TypeScript config
- **License check:** Bootstrap flow validates license on startup (states: VALID, EXPIRED, SUSPENDED, OVERDUE, etc.)
- **IPC:** Context bridge exposes `mainApi` and `electronAPI` to renderer

### Deployment

- GitHub Actions CI/CD (`api/.github/workflows/deploy.yml`) deploys the Go API to 3 service instances (app1, app2, app3) on push to `main`/`master`
- Electron app distributes via electron-builder with auto-updates from `updates.aart-enterprise.com`

## Key Conventions

- Go: follows modern Go guidelines (enabled via Claude marketplace plugin)
- Frontend: Vue 3 `<script setup>` syntax with TypeScript
- Data grids/charts use AG Grid and AG Charts (enterprise editions)
- Multiple date libraries in use: date-fns, moment, luxon
- Package manager: yarn (app has `.yarnrc`)
