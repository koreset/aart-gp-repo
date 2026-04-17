# Bank Account Verification (BAV)

**Status:** Phases 1–6 shipped. Phase 7 (second real provider) and Phase 8
(deployment hardening) pending.
**Owner:** Backend (`api/services/bav/`) with frontend in
`app/src/renderer/`.
**Companion doc:** [`bav-provider-abstraction-plan.md`](./bav-provider-abstraction-plan.md)
— the implementation plan that drove this work.

This document describes the provider-agnostic BAV module: what it does,
how to operate it, and how to extend it. If you are planning changes,
read the plan first; if you are using or deploying the system, start
here.

## 1. Overview

AART Group Risk integrates with whichever Bank Account Verification
provider the insurer holds a contract with (VerifyNow, LexisNexis,
TransUnion, XDS, ThisIsMe, direct BankservAfrica AVS, etc.). Every
real provider front-ends **BankservAfrica AVS**, so the semantic
contract is standard; providers differ in auth scheme, wire format,
sync vs async, and error shape.

The `services/bav/` package exposes a canonical domain model and a
`Provider` interface. Controllers, services, and the frontend only ever
speak the canonical shape. Adapters per provider live under
`services/bav/providers/`. Swapping providers is a config change:
`BAV_PROVIDER=<name>` plus the provider's credentials.

## 2. Architecture

```
api/
├── services/bav/
│   ├── types.go            VerifyRequest, VerifyResult, TriState, Status
│   ├── provider.go         Provider interface, Registry, sentinel errors
│   ├── registry.go         package-level SetDefault/Active/Verify/Poll
│   ├── logger.go           Logger interface, LogEntry, DedupeWindow (24h)
│   ├── audit/
│   │   └── gorm_logger.go  GORM-backed Logger implementation
│   └── providers/
│       ├── registry.go     NewRegistry(Config) factory
│       ├── verifynow.go    VerifyNow adapter (sync)
│       └── mock.go         Mock adapter (sync or async)
├── models/
│   └── bav_verification_log.go    BAVVerificationLog GORM model
├── controllers/
│   └── claim_payments.go   VerifyBankAccount{,V2,V2Status}, BAVWebhook
├── migrations/{mysql,postgresql,mssql}/
│   └── 20260417080000_create_bav_verification_log.sql
└── main.go                 wires the Registry + GORM logger at startup

app/src/renderer/
├── types/bav.ts            TriState, VerifyResult, triStateIcon()
├── api/GroupPricingService.ts   verifyBankAccount, getBankVerificationStatus
└── screens/group_pricing/claims_management/components/
    └── ClaimRegistrationForm.vue    poll loop (3s × up to 60s)
```

**Layering rules**

- `services/bav` has **zero** dependencies on other internal packages.
  Types + interfaces only.
- `services/bav/providers` imports `services/bav` (to satisfy the
  interface). Per-provider adapters import only what they need.
- `services/bav/audit` imports `api/models` and `gorm.io/gorm`.
- Nothing outside these three packages imports provider-specific code.

## 3. Canonical types reference

All fields exposed below are shared between the Go model and the
TypeScript mirror (`app/src/renderer/types/bav.ts`).

### `VerifyRequest` (input)

| Field | Type | Notes |
|---|---|---|
| `FirstName` | string | |
| `Surname` | string | |
| `IdentityNumber` | string | RSA ID or passport |
| `IdentityType` | string | Defaults to `"IDNumber"` when empty |
| `BankAccountNumber` | string | |
| `BankBranchCode` | string | Universal branch code |
| `BankAccountType` | string | e.g. `"Cheque"`, `"Savings"` |
| `ClaimID` | `*int` | Optional — present when verification is tied to a claim |
| `Attempt` | int | Retry counter, defaults to 1. Used in the idempotency key |
| `IdempotencyKey` | string | Caller-supplied or auto-derived (see §9) |

### `VerifyResult` (output)

| Field | Type | Notes |
|---|---|---|
| `Status` | `Status` | `"complete"` \| `"pending"` \| `"failed"` |
| `Verified` | bool | Provider's overall verdict |
| `Summary` | string | Human-readable |
| `AccountFound` | `TriState` | `"yes"` \| `"no"` \| `"unknown"` |
| `AccountOpen` | `TriState` | |
| `IdentityMatch` | `TriState` | |
| `AccountTypeMatch` | `TriState` | |
| `AcceptsCredits` | `TriState` | |
| `AcceptsDebits` | `TriState` | |
| `Provider` | string | e.g. `"verifynow"` |
| `ProviderRequestID` | string | Provider's own request ID |
| `ProviderJobID` | string | Populated when `Status == "pending"` |
| `ProviderStatusText` | string | Provider's native status string (`json:"-"`) |
| `RawPayload` | []byte | Raw provider response (`json:"-"`) |

`TriState` always resolves to one of three values. Unrecognised provider
input maps to `"unknown"` via `bav.ParseTriState`.

### Sentinel errors

Matched with `errors.Is`. Adapters wrap these with additional context:

| Error | Meaning |
|---|---|
| `ErrUnauthorized` | 401 from the provider — credentials invalid |
| `ErrRateLimited` | 429 — back off and retry later |
| `ErrProviderUnavailable` | 5xx or network error |
| `ErrInvalidInput` | 400 or adapter-level validation failure |
| `ErrProviderNotConfigured` | No provider wired (missing API key, unknown name, or registry absent) |
| `ErrNotSupported` | Sync provider asked to `Poll`, or webhook for provider with no wiring |

## 4. HTTP endpoints

All v2 endpoints are behind the normal auth middleware. The webhook is
intentionally unauthenticated at the route level — auth happens via
HMAC inside the handler.

| Method | Path | Purpose | Since |
|---|---|---|---|
| `POST` | `/group-pricing/claims/verify-bank-account` | Legacy v1. Emits VerifyNow-shaped payload (`"Yes"`/`"No"`/`"Unknown"`). **Scheduled for removal one release after v2 reaches all clients.** | Phase 3 |
| `POST` | `/v2/group-pricing/claims/verify-bank-account` | Canonical v2. Emits `VerifyResult`. | Phase 4 |
| `POST` | `/v2/group-pricing/claims/verify-bank-account/status/:job_id` | Poll a pending async verification. Returns 501 if the active provider is synchronous. | Phase 6 |
| `POST` | `/bav/webhook/:provider` | Inert stub returning 501. Phase 7 wires HMAC + dispatch per provider. | Phase 6 |

### v2 request body

```json
{
  "first_name": "Thandi",
  "surname": "Nkosi",
  "identity_number": "9001015009087",
  "identity_type": "IDNumber",
  "bank_account_number": "1234567890",
  "bank_branch_code": "250655",
  "bank_account_type": "Cheque",
  "claim_id": null,
  "attempt": 1
}
```

`claim_id` and `attempt` are optional; leaving them out is fine for
pre-claim verification. See §9 for how they affect idempotency.

### v2 response body

```json
{
  "success": true,
  "data": {
    "status": "complete",
    "verified": true,
    "summary": "All checks passed",
    "accountFound": "yes",
    "accountOpen": "yes",
    "identityMatch": "yes",
    "accountTypeMatch": "yes",
    "acceptsCredits": "yes",
    "acceptsDebits": "no",
    "provider": "verifynow",
    "providerRequestId": "req-abc-123"
  }
}
```

When `status === "pending"`, `providerJobId` is populated and the caller
should poll the status endpoint until `status` flips to `"complete"` or
`"failed"`.

### Error responses

| HTTP | When |
|---|---|
| 400 | Request binding failed (missing required fields) |
| 404 | Status endpoint called with an unknown `job_id` |
| 500 | Any adapter error other than `ErrNotSupported` / `ErrInvalidInput` |
| 501 | Status endpoint called against a sync provider; webhook endpoint |

## 5. Configuration reference

All config is via environment variables, read in `api/config/config.go`.

| Var | Default | Purpose |
|---|---|---|
| `BAV_PROVIDER` | `verifynow` | Short name of the active adapter. Currently accepted: `verifynow`, `mock`. |
| `BAV_API_KEY` | *empty* | API key for the active provider. Falls back to `VERIFYNOW_API_KEY` if empty (backward-compat). |
| `BAV_BASE_URL` | *provider default* | Override the provider's base URL. Adapters supply their own production default. |
| `BAV_OAUTH_CLIENT_ID` | *empty* | For OAuth2 providers (reserved for Phase 7). |
| `BAV_OAUTH_CLIENT_SECRET` | *empty* | |
| `BAV_OAUTH_TOKEN_URL` | *empty* | |
| `BAV_TIMEOUT_SECONDS` | `45` | Per-request HTTP timeout. |
| `MOCK_BAV_ASYNC` | *unset* | When `true`, the mock provider returns `pending` and resolves via `Poll` after ~3s. Only consulted when `BAV_PROVIDER=mock`. |
| `VERIFYNOW_API_KEY` | *empty* | Legacy. Still honoured: populates `BAV_API_KEY` when that is empty. |
| `VERIFYNOW_MODE` | `production` | VerifyNow-specific mode passed through to the provider. |

**Backward-compat guarantees**

Deployments that previously set only `VERIFYNOW_API_KEY` continue to
work with no config changes: the default `BAV_PROVIDER` is `verifynow`
and `BAV_API_KEY` falls back to `VERIFYNOW_API_KEY`.

## 6. Using BAV from Go code

```go
import "api/services/bav"

result, err := bav.Verify(ctx, bav.VerifyRequest{
    FirstName:         "Thandi",
    Surname:           "Nkosi",
    IdentityNumber:    "9001015009087",
    BankAccountNumber: "1234567890",
    BankBranchCode:    "250655",
    BankAccountType:   "Cheque",
    ClaimID:           &claimID, // optional
    Attempt:           1,         // optional, defaults to 1
})
if err != nil {
    if errors.Is(err, bav.ErrProviderNotConfigured) {
        // no adapter wired; fall back or surface to operator
    }
    return err
}

if result.Status == bav.StatusPending {
    // Poll with result.ProviderJobID
    final, err := bav.Poll(ctx, result.ProviderJobID)
    ...
}
```

- `bav.Verify` and `bav.Poll` route through the process-wide `Registry`
  installed in `main.go` during startup.
- Never talk to providers directly; always go through the package. The
  registry handles idempotency key derivation, 24h success dedup, and
  audit logging uniformly.
- Logger errors are swallowed inside the registry — a persistence
  failure will not break the user-facing call. Check the app logs for
  `bav/audit: failed to persist verification log:` lines if audit rows
  are missing.

## 7. Using BAV from the frontend

Import types and helpers:

```ts
import {
  triStateIcon,
  type VerifyResult,
  type TriState,
} from '@/renderer/types/bav'
import GroupPricingService from '@/renderer/api/GroupPricingService'
```

Call and handle sync + async uniformly:

```ts
const res = await GroupPricingService.verifyBankAccount({
  first_name: firstName,
  surname,
  identity_number: idNumber,
  bank_account_number,
  bank_branch_code,
  bank_account_type,
})
let result = res.data.data

if (result.status === 'pending' && result.providerJobId) {
  // Poll every 3s up to 60s total — see pollUntilResolved in
  // ClaimRegistrationForm.vue for the reference implementation.
}

if (result.verified) {
  // success UI
} else if (result.status === 'failed') {
  // error UI; use granular TriState fields to explain which checks
  // failed. 'unknown' is distinct from 'no' — bank could not confirm.
}
```

The `triStateIcon()` helper returns `{ icon, color, label }` for Vuetify
`mdi-check-circle` / `mdi-close-circle` / `mdi-help-circle` rendering.
Use it when you want each granular check's outcome visible; the
reference component shows them as a row of `v-chip`s.

## 8. Adding a new provider

End-to-end recipe. Takes roughly half a day per provider, more if the
wire format is exotic.

**Step 1 — Adapter file.** Create
`api/services/bav/providers/<name>.go`:

```go
package providers

import "api/services/bav"

type MyProvider struct { /* config, HTTP client */ }

func NewMyProvider(cfg MyProviderConfig) *MyProvider { ... }

func (a *MyProvider) Name() string { return "myprovider" }

func (a *MyProvider) Verify(ctx context.Context, req bav.VerifyRequest) (*bav.VerifyResult, error) {
    // 1. Marshal req to provider's wire format
    // 2. Dispatch HTTP call (use ctx, honour req.IdempotencyKey as the provider's idempotency header)
    // 3. Map HTTP status to bav.Err* sentinels (401→ErrUnauthorized, 429→ErrRateLimited, 5xx→ErrProviderUnavailable, 400→ErrInvalidInput)
    // 4. Unmarshal response; translate provider's "Y"/"yes"/true/etc. via bav.ParseTriState
    // 5. Return *bav.VerifyResult with Status=Complete (sync) or Pending (async); set RawPayload
}

func (a *MyProvider) Poll(ctx context.Context, jobID string) (*bav.VerifyResult, error) {
    // For sync providers:
    return nil, bav.ErrNotSupported
    // For async providers, hit the provider's status endpoint and map the result.
}

var _ bav.Provider = (*MyProvider)(nil)
```

**Step 2 — Register in the factory.** Add a case to
`api/services/bav/providers/registry.go::NewRegistry`:

```go
case "myprovider":
    return bav.NewRegistry(NewMyProvider(MyProviderConfig{
        APIKey:  cfg.APIKey,
        BaseURL: cfg.BaseURL,
        // map OAuth2 fields if the provider uses OAuth
    })), nil
```

**Step 3 — Tests.** Create `<name>_test.go` using `httptest.Server` to
cover: happy path, 401, 429, 5xx, 400, malformed JSON, context
cancellation. Target ≥80% coverage on the adapter (existing VerifyNow
adapter is at 93%).

**Step 4 — Webhook handler** (async providers only). Wire the provider
into `controllers.BAVWebhook`, replacing the current 501 stub with
HMAC verification per `cfg.OAuthClientSecret` (or a dedicated webhook
secret) and dispatch into the adapter.

**Step 5 — Docs.** Add a subsection under this file noting: sandbox
URL, required env vars, test IDs, known quirks.

**Step 6 — Operations.** Ensure the deployment's env has the new
`BAV_PROVIDER=<name>` and `BAV_API_KEY` (or OAuth creds). Nothing in
controllers, frontend, or log table changes. **If any of them need to
change, the abstraction has leaked — revisit Phases 1–5.**

## 9. Audit log & idempotency

### Schema

Table: `bav_verification_logs`

| Column | Type | Notes |
|---|---|---|
| `id` | INT PK | |
| `claim_id` | INT NULL | Nullable — verification may precede claim creation |
| `provider` | VARCHAR(64) | Indexed |
| `provider_request_id` | VARCHAR(128) | |
| `idempotency_key` | VARCHAR(128) | Indexed (composite with `created_at`) |
| `status` | VARCHAR(32) | `complete` \| `pending` \| `failed` |
| `request_payload` | TEXT (NVARCHAR(MAX) on SQL Server) | Canonical `VerifyRequest` JSON |
| `response_payload` | TEXT | Canonical `VerifyResult` + raw provider payload, reconstitutable via `audit.unmarshalStoredResult` |
| `error_message` | VARCHAR(1024) | Truncated if longer |
| `created_at` | DATETIME / TIMESTAMP | Indexed |

Stored as plain text (via the existing `models.JSON` type) rather than
`datatypes.JSON` — this sidesteps the SQL Server compatibility
question still open in the plan.

### Idempotency key

When a caller does not supply `VerifyRequest.IdempotencyKey`, the
registry derives it:

```
sha256(claim_id | attempt | provider | identity_number | bank_account_number | bank_branch_code)
```

Stable across retries as long as the banking fields don't change. If
the user edits account or branch and retries, the key changes — new
call, fresh audit row.

### Dedup window (24h)

Before calling the adapter, the registry looks up
`bav_verification_logs` for any row where:

- `idempotency_key` matches, **and**
- `status = 'complete'` (only successes are cached), **and**
- `created_at >= now - 24h`

On hit, the cached `VerifyResult` is reconstituted and returned —
adapter is not called. On miss, the adapter runs and a new log row is
written regardless of success or failure.

**Failures are never cached.** A transient provider error must not
block the user's next retry.

### Logger errors

`GormLogger.Record` surfaces storage errors to its caller, but the
`Registry` swallows them (`_ = r.logger.Record(...)`). A DB outage will
not break verification. Monitor for `bav/audit: failed to persist
verification log:` lines in app logs.

## 10. Async (pending/poll) flow

```
┌──────────┐  Verify       ┌──────────┐
│ Frontend │──────────────▶│ Backend  │
│          │  status=      │          │
│          │◀───pending────│ Registry │
│          │  jobId=J      └──────────┘
│          │
│ wait 3s  │
│          │  /status/J    ┌──────────┐
│          │──────────────▶│ Backend  │
│          │  status=      │          │
│          │◀───pending────│ Registry │
│          │               └──────────┘
│ wait 3s  │
│          │  /status/J    ┌──────────┐
│          │──────────────▶│          │
│          │  status=      │ Registry │
│          │◀───complete───│          │
└──────────┘  verified=Y   └──────────┘
```

- **Frontend loop**: every 3s, cap 60s total
  (`ClaimRegistrationForm.vue::pollUntilResolved`). On timeout, the UI
  shows a "still in progress — retry in a few minutes" message.
- **Backend**: `bav.Poll(ctx, jobID)` routes through `Registry.Poll`.
  Only **terminal** polls (complete, failed, or transport error) write
  an audit row; intermediate `StatusPending` polls are skipped so a
  single 60s verification produces ~2 rows (initial pending + final
  resolution), not ~20.
- **Sync providers** (VerifyNow) return `ErrNotSupported` from `Poll`,
  which the status endpoint maps to HTTP 501.

## 11. Testing & local development

### Unit tests

```bash
cd api
go test ./services/bav/...
# Coverage: bav 83%, providers 93%, audit 26% (DB paths not tested)
```

### End-to-end with the mock provider

```bash
# API side — async mock returns pending, resolves in ~3s:
export BAV_PROVIDER=mock
export MOCK_BAV_ASYNC=true
air  # or go run .

# App side — unchanged:
cd app && yarn dev
```

Open the Claim Registration form, fill banking details, click **Verify
Account**. Expected:

1. "Verification in progress..." appears immediately.
2. ~3s later the success alert + green tick chips appear for all
   granular checks.

Switch `MOCK_BAV_ASYNC=false` (or unset) to exercise the sync path with
the same mock adapter.

### Frontend type check

```bash
cd app
npx vue-tsc --noEmit
```

## 12. Known open questions (deferred)

Carried over from the plan. Each will decide final design before Phase
5's log table is touched at scale, or before Phase 7 ships a second
real provider.

1. **POPIA data-retention policy** for `bav_verification_logs`. ID
   numbers and account numbers live in `request_payload`. Need a
   read-time redaction strategy (non-admin users see masked values) or
   a scheduled purge. Currently: raw data retained indefinitely.
2. **Existing insurer contracts.** Which provider should Phase 7
   target first — LexisNexis, TransUnion, or something else? Unknown
   until procurement shares the contract list.
3. **Log table location.** Alongside claims (same DB) or separate
   audit schema? Currently: same DB as claims.
4. **`datatypes.JSON` on SQL Server.** Not validated in this codebase.
   Currently working around by storing payloads as plain text via
   `models.JSON`. Upgrade is a one-line model change + migration if
   JSON column performance matters.
5. **Polling row volume.** Resolved: `Registry.Poll` skips audit writes
   for intermediate `StatusPending` responses and only logs terminal
   transitions. A 60s verification now produces ~2 rows (initial
   pending from `Verify` + final resolution) rather than ~20. Trade-off
   accepted: the audit log no longer captures individual poll
   timestamps — only the start and end of each verification. If
   per-poll timing becomes necessary, revisit with an
   UPDATE-in-place-with-`poll_count` approach.

## 13. Phase status

| Phase | Delivery | Status |
|---|---|---|
| 1. Canonical model | `types.go`, `provider.go`, `ParseTriState` | ✅ |
| 2. Port VerifyNow | `providers/verifynow.go` + tests | ✅ |
| 3. Registry + controller cutover | Config plumbing, v1 legacy wire-shape preserved | ✅ |
| 4. Canonical frontend contract | `/v2` endpoint, `types/bav.ts`, per-check chip row | ✅ |
| 5. Audit log + idempotency | `bav_verification_log`, 24h dedup, GORM logger | ✅ |
| 6. Async support | `Poll`, mock provider, status endpoint, poll loop, webhook stub | ✅ |
| 7. Second real provider | TBD — driven by insurer requirement | ⏳ |
| 8. Deployment hardening | Per-env secrets in CI, runbook, startup health check | ⏳ |

## 14. References

- Implementation plan: [`bav-provider-abstraction-plan.md`](./bav-provider-abstraction-plan.md)
- VerifyNow adapter: `api/services/bav/providers/verifynow.go`
- Mock adapter: `api/services/bav/providers/mock.go`
- Reference frontend component:
  `app/src/renderer/screens/group_pricing/claims_management/components/ClaimRegistrationForm.vue`
