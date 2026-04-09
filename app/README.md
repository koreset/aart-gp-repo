# AART Group Pricing

**AART Group Pricing** is a standalone cross-platform desktop application for group insurance pricing, administration, bordereaux management, and premium lifecycle management. It is a focused build of the AART platform, scoped exclusively to Group Pricing workflows.

- **Version:** 1.0.0
- **Publisher:** Actuaries and Digital Solutions (Pty) Ltd
- **App ID:** `za.co.adsolutions.aart-gp`

---

## Tech Stack

| Layer     | Technology                                   |
| --------- | -------------------------------------------- |
| Shell     | Electron 33                                  |
| Frontend  | Vue 3 + Vite 5 + TypeScript                  |
| UI        | Vuetify 3 + Material Design Icons            |
| Grids     | AG Grid 31 (Enterprise) + AG Charts 10       |
| State     | Pinia                                        |
| Routing   | Vue Router 4                                 |
| Backend   | Go REST API (shared with AART, port 9090)    |
| IPC       | Electron context bridge (`preload/index.ts`) |
| Packaging | electron-builder 25                          |

---

## Prerequisites

- Node.js >= 18
- Yarn 1.22 (Classic)
- The [AART backend](../aart_project/backend/) running on port 9090

---

## Environment Setup

Create a `.env.development.local` file in the project root before running in development:

```env
VITE_APP_BASE_API_URL=http://localhost:9090/api/v1/
VITE_APP_LICENSE_SERVER=https://licenses.aart-enterprise.com
VITE_APP_AUTH_URL=https://auth.aart-enterprise.com/api/v1
```

For production builds, create `.env.production.local` with the appropriate values.

---

## Development

```bash
# Install dependencies
yarn install

# Start with Vite hot reload
yarn dev

# Debug mode
yarn dev:debug
```

The app opens as an Electron window connecting to `http://localhost:5173`. The backend must be running separately.

---

## Building

The build pipeline runs `prettier → vue-tsc --noEmit → vite build` before packaging.

```bash
# macOS (dmg + zip, arm64 / x64 / universal)
yarn build:mac

# Windows (NSIS installer, x64)
yarn build:win

# Linux
yarn build:linux

# All platforms
yarn build:all

# Directory output only (no installer, useful for testing)
yarn build:dir
```

Output goes to `release/<version>/`.

---

## Code Quality

```bash
yarn lint           # ESLint check
yarn lint:fix       # ESLint auto-fix
yarn format:fix     # Prettier auto-format
```

---

## Testing

```bash
yarn test           # Build then run Playwright E2E tests
yarn test:linux     # Linux with Xvfb virtual display
```

---

## Project Structure

```
aart-gp/
├── src/
│   ├── main/                  # Electron main process
│   │   ├── index.ts           # App lifecycle, window creation
│   │   ├── IPCs.ts            # All IPC handlers (auth, licensing, store)
│   │   ├── MainRunner.ts      # Window management helper
│   │   └── utils/
│   │       ├── Constants.ts   # App name, version, window defaults
│   │       └── encryption.js  # AES-256 encrypt/decrypt for local store
│   ├── preload/
│   │   └── index.ts           # Context bridge — exposes mainApi to renderer
│   └── renderer/              # Vue 3 application
│       ├── main.ts            # Vue app bootstrap
│       ├── App.vue            # Root shell with NavigationDrawer
│       ├── AppLogin.vue       # Login screen
│       ├── AppSetup.vue       # First-run setup wizard
│       ├── api/               # Axios service layer (one file per domain)
│       ├── auth/              # Auth providers (Internal + SSO/Okta)
│       ├── composables/       # Reusable Vue 3 composition functions
│       ├── components/        # Shared UI components
│       ├── constants/         # Static reference data (metadata, etc.)
│       ├── locales/           # i18n JSON files
│       ├── plugins/           # Vuetify + vue-i18n initialization
│       ├── router/            # Vue Router (index.ts)
│       ├── screens/           # Feature screens
│       │   └── group_pricing/ # All GP screens (see below)
│       ├── store/             # Pinia stores
│       ├── types/             # TypeScript type definitions
│       └── utils/             # Formatting and helper utilities
├── buildAssets/
│   ├── builder/config.js      # electron-builder configuration
│   └── icons/                 # App icons (icns, ico, png)
└── .env.development.local     # Local environment variables (not committed)
```

---

## Screens

### Group Pricing Core

| Screen           | Route                                   |
| ---------------- | --------------------------------------- |
| Dashboard        | `/group-pricing/dashboard`              |
| Quote Generation | `/group-pricing/quote-generation`       |
| Quotes List      | `/group-pricing/quotes`                 |
| Quote Output     | `/group-pricing/quotes/output/:quoteId` |
| Actuarial Tables | `/group-pricing/tables`                 |
| Metadata         | `/group-pricing/metadata`               |
| Schemes          | `/group-pricing/schemes`                |
| Scheme Detail    | `/group-pricing/schemes/:id`            |

### Administration

| Screen                 | Route                                             |
| ---------------------- | ------------------------------------------------- |
| Member Management      | `/group-pricing/administration/member-management` |
| Beneficiary Management | `/group-pricing/administration/beneficiaries`     |
| Claims List            | `/group-pricing/claims-list`                      |
| Claims Management      | `/group-pricing/claims-management`                |

### Bordereaux Management

| Screen | Route |
| --- | --- |
| Management Hub | `/group-pricing/bordereaux-management` |
| Generation | `/group-pricing/bordereaux-management/generation` |
| Submission Tracking | `/group-pricing/bordereaux-management/tracking` |
| Reconciliation | `/group-pricing/bordereaux-management/reconciliation` |
| Templates | `/group-pricing/bordereaux-management/templates` |
| Analytics Dashboard | `/group-pricing/bordereaux-management/analytics` |
| Deadline Calendar | `/group-pricing/bordereaux-management/deadline-calendar` |
| Inbound Submissions | `/group-pricing/bordereaux-management/inbound-submissions` |
| Claim Notifications | `/group-pricing/bordereaux-management/claim-notifications` |
| Reinsurer Tracking | `/group-pricing/bordereaux-management/reinsurer-tracking` |
| RI Treaties | `/group-pricing/bordereaux-management/ri-treaties` |
| RI KPI Dashboard | `/group-pricing/bordereaux-management/ri-kpi-dashboard` |
| RI Submission Register | `/group-pricing/bordereaux-management/ri-submission-register` |
| RI Bordereaux | `/group-pricing/bordereaux-management/ri-bordereaux` |
| RI Claims | `/group-pricing/bordereaux-management/ri-claims` |
| RI Technical Accounts | `/group-pricing/bordereaux-management/ri-settlement` |

### PHI Valuations

| Screen         | Route                               |
| -------------- | ----------------------------------- |
| PHI Tables     | `/group-pricing/phi/tables`         |
| Shock Settings | `/group-pricing/phi/shock-settings` |
| Run Settings   | `/group-pricing/phi/run-settings`   |
| Run Results    | `/group-pricing/phi/run-results`    |

### Premium Management

| Screen             | Route                                           |
| ------------------ | ----------------------------------------------- |
| Premium Dashboard  | `/group-pricing/premiums/dashboard`             |
| Schedules          | `/group-pricing/premiums/schedules`             |
| Schedule Detail    | `/group-pricing/premiums/schedules/:scheduleId` |
| Invoices           | `/group-pricing/premiums/invoices`              |
| Invoice Detail     | `/group-pricing/premiums/invoices/:invoiceId`   |
| Payments           | `/group-pricing/premiums/payments`              |
| Reconciliation     | `/group-pricing/premiums/reconciliation`        |
| Arrears Management | `/group-pricing/premiums/arrears`               |
| Statements         | `/group-pricing/premiums/statements`            |

---

## Architecture

```
Vue 3 Screen
  → renderer/api/*.ts (Axios, base URL from electron-store)
  → Go controller  (backend/controllers/)
  → Go service     (backend/services/)
  → GORM → MySQL / PostgreSQL / SQL Server
```

The Electron context bridge (`preload/index.ts`) exposes `window.mainApi` to the renderer. All IPC calls (auth token storage, license validation, app settings) go through this bridge — the renderer never has direct Node.js access.

### State Management

| Store                   | Purpose                                    |
| ----------------------- | ------------------------------------------ |
| `app.ts`                | Active entitlements, global UI state       |
| `auth.ts`               | User session and authentication            |
| `group_pricing.ts`      | Shared GP data (schemes, risk types, etc.) |
| `group_user.ts`         | Group user permissions                     |
| `premium_management.ts` | Premium lifecycle state                    |
| `flash.ts`              | Toast / snackbar messages                  |
| `network_status.ts`     | Online/offline state                       |

### Entitlement-Based UI

Route guards and `v-if` conditions check entitlements loaded into `store/app.ts`. The `entitlementGuard` composable blocks navigation to routes whose `meta.required_permission` the user does not hold.

---

## Backend

The shared AART backend serves this application. See [`../aart_project/backend/README.md`](../aart_project/backend/) or the Swagger UI at `http://localhost:9090/api/v1/swagger/` when the backend is running.

All API responses follow the envelope:

```json
{ "success": true, "data": { ... } }
{ "success": false, "message": "error description" }
```

---

## Auto-Updates

The app is configured for auto-updates via `electron-updater`. Update artifacts are published to `https://updates.aart-enterprise.com/update-gp/`. Use `yarn build:publish` to build and publish a new release.
