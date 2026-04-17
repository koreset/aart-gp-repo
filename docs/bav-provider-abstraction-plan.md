# Bank Account Verification — Provider Abstraction Implementation Plan

**Status:** Proposed
**Owner:** Backend (api/) — with coordinating changes in app/
**Target release:** TBD
**Last updated:** 2026-04-17

## 1. Context

The current bank account verification path is hard-wired to **VerifyNow** (`api/utils/verifynow.go`). Each insurer deployment of AART Group Risk will use a different BAV provider (LexisNexis, TransUnion, Experian AVS, XDS, ThisIsMe, direct BankservAfrica AVS, or VerifyNow). AART does not supply the BAV service — we integrate into whichever the insurer already holds a contract with.

The good news: in the South African market, virtually every BAV provider front-ends **BankservAfrica AVS**, so the *semantic* contract is standardised. The delta across providers is auth scheme, wire format, sync vs async, and error shape — all solvable by an adapter layer.

## 2. Goals

1. Remove VerifyNow coupling from controllers, services, frontend, and DB records.
2. Introduce a provider-agnostic canonical model that the rest of the stack depends on.
3. Support selection of a single active provider per deployment via configuration.
4. Support sync and async providers without requiring a second refactor later.
5. Capture raw provider payloads for audit/dispute resolution.
6. Zero-downtime rollout — VerifyNow-backed deployments keep working with no config changes on day one.

## 3. Non-goals (for this plan)

- Multi-tenant provider selection within a single deployment. Deferred until an insurer actually requires it.
- Batch BAV. Current flow is single-claim verification only.
- Replacing the existing ACB file generation path — unrelated.

## 4. Architecture summary

New Go package `api/services/bav/` owns the abstraction:

```
api/services/bav/
├── types.go           # VerifyRequest, VerifyResult, TriState, Status enums
├── provider.go        # Provider interface, Registry, errors
├── registry.go        # Loads active provider from config
├── logging.go         # bav_verification_log persistence helpers
└── providers/
    ├── verifynow.go   # ported from utils/verifynow.go
    ├── lexisnexis.go  # (written on demand)
    ├── transunion.go  # (written on demand)
    ├── mock.go        # deterministic, for tests + local dev
    └── common/
        ├── oauth2.go  # shared OAuth2 client-credentials helper
        └── tristate.go# "Y"/"N"/"yes"/"no"/bool → TriState normalisation
```

Controller (`controllers/claim_payments.go`) calls `bav.Active().Verify(ctx, req)` instead of `utils.VerifyBankAccount(...)`. `utils/verifynow.go` is removed at the end of Phase 3.

## 5. Phased delivery

Each phase is independently shippable. Phases 1–5 must ship before any second provider is integrated; Phases 6–8 are driven by real insurer requirements.

### Phase 1 — Canonical domain model (no behaviour change)

**Scope**
- Create `api/services/bav/types.go` with `VerifyRequest`, `VerifyResult`, `TriState` (`TriYes`/`TriNo`/`TriUnknown`), `Status` (`StatusComplete`/`StatusPending`/`StatusFailed`).
- Create `api/services/bav/provider.go` with the `Provider` interface, `Registry` struct, and typed errors (`ErrUnauthorized`, `ErrRateLimited`, `ErrProviderUnavailable`, `ErrInvalidInput`).
- No caller changes yet.

**Acceptance**
- `go build ./...` passes.
- Unit tests for `TriState` parsing (`"Y"`, `"N"`, `"yes"`, `"no"`, `true`, `false`, `""`, `"unknown"`) in `types_test.go`.

**Estimated effort:** 0.5 day.

### Phase 2 — Port VerifyNow behind the interface

**Scope**
- Create `api/services/bav/providers/verifynow.go` implementing `Provider`.
- Move all HTTP, auth, and marshalling logic from `api/utils/verifynow.go` into the new adapter.
- Translate VerifyNow's `"Y"/"N"` response fields into `TriState` via `common/tristate.go`.
- Preserve the raw response bytes on `VerifyResult.RawPayload`.
- Add unit tests using `httptest.Server` covering: happy path, 401 → `ErrUnauthorized`, 429 → `ErrRateLimited`, 5xx → `ErrProviderUnavailable`, malformed JSON.

**Acceptance**
- New adapter has ≥80% line coverage.
- Existing VerifyNow integration test (if present) migrated; otherwise add one.

**Estimated effort:** 1 day.

### Phase 3 — Registry, config plumbing, controller cutover

**Scope**
- Extend `api/config/config.go` with a `BAV` sub-struct:
  - `BAV_PROVIDER` (default `verifynow`)
  - `BAV_BASE_URL` (override; defaults per-provider)
  - `BAV_API_KEY`, `BAV_OAUTH_CLIENT_ID`, `BAV_OAUTH_CLIENT_SECRET`, `BAV_OAUTH_TOKEN_URL`
  - `BAV_TIMEOUT_SECONDS` (default 45)
- Backward-compat aliases: if `BAV_API_KEY` is empty fall back to `VERIFYNOW_API_KEY`; if `BAV_PROVIDER` is empty and `VERIFYNOW_API_KEY` is set, treat as `verifynow`.
- Implement `bav.NewRegistry(cfg)` wiring the active provider once at startup (called from `main.go` or a fresh `api/bootstrap` package).
- Rewrite `controllers/claim_payments.go::VerifyBankAccount` to build a `bav.VerifyRequest`, inject an idempotency key derived from `claim_id + attempt_number` (see Phase 5), and call `bav.Active().Verify(ctx, req)`.
- Translate `VerifyResult` to the existing frontend JSON shape for now — **do not** change the wire contract yet (Phase 4 owns that).
- Delete `api/utils/verifynow.go`.

**Acceptance**
- No deployment requires a config change to keep working (old `VERIFYNOW_*` env vars still take effect).
- Existing e2e test that calls `/group-pricing/claims/verify-bank-account` still passes.
- Grep for `VerifyNow` in `api/controllers/`, `api/services/` (outside `bav/`), `api/utils/` returns zero results.

**Estimated effort:** 1 day.

### Phase 4 — Canonical frontend contract

**Scope**
- Introduce `app/src/renderer/types/bav.ts` with a TypeScript mirror of `VerifyResult` (`TriState` as `'yes' | 'no' | 'unknown'`).
- Update `GroupPricingService.ts::verifyBankAccount` return type.
- Update the consuming Vue component (trace callers of `verifyBankAccount` — there should be ≤2) to render `TriState` via a shared helper that maps to tick/cross/question icons.
- Update the Go controller response to emit the canonical shape instead of the VerifyNow-shaped response. Bump the endpoint behind a minor version: `POST /v2/group-pricing/claims/verify-bank-account`, and keep `/v1/...` returning the legacy shape for one release to avoid a lockstep deploy of api+app.

**Acceptance**
- Type check passes (`yarn type-check`).
- Manual verification in dev: verify a known-good account, known-bad account, and an account where bank returns "unknown" for `accountOpen` — UI must visually differentiate the three states.

**Estimated effort:** 1.5 days.

### Phase 5 — Audit logging + deterministic idempotency

**Scope**
- New GORM model `api/models/bav_verification_log.go` with fields: `ID`, `ClaimID`, `Provider`, `ProviderRequestID`, `IdempotencyKey`, `Status`, `RequestPayload` (JSON string — use `datatypes.JSON` for cross-DB portability since the project supports MySQL/Postgres/SQL Server), `ResponsePayload`, `ErrorMessage`, `CreatedAt`.
- Register the model in `api/services/migrations.go` so `AutoMigrate` picks it up on next boot.
- `bav.Registry.Verify(...)` wraps each adapter call with log-before and log-after, writing both success and failure rows.
- Idempotency key = `sha256(claim_id + attempt_number + provider)`. Attempt number increments on retry.
- Replace the UUID generated inside the adapter with the caller-supplied key.

**Acceptance**
- Replaying the same `claim_id + attempt` returns the cached log row instead of re-calling the provider (deduplication window: 24 hours).
- A failed provider call still produces a log row so auditors can reconstruct the timeline.

**Estimated effort:** 1.5 days.

### Phase 6 — Async provider support (built but inert)

**Scope**
- Extend `Provider` interface: `Verify(...)` may return `Status = StatusPending` with a `ProviderJobID`.
- Add `Provider.Poll(ctx, jobID) (*VerifyResult, error)` method; `VerifyNow` returns `ErrNotSupported`.
- New endpoint `POST /v2/.../verify-bank-account/status/:job_id` for the frontend to poll.
- Optional webhook endpoint `POST /bav/webhook/:provider` with HMAC verification — disabled until a provider needs it.
- Frontend: when `status === 'pending'`, show a spinner and poll every 3s for up to 60s.

**Acceptance**
- Mock async provider (`providers/mock.go` with `MOCK_BAV_ASYNC=true`) exercises the pending → complete flow end-to-end.
- VerifyNow path unchanged — still returns synchronously on the first call.

**Estimated effort:** 2 days.

### Phase 7 — Second real provider (validation of the abstraction)

Pick the first insurer's real provider. Historically in this market the likely next adapter is **LexisNexis** or **TransUnion** — both OAuth2 client-credentials.

**Scope**
- `providers/lexisnexis.go` (or whichever is needed first) implementing `Provider`.
- Shared OAuth2 helper in `providers/common/oauth2.go` with token caching.
- Sandbox-env smoke tests documented in `docs/bav-providers/lexisnexis.md` (credentials, sandbox URL, test IDs).
- Zero changes to controller, frontend, or log table — if any are needed, the abstraction leaked and Phase 1–5 must be revisited.

**Acceptance**
- Flipping `BAV_PROVIDER=lexisnexis` in a sandbox env routes all verifications to LexisNexis with no other code change.
- Response normalisation produces the same `VerifyResult` shape the frontend expects.

**Estimated effort:** 2 days per additional provider.

### Phase 8 — Deployment hardening

**Scope**
- Update the CI/CD workflow (`api/.github/workflows/deploy.yml`) so each deployment target (app1/app2/app3) can be configured with its own `BAV_*` secrets.
- Update operational runbook / README with the new env var matrix.
- Add a startup health check: on boot, `Registry` pings the provider's health endpoint (if one exists) and logs a warning — never a hard failure — if unreachable.

**Acceptance**
- A fresh deployment with no `BAV_*` or `VERIFYNOW_*` vars set logs a clear warning and fails any verification call with `ErrProviderNotConfigured` rather than a panic.

**Estimated effort:** 0.5 day.

## 6. Risks and mitigations

| Risk | Mitigation |
|---|---|
| Async provider contract isn't as clean as sync; interface gets uglier over time. | Phase 6 designs the seam *before* a second provider is added. Re-evaluate after the first async integration. |
| Frontend coupling to VerifyNow wire shape is wider than `GroupPricingService.ts`. | Phase 4 begins with a grep audit (`identityMatch`, `accountFound`, etc.). If call sites > 3, scope grows. |
| Insurer demands a provider whose semantic model truly differs (e.g. no `acceptsCredits` flag). | `TriState` already accommodates `Unknown`. UI shows "unknown" rather than failing. |
| Legacy `VERIFYNOW_*` env vars stay in deployments indefinitely. | Schedule removal: one release after Phase 3, log a deprecation warning on boot. |
| Raw payload logging captures PII (ID numbers). | `bav_verification_log.request_payload` must be redacted at read-time for non-admin users, and the table should be covered by existing data-retention policy. Coordinate with whoever owns data retention before Phase 5 ships. |

## 7. Rollout sequence

1. Ship Phases 1–3 in a single release. VerifyNow deployments unchanged.
2. Ship Phase 4 (canonical frontend contract) the next release. Requires coordinated api+app deploy or use the `/v2` endpoint split.
3. Ship Phase 5 (audit log) as a standalone release — has a migration, deserves its own blast radius.
4. Phase 6 ships when the first async-provider insurer is on the near horizon, not before.
5. Phase 7 ships per insurer, on demand.
6. Phase 8 can land any time after Phase 3.

## 8. Open questions

1. Is there a canonical data-retention policy for BAV request/response payloads under POPIA? (Affects Phase 5 log table design.)
2. Do any existing insurer contracts already specify a BAV provider we should line up as Phase 7?
3. Should `bav_verification_log` live alongside claims (same DB) or in a separate audit schema?
4. `AutoMigrate` will handle table creation across MySQL/Postgres/SQL Server, but JSON column support differs (Postgres `JSONB` vs MySQL `JSON` vs SQL Server `NVARCHAR(MAX)` with CHECK constraints). Confirm `datatypes.JSON` behaves acceptably on SQL Server in this codebase, or fall back to `text` for the two payload columns.

## 9. Effort summary

| Phase | Effort |
|---|---|
| 1. Canonical model | 0.5 day |
| 2. Port VerifyNow | 1 day |
| 3. Registry + controller cutover | 1 day |
| 4. Canonical frontend contract | 1.5 days |
| 5. Audit log + idempotency | 1.5 days |
| 6. Async support | 2 days |
| 7. Second provider (per provider) | 2 days |
| 8. Deployment hardening | 0.5 day |
| **Total (Phases 1–6, 8)** | **8 days** |
