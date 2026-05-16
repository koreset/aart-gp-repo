# Broker Portal — Implementation Plan

**Status:** Draft
**Owner:** Engineering
**Last updated:** 2026-05-11
**Audience:** Engineers and tech leads. Business stakeholders should read sections 1, 2, and 12.

---

## 1. Why we are building this

Brokers currently have no way to access AART Group Risk. Quotes attributed to brokers are entered into the system by internal staff on the broker's behalf, with the broker represented only as a name on the quote (`QuoteBroker` struct at `api/models/group_pricing.go:516-519`). This creates three problems:

1. **Internal staff are the bottleneck for new business.** Brokers cannot quote without phoning or emailing the office.
2. **No broker accountability or audit trail.** A broker is a free-text name on a quote, not an authenticated identity.
3. **Distribution is capped.** We cannot scale broker numbers without scaling internal headcount proportionally.

Giving the desktop Electron app to brokers is not viable — installation friction, internal-tool UX, and the security surface of a fat client all rule it out. The answer is a web-based broker portal that talks to the existing Go API.

## 2. Scope

### In scope for v1

- Broker authentication (invite-based onboarding, password + optional MFA)
- Brokerage and broker entity management
- Scheme creation by brokers (for new clients)
- Quote creation against new or existing brokerage-owned schemes
- Quote status visibility (draft → submitted → approved/rejected → accepted → in-force)
- Auto-approval for quotes that meet a configurable rule set; underwriter review otherwise
- Quote document download (PDF/DOCX) once approved
- Member data upload against a scheme
- Multi-tenant isolation: a broker only sees data belonging to their brokerage
- FSCA verification step at broker invitation
- PWA installation support for tablet use

### Out of scope for v1 (deferred)

- Public broker API (key-based access for AMS integrations) — design with this in mind, ship later
- Claims visibility and claim submission
- Commission statements and broker payouts
- Bordereaux upload and management
- Bind / in-force transition workflows
- Scheme transfer between brokerages
- Automated FSCA register integration
- Member self-service and HR admin portals (separate workstream)

### Explicit non-goals

- The portal is **not** a stripped-down web version of the desktop app. It is a different product for a different user.
- Brokers will **not** have access to internal underwriting controls (custom TIR tables, calculation method overrides, credibility settings).
- The portal will **not** support offline-first quote creation in v1; PWA installability is for convenience, not data sync.

---

## 3. User personas

**Broker (primary user)**
FSCA-registered individual at a brokerage firm. Quotes group risk cover for employer clients. Spends most of their time outside the office — laptop and tablet, sometimes at a client site. Expects fast quote turnaround. Needs to download and present a professional quote document the same day a client briefs them.

**Brokerage administrator**
Senior broker or operations person at the firm. In v1 has the same capabilities as a regular broker but sees firm-wide data. In a later phase will manage their firm's broker roster and approval limits.

**Internal underwriter (existing user, new workflow)**
Receives broker-submitted quotes that exceed auto-approval thresholds. Reviews, requests changes via comments, approves or rejects. Continues to use the existing desktop app.

**Internal admin (existing user, new workflow)**
Onboards new brokerages and brokers. Verifies FSCA registration. Sets per-brokerage approval rules.

---

## 4. Architecture overview

The portal is a new SPA in a sibling workspace `broker-portal/`, communicating with the existing Go API via a new namespace `/api/v2/broker/*`. The same Go binary serves both internal and broker endpoints — there is no separate service. Audience separation is enforced at the middleware layer via a new JWT `aud` claim.

```
                 ┌──────────────────────┐
                 │  Electron desktop    │
                 │  app/                │  internal users
                 └──────────┬───────────┘
                            │ /api/v2/...
                            │ (aud: internal)
                            │
                 ┌──────────▼───────────┐
                 │   Go API (api/)      │
                 │   Gin + GORM         │
                 │                      │
                 │   middleware:        │
                 │     RequireAud       │
                 │     RequirePerm      │
                 │     ScopeToBrokerage │
                 └──────────▲───────────┘
                            │ /api/v2/broker/...
                            │ (aud: broker)
                            │
                 ┌──────────┴───────────┐
                 │  Broker portal       │
                 │  broker-portal/      │   brokers
                 │  (Vue 3 SPA / PWA)   │
                 └──────────────────────┘
```

Shared code between the two frontends — the API client, validators, formatters, and a small set of UI primitives — lives in `packages/` workspaces consumed by both `app/` and `broker-portal/`. This is added as part of phase 3.

## 5. Data model changes

### New tables

**`brokerages`**

| Column | Type | Notes |
|---|---|---|
| `id` | int, PK | |
| `name` | varchar(255) | Display name |
| `fsca_license_number` | varchar(64) | The firm's FSCA Financial Services Provider number |
| `fsca_verified_at` | datetime, nullable | Set when internal admin verifies |
| `fsca_verified_by` | int, FK → org_users.id | Named internal user |
| `status` | enum | `pending`, `active`, `suspended`, `terminated` |
| `approval_rules` | json, nullable | Auto-approval config; NULL means use system default |
| `primary_contact_email` | varchar(255) | |
| `created_at`, `updated_at` | datetime | |

**`brokers`**

| Column | Type | Notes |
|---|---|---|
| `id` | int, PK | |
| `brokerage_id` | int, FK → brokerages.id, NOT NULL | Every broker belongs to a firm |
| `name` | varchar(255) | |
| `email` | varchar(255), unique | Login identity |
| `fsca_rep_number` | varchar(64) | Individual representative number |
| `fsca_verified_at` | datetime, nullable | |
| `fsca_verified_by` | int, FK → org_users.id | |
| `status` | enum | `invited`, `active`, `suspended`, `terminated` |
| `password_hash` | varchar(255), nullable | NULL until first login |
| `invitation_token` | varchar(64), nullable | One-shot magic link |
| `invitation_expires_at` | datetime, nullable | |
| `last_login_at` | datetime, nullable | |
| `created_at`, `updated_at` | datetime | |

**`broker_audit_log`**

Append-only record of broker actions. At minimum: `id`, `broker_id`, `action`, `entity_type`, `entity_id`, `payload` (json), `ip_address`, `user_agent`, `created_at`. Used for compliance and for reconstructing why an auto-approval fired.

### Modified tables

**`group_schemes`** — add `brokerage_id` (int, FK, nullable). NULL means internally-acquired. Populated automatically when a broker creates a scheme.

**`group_pricing_quotes`** — add:
- `brokerage_id` (int, FK, nullable) — denormalised from scheme for query speed and to handle scheme transfer cleanly
- `created_by_broker_id` (int, FK → brokers.id, nullable) — who in the brokerage created it
- `auto_approved_at` (datetime, nullable)
- `auto_approval_rule_snapshot` (json, nullable) — frozen copy of the rule that triggered auto-approval

The existing `QuoteBroker` embedded struct (`broker_id`, `broker_name`) stays in place for backward compatibility with existing quotes but is deprecated. Newly-created broker quotes derive these fields from the FK.

### Migration strategy

Use the existing diff-based migration system (`api/tools/generate_migration.go` → `api/migrations/<dialect>/`, applied by `services.RunMigrationsOnStartup`). All new columns nullable; backfill is unnecessary for v1 since existing quotes have no broker entity to attach. Add a check at boot time that warns if `brokerages` or `brokers` is empty when broker routes are enabled — this catches misconfigured environments.

---

## 6. Authentication and authorization

### JWT audience separation

Add an `aud` claim to all issued tokens. Existing internal-user tokens get `aud: "internal"`; broker tokens get `aud: "broker"`. Add a middleware `RequireAudience(aud string)` mounted on the appropriate route groups in `api/routes/routes.go`. Existing internal middleware (`GetActiveUser`, in `api/routes/middleware.go:140-226`) continues to work for internal endpoints, gaining only a check that `aud == "internal"`. A new `GetActiveBroker` middleware handles broker tokens — looks up the broker by email, attaches `broker` and `brokerage_id` to the Gin context.

The two audiences share no endpoints. A broker token presented to an internal endpoint returns 401, and vice versa. This is non-negotiable: it is the foundation of audience isolation.

### Brokerage scoping

A second middleware `ScopeToBrokerage` reads `brokerage_id` from the Gin context and exposes it via a context helper used by every broker service. The rule: **no broker SQL query is allowed to omit the `brokerage_id` filter.** Enforce this with a code review checklist and, ideally, a small linter or test that scans broker service files for raw queries missing the filter.

### Permission slugs

New broker permission slugs prefixed `broker:*` to keep the namespace separate from internal slugs:

- `broker:scheme.create`
- `broker:scheme.view`
- `broker:quote.create`
- `broker:quote.view`
- `broker:quote.submit`
- `broker:quote.export`
- `broker:member.upload`

For v1 every active broker has all of these. The slug structure exists so that the brokerage admin role (later) can grant subsets.

### Internal-side new slugs

- `broker_admin:brokerage.manage` — create/edit brokerages, invite brokers, set approval rules
- `broker_admin:quote.review` — approve or reject submitted broker quotes
- `broker_admin:quote.comment` — exchange messages with brokers on a quote

### Onboarding flow

1. Internal admin (with `broker_admin:brokerage.manage`) creates a `brokerages` row, captures and verifies the FSCA license number, sets `status = active`.
2. Internal admin creates `brokers` rows for each individual broker the firm wants onboarded, captures FSCA rep number, sets `status = invited`, generates an `invitation_token` valid for 7 days.
3. Email service sends the broker a link `https://broker.aart-enterprise.com/accept-invite?token=...`.
4. Broker sets a password, optionally enables MFA, status flips to `active`.

### MFA

TOTP-based (Google Authenticator, Authy) is the recommended baseline. Make it optional in v1, mandatory for brokerage admins. Plan to make it mandatory for all brokers once the portal is past pilot.

### Password and account hygiene

Standard set: bcrypt hashing (cost 12), rate-limited login (5 attempts per 15 min per email + per IP), password reset by emailed token, session JWTs with 1-hour expiry and refresh token rotation, no concurrent session limit in v1.

---

## 7. API surface

All endpoints under `/api/v2/broker/*`. Designed as a versioned, documented public API from day one — same shape that would later be exposed to broker partners via API keys.

### Auth

- `POST /api/v2/broker/auth/accept-invite` — body: `{ token, password }` → sets password, returns JWT
- `POST /api/v2/broker/auth/login` — body: `{ email, password, mfa_code? }` → returns JWT + refresh token
- `POST /api/v2/broker/auth/refresh` — body: `{ refresh_token }` → new JWT
- `POST /api/v2/broker/auth/password-reset` — body: `{ email }` → sends reset email
- `POST /api/v2/broker/auth/password-reset/confirm` — body: `{ token, new_password }`

### Self

- `GET /api/v2/broker/me` — returns current broker + brokerage profile
- `GET /api/v2/broker/me/brokerage` — brokerage details (read-only in v1)

### Schemes

- `GET /api/v2/broker/schemes` — list firm's schemes
- `POST /api/v2/broker/schemes` — create new scheme, `brokerage_id` auto-assigned from token
- `GET /api/v2/broker/schemes/:id` — fetch (404 if not owned by token's brokerage)
- `PUT /api/v2/broker/schemes/:id` — edit limited fields
- `POST /api/v2/broker/schemes/:id/members` — upload member data (multipart CSV or JSON)

### Quotes

- `GET /api/v2/broker/quotes` — list firm's quotes, paginated, filterable by status/scheme
- `POST /api/v2/broker/schemes/:scheme_id/quotes` — create quote
- `GET /api/v2/broker/quotes/:id` — fetch quote
- `PUT /api/v2/broker/quotes/:id` — edit while in `draft`
- `POST /api/v2/broker/quotes/:id/calculate` — trigger pricing engine (proxies to existing `calculate-quote` job)
- `GET /api/v2/broker/quotes/:id/result-summary` — calculation results
- `POST /api/v2/broker/quotes/:id/preview-approval` — pre-submission check: would this auto-approve, and why/why not?
- `POST /api/v2/broker/quotes/:id/submit` — submit; runs approval evaluator, moves to `approved` or `submitted`
- `GET /api/v2/broker/quotes/:id/document.pdf` — export, only if status in `[approved, accepted, in_force]`
- `GET /api/v2/broker/quotes/:id/document.docx` — same
- `GET /api/v2/broker/quotes/:id/comments` — thread with underwriter
- `POST /api/v2/broker/quotes/:id/comments` — add comment

### Notifications

- `GET /api/v2/broker/notifications` — list
- `POST /api/v2/broker/notifications/:id/read` — mark read

### Shared service layer

Where calculation logic lives in `api/services/`, broker endpoints reuse it without forking. Controllers under `api/controllers/broker/` are thin: validate input, enforce brokerage scope, call shared service, marshal output.

### OpenAPI / Swagger

Continue using `swag init` to generate docs. Tag every broker endpoint with `broker` so the generated spec can be split into two documents — `swagger-internal.json` and `swagger-broker.json` — making the future public API a publication step rather than an extraction job.

---

## 8. Auto-approval

### Rule schema

Stored as JSON on `brokerages.approval_rules`, defaulting to a system-wide config if NULL:

```json
{
  "max_member_count": 50,
  "max_total_sum_insured": 50000000,
  "allowed_benefits": ["GLA", "Funeral"],
  "max_average_age": 55,
  "excluded_industry_codes": ["mining_underground", "deep_sea_fishing"],
  "require_clean_claims_history": true
}
```

A quote auto-approves only if **all** conditions pass. Any failure routes it to internal underwriter review.

### Evaluator

`services.EvaluateAutoApproval(quote, rules) (eligible bool, reasons []string)` in `api/services/broker_approval.go`. Pure function over the quote, member data, and the rule snapshot. Returns the failure reasons (e.g., `"member_count 78 exceeds max 50"`) which become both UI feedback and the audit-log entry.

Called from two endpoints:
- `POST /quotes/:id/preview-approval` — so the broker UI can show "this quote will auto-approve" or "this will go to underwriter — reasons: ..." before submission
- `POST /quotes/:id/submit` — actual decision point

### Audit

When a quote auto-approves, store `auto_approved_at` and a `auto_approval_rule_snapshot` JSON blob containing the rules that were in effect at the moment of approval. This lets compliance reconstruct *why* an auto-approval was granted even after rules change.

### Edge cases

- A broker resubmits an edited quote that previously auto-approved. The evaluator runs again. If it no longer auto-approves, status reverts to `submitted` and the previous auto-approval is invalidated (logged).
- A rule change while a quote is in flight does not retroactively affect approved quotes.
- Internal admin can manually override auto-approval by sending a quote back to `pending_review`.

---

## 9. Quote and scheme workflows

### Scheme creation (broker)

1. Broker fills slim form: employer name, industry classification, anticipated headcount, intended effective date, primary contact.
2. Backend creates `group_schemes` row with `brokerage_id` from token.
3. Broker is taken to the scheme detail page with a "Upload members" call to action.

The internal scheme creation flow remains unchanged and continues to expose all the configuration that brokers don't need.

### Quote lifecycle

The existing status enum (`api/models/group_pricing.go:60-80`) is sufficient:

```
draft → submitted → pending_review → approved → accepted → in_force
                                  ↘ rejected
                                  ↘ cancelled
```

Broker-facing transitions:

- `draft → submitted` — broker submits. Evaluator decides next state.
- `submitted → approved` — auto-approval path (no human review).
- `submitted → pending_review` — sent to internal underwriter.
- `pending_review → approved | rejected` — underwriter decision (internal app).
- `approved → accepted` — broker marks client has accepted the quote.
- `accepted → in_force` — internal-only, once cover is bound.
- Any non-terminal state → `cancelled` — broker or internal.

### Comment thread per quote

Existing schema doesn't have one; add `quote_comments` table: `id`, `quote_id`, `author_type` (`internal | broker`), `author_id`, `body`, `created_at`. Brokers can comment on their own quotes; internal underwriters can comment on any. WebSocket push for real-time when both parties are connected.

### Notifications

Email on the following events to the broker who created the quote:
- Quote auto-approved
- Quote routed to underwriter
- Quote approved or rejected by underwriter
- New comment from underwriter
- Quote expiring (configurable, e.g., 7 days before)

Email on the following events to internal underwriting team:
- New quote awaiting review
- Broker comment on a quote in `pending_review`

Email service can be the existing one used elsewhere in the API; if no such abstraction exists, add `api/services/email/` and start there.

---

## 10. Frontend architecture

### Workspace layout

```
aart-group-risk/
├── api/                    (unchanged)
├── app/                    (Electron, unchanged)
├── broker-portal/          NEW
│   ├── src/
│   ├── public/
│   ├── vite.config.ts
│   └── package.json
├── packages/               NEW
│   ├── api-client/         shared axios client + typed responses
│   ├── shared-ui/          shared Vuetify components (logo, layout primitives)
│   └── shared-utils/       date, currency, validators, formatters
└── docs/
    └── broker-portal-plan.md   (this file)
```

The `packages/` workspace structure may already partially exist via yarn workspaces in `app/`; verify before adding root-level `package.json` workspaces config.

### Stack

- Vue 3 with `<script setup>` and TypeScript — matches `app/` team expertise.
- Vuetify 3 for layout and form components — same as `app/`, lower learning curve.
- Pinia for state, Vue Router for routing, Vuelidate for form validation — same as `app/`.
- Vite build, SPA output.
- `vite-plugin-pwa` for service worker and manifest. Installable on iPadOS and Android tablets out of the box.

### Authentication in the SPA

JWT in memory (Pinia store), refresh token in `httpOnly` cookie set by the API. On `401`, attempt refresh once before redirecting to login. No tokens in `localStorage`.

### Routing structure

```
/login
/accept-invite/:token
/password-reset
/                  → /dashboard
/dashboard         (KPIs: in-flight quotes, recent activity)
/schemes
/schemes/new
/schemes/:id
/schemes/:id/members
/quotes
/quotes/new        → really /schemes/:id/quotes/new
/quotes/:id
/profile
```

### Shared components extracted from `app/`

To avoid drift, the following are reasonable candidates to extract into `packages/shared-ui/` over phases 3-4:

- API axios instance with interceptors
- Currency / date formatters
- Member CSV parsing and validation
- Quote summary card
- Benefit configuration display widget

Don't try to extract everything. Anything that has different desktop vs. broker UX (full quote editor, member grids, dashboards) stays bespoke.

---

## 11. Phased delivery plan

Each phase is intended to land independently and be reviewable in isolation. Estimates are rough — they assume one mid-senior backend dev and one mid-senior frontend dev on the work, both with context on the existing codebase. Adjust for your team.

### Phase 0 — Data model and migrations (≈1 week)

- Generate migrations for `brokerages`, `brokers`, `broker_audit_log`, `quote_comments` tables
- Add `brokerage_id` to `group_schemes` and `group_pricing_quotes`
- GORM models in `api/models/broker.go` and additions to existing models
- Seed script for one test brokerage and two test brokers (dev only)
- No public-facing behaviour change yet

### Phase 1 — Auth foundation (≈1.5 weeks)

- `aud` claim added to JWT issuer
- `RequireAudience` middleware
- `GetActiveBroker` middleware, sets context
- Broker invitation flow: internal admin endpoint to create + email invite
- `POST /accept-invite`, `POST /login`, `POST /refresh` for brokers
- Password reset
- TOTP MFA scaffolding (optional toggle per broker)
- Internal admin UI in the desktop app for creating brokerages and inviting brokers (gated behind `broker_admin:brokerage.manage`)

### Phase 2 — Broker API (≈2 weeks)

- `/api/v2/broker/schemes` CRUD
- `/api/v2/broker/quotes` CRUD
- Calculate, result-summary, document export endpoints (thin wrappers over existing services)
- `EvaluateAutoApproval` service and preview-approval endpoint
- Submit endpoint with state transitions and audit
- Comments endpoints
- Notifications scaffolding (email send on key events)
- Swagger tagging and split

### Phase 3 — Portal skeleton (≈1 week, parallel with phase 2)

- `broker-portal/` workspace, Vite + Vue + Vuetify scaffold
- `packages/api-client/`, `packages/shared-utils/` initial extractions
- Login, accept-invite, password-reset screens
- Auth Pinia store and route guards
- App shell, navigation, profile page
- Deployable to staging behind basic auth

### Phase 4 — Quote and scheme flows (≈3 weeks)

- Scheme list, scheme creation wizard, scheme detail
- Member CSV upload with client-side preview and server-side validation
- Quote creation wizard (slim version of `NewQuoteDetail.vue`)
- Quote list with filters
- Quote detail with status indicator
- Pre-submission approval preview UI
- Submit flow with appropriate next-step UI
- Document download

### Phase 5 — Workflow loop and polish (≈1.5 weeks)

- Comments UI with WebSocket push
- Notifications inbox
- Internal underwriter review queue in desktop app (`broker_admin:quote.review`)
- Approve / reject / send-back actions from internal app
- PWA manifest, service worker, install prompt
- Production deployment, DNS, TLS, WAF rules
- Pilot with 2-3 friendly brokerages

### Phase 6 — Hardening before general availability (≈1 week + bake time)

- Penetration test
- Load test the broker API (target: 100 concurrent brokers, 1000 quotes/day)
- Rate limit tuning
- Audit log review with compliance
- Documentation for brokerages (admin and end-user)
- Training material

**Total elapsed:** roughly 10–11 weeks of engineering, plus pilot bake time. Phase 2 and 3 run in parallel.

---

## 12. Risks and open questions

### Risks

1. **Multi-tenant data leakage.** A single missing `brokerage_id` filter on a broker endpoint exposes another firm's data. Mitigations: code review checklist, integration tests that verify cross-tenant 404s, a static-analysis pass over `api/controllers/broker/` looking for service calls without scope.

2. **FSCA compliance scope.** Plan assumes invitation-time manual verification is sufficient. Compliance team should confirm whether ongoing verification, register integration, or specific record-keeping rules apply.

3. **Calculation engine reuse breaks under broker concurrency.** The existing pricing engine was sized for internal users. Broker submissions may spike at month-end or renewal periods. Phase 6 load test will surface this; mitigation may be additional queue workers.

4. **Drift between Electron app and portal.** Two frontends backed by the same API can develop subtle inconsistencies. Mitigations: shared `packages/api-client/` for types, shared validation rules, and a single OpenAPI spec as the contract.

5. **Auto-approval rule misconfiguration.** A too-permissive rule lets risky quotes through without human review. Mitigations: rule changes require `system:admin` permission, every auto-approval logs the rule snapshot, monthly compliance review of approved quotes.

### Open questions

These need answers before phases 4-5 finalize:

- **Scheme transfer between brokerages.** What happens when an employer switches brokers? Options: (a) scheme stays with original brokerage, new firm creates a new scheme; (b) scheme reassigned by internal admin with full history retained; (c) scheme history split. Recommendation: defer the implementation but design the audit log to support reconstructing scheme ownership over time.

- **Brokerage admin role.** Is "brokerage admin" a distinct permission level in v1, or is it acceptable that every broker sees firm-wide quotes? Current plan assumes the latter for simplicity; a brokerage-admin tier is a v1.5 addition.

- **Pricing transparency.** Should brokers see the full pricing methodology (loading factors, credibility weights), or only the final quote? Internal users see everything; broker visibility is a business call.

- **Re-quote vs. new-quote.** When a broker tweaks a submitted quote, is it a new quote or a revision of the existing one? Affects audit trail and comment thread continuity.

- **MFA enforcement timeline.** When does optional MFA become mandatory? Recommendation: optional for v1 pilot, mandatory before general availability.

- **Quote validity period.** How long does an approved quote remain valid before it must be re-priced? Existing internal practice should set the default.

- **Document branding.** Does the exported quote document need brokerage branding (logo, footer), or always AART Group Risk branding? Affects the document template service.

### Decisions made

- Brokerage is a first-class entity from v1, not a v1.5 addition.
- Brokers see all quotes belonging to their brokerage.
- Auto-approval is configurable per brokerage; rule schema as specified in section 8.
- FSCA verification is manual at onboarding, with annual re-verification reminders.
- Future broker-facing public API is not in v1 scope but the v1 endpoints are designed to become that API.

---

## 13. Appendix — references into the existing codebase

Key files to read before starting each phase:

- Quote model: `api/models/group_pricing.go` — status enum at lines 60-80, `GroupPricingQuote` struct from line 107, `QuoteBroker` embedded struct at lines 516-519
- Quote controller: `api/controllers/group_pricing.go` (lines 40-99 for quote creation entry point)
- Routes: `api/routes/routes.go` (lines 59-189 cover group pricing endpoints; this is where the `/broker/*` namespace will be mounted)
- Auth middleware: `api/routes/middleware.go` (lines 140-308 cover `GetActiveUser`, `RequirePermission`, permission lookup)
- Migrations system: `api/MIGRATIONS.md` and `api/tools/generate_migration.go`
- Existing quote frontend: `app/src/renderer/screens/group_pricing/NewQuoteDetail.vue` — the canonical reference for what a quote-editing UI looks like today; the broker portal version is a slim derivative
- Scheme frontend: `app/src/renderer/screens/group_pricing/GroupSchemeList.vue` and `GroupSchemeDetail.vue`

---

*End of plan. Comments and suggested edits welcome via PR or in the project chat.*
