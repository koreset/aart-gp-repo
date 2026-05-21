# Email Setup & Operations Guide

This guide covers everything needed to make outbound email work in AART Group
Risk — for the **server admin** who prepares the environment and for the
**application admin** who fills in the settings screen. Read §2 (server
prerequisites) before configuring anything in the UI.

> TL;DR: set `APP_SECRET` on the server (same value on every instance) → apply
> the DB migration → seed the email templates → pick a provider in
> **Group Pricing → Email → Settings** → **Send Test Email** → watch the
> **Outbox**.

---

## 1. How email works

Email is **queue-based**, not send-on-click:

```
Settings (per license)        Template (per license)
        │                            │
        └──────────► EnqueueEmail ◄──┘
                          │
                   email_outbox row  (status = pending)
                          │
              background outbox worker (every 4–7s)
                          │
              BuildMailer(settings)  ── provider? ──┐
                          │                         │
                  ┌───────┴────────┐        ┌───────┴────────┐
                  │  SMTPMailer    │        │  GraphMailer   │
                  │ (relay / SMTP) │        │ (Microsoft 365)│
                  └───────┬────────┘        └───────┬────────┘
                          └─────────► sent / failed ◄┘
```

Key facts:

- **One config row per license.** A license is either SMTP **or** Microsoft 365,
  never both at once.
- **A background worker** polls `email_outbox` and sends. Clicking *Save* or
  *Send Test Email* only **queues** — delivery happens seconds later. Always
  confirm the outcome in the **Outbox**, not the toast message.
- **Secrets are encrypted at rest** (SMTP password, Graph client secret) with
  `APP_SECRET`. They are never returned by the API or written to logs.
- **Multi-instance safe.** Each row is claimed with an atomic
  `UPDATE … WHERE status = 'pending'`, so running multiple API instances
  against the same database never double-sends.
- **Failed sends retry** with exponential backoff: 1 → 2 → 4 → 8 → 16 minutes,
  up to `max_attempts` (default 5), then the row is marked `failed`. Some
  failures are **terminal** (no retry) — see §6.

Relevant source: `services/email_worker.go`, `services/email_mailer.go`,
`services/email/smtp.go`, `services/email/graph.go`,
`services/email_settings.go`, `models/email.go`.

---

## 2. Server prerequisites (do this first)

### 2.1 `APP_SECRET` — **required**

`APP_SECRET` is the key used to AES-256-GCM encrypt the SMTP password and the
Graph client secret before they are stored in the database.

- Set it to a long random string (e.g. `openssl rand -hex 32`).
- It is read at startup from the environment / the instance's `.env`
  (`config/config.go` calls `godotenv.Load()`).

> **Gotcha — missing.** If `APP_SECRET` is unset, saving or sending fails with
> `crypto: APP_SECRET environment variable is not set`. There is no fallback by
> design.

> **Gotcha — must be identical on every instance.** Production runs **three**
> services (`app1`, `app2`, `app3` under `/home/aart/api1|api2|api3`), each
> loading its **own** `.env`. The deploy workflow does **not** set `APP_SECRET`.
> If the value differs between instances, a secret encrypted by one instance
> cannot be decrypted by another, so sends fail **intermittently** depending on
> which instance's worker happens to claim the row. Put the **same**
> `APP_SECRET` in `api1/.env`, `api2/.env`, and `api3/.env` (or a shared
> systemd `EnvironmentFile`).

> **Gotcha — rotation.** Changing `APP_SECRET` makes every previously stored
> secret undecryptable. After a rotation, every license must re-enter its SMTP
> password / Graph client secret. Treat it as a long-lived secret.

### 2.2 Database migration

The provider columns are added by
`migrations/mysql/20260603000000_add_email_graph_provider.sql`. It is applied
automatically on startup by `RunMigrationsOnStartup`.

- Ensure the migration file reaches the server and the API is **restarted** so
  the runner applies it. (The deploy workflow copies `migrations/*` to
  `api1`; the shared `migrations` table prevents re-application across
  instances.)
- Existing rows default to `provider = 'smtp'`, so current configs keep working.

> **Gotcha — non-MySQL deployments.** The email tables currently exist only in
> `migrations/mysql/`. A PostgreSQL or SQL Server deployment has no email tables
> at all yet — the base email migrations must be ported to that dialect before
> any email feature (including this one) will work there.

> **Gotcha — DB not configured.** If `DB_HOST`/`DB_USER`/`DB_PWD`/`DB_NAME` are
> blank the API can't connect and migrations never run. Confirm the DB env vars
> are filled before expecting email to work.

### 2.3 Seed the email templates

Emails render from templates; nothing sends without one. In particular the
**Send Test Email** button uses a template with code **`system_test`**, which
must exist and be **active**.

- Open **Group Pricing → Email → Templates** and seed the starter templates
  (the screen offers this), or create `system_test` manually.

> **Gotcha.** If the test button returns *"is a template with code system_test
> configured?"*, the template is missing or still in `draft` — seed/activate it.

### 2.4 Network egress

The API server makes the outbound connections, not the desktop app. The
server's firewall/proxy must allow:

- **SMTP provider:** outbound TCP **587** (STARTTLS) or **465** (implicit TLS)
  to the relay host. Corporate networks frequently block SMTP egress.
- **Microsoft 365:** outbound **HTTPS (443)** to `login.microsoftonline.com`
  and `graph.microsoft.com`.

> **Gotcha — TLS-inspecting proxies.** A proxy that re-signs TLS can break
> certificate validation to Microsoft endpoints. Allow-list these hosts to
> bypass inspection if token/sendMail calls fail with TLS errors.

---

## 3. Choosing a provider

| | **SMTP relay** | **Microsoft 365 (Graph)** |
|---|---|---|
| Use when | Client uses (or accepts) a transactional relay | Client must send from their Microsoft 365 / Exchange Online mailbox |
| Examples | SendGrid, Postmark, Brevo, Mailgun, Amazon SES | A licensed or shared mailbox in the client's M365 tenant |
| Auth | Username + password (or API key) | Azure app — OAuth2 app-only |
| Setup effort | Low (sign up, verify sender) | Medium (Azure app registration + admin consent) |
| Best for bulk | Yes (relays are built for it) | No — M365 throttles; use a relay for volume |

> **Gotcha — do not use Office 365 / Outlook.com over SMTP.** Microsoft retired
> basic SMTP auth. An Outlook/M365 mailbox configured as an *SMTP* provider
> returns `504 5.7.4 Unrecognized authentication type` no matter what password
> (or app password) you use. For Microsoft mailboxes you **must** use the
> Microsoft 365 (Graph) provider — see §5.

---

## 4. Provider A — SMTP relay

### Steps (application admin)

1. **Group Pricing → Email → Settings**, set **Provider = SMTP relay**.
2. Fill the fields (table below), **Save**, then **Send Test Email**.
3. Confirm `sent` in the **Outbox**.

### Field reference

| Field | Required | Notes |
|---|---|---|
| Host | ✅ | Relay SMTP host |
| Port | ✅ (default 587) | 587 = STARTTLS, 465 = implicit TLS |
| TLS mode | default `starttls` | `starttls` / `tls` / `none` |
| Auth user | — | Login / API-key username |
| Auth password | — | Encrypted at rest; **leave blank on re-save to keep the stored one** |
| From address | ✅ | Must be a verified sender at the relay |
| From name | — | Display name |
| Reply-to | — | Optional |

### Quick config by relay

All use port **587** + **STARTTLS**:

| Provider | Host | Auth user | Auth password |
|---|---|---|---|
| Brevo | `smtp-relay.brevo.com` | account login email | SMTP **key** (not login password) |
| SendGrid | `smtp.sendgrid.net` | the literal `apikey` | the API key |
| Postmark | `smtp.postmarkapp.com` | server token | same server token |
| Mailgun | `smtp.mailgun.org` | `postmaster@<domain>` | SMTP password |
| Amazon SES | `email-smtp.<region>.amazonaws.com` | SES SMTP username | SES SMTP password |

### SMTP gotchas

- **Sender/domain not verified.** Every relay rejects mail from an unverified
  From address (`550`, "sender not allowed"). Verify the single sender **or** the
  domain in the relay dashboard before sending. You can keep a client's own
  From address (e.g. their Outlook address) by verifying it as a single sender.
- **`tls_mode = none` + auth fails.** Go's SMTP client refuses to send
  credentials over an unencrypted connection (error mentions an *unencrypted
  connection*). Only use `none` for an internal, no-auth relay; otherwise use
  `starttls`/`tls`.
- **Wrong key vs password.** For API-key relays (Brevo SMTP key, SendGrid
  `apikey`, Postmark token), the *password* field is the key/token — not the
  dashboard login password. Wrong value → `535` authentication failure.
- **Firewall blocks egress** → connection timeout / `dial tcp` error. See §2.4.
- **Deliverability (lands in spam).** Configure **SPF**, **DKIM**, and
  **DMARC** for the sending domain at the relay. Without them, mail is filtered
  even though the Outbox shows `sent`.
- **Rate limits.** Free tiers throttle (e.g. Brevo ~300/day). Bulk sending from
  a personal mailbox or free tier will be throttled or blocked — size the plan
  to the volume.

---

## 5. Provider B — Microsoft 365 (Graph)

The app authenticates as an **Azure app (app-only / client-credentials)** and
calls Graph `POST /users/{from}/sendMail`. This sends unattended from the
background worker — no per-user login.

### 5.1 Register the Azure app (tenant admin)

1. **Microsoft Entra admin center → App registrations → New registration.**
   Give it a name (e.g. *AART Mailer*); single-tenant is fine for the
   per-client model.
2. **API permissions → Add a permission → Microsoft Graph → Application
   permissions → `Mail.Send`.** Then **Grant admin consent**.
   - It **must** be an **Application** permission, not Delegated.
   - **Grant admin consent** is mandatory and requires a tenant admin.
3. **Certificates & secrets → New client secret.** Copy the **Value**
   immediately.
4. **Overview** → copy **Directory (tenant) ID** and **Application (client) ID**.

### 5.2 Restrict which mailboxes the app can send as — **important**

App-only `Mail.Send` lets the app send as **any** mailbox in the tenant. Lock it
down with an **Application Access Policy** (Exchange Online PowerShell):

```powershell
# Group containing only the mailbox(es) the app may send from
New-DistributionGroup -Name "AART Mail Senders" -Type Security `
  -PrimarySmtpAddress aart-senders@contoso.com
Add-DistributionGroupMember -Identity "AART Mail Senders" `
  -Member bordereaux@contoso.com

New-ApplicationAccessPolicy -AppId <client-id> `
  -PolicyScopeGroupId aart-senders@contoso.com `
  -AccessRight RestrictAccess `
  -Description "Restrict AART app to the bordereaux mailbox"

# Verify
Test-ApplicationAccessPolicy -Identity bordereaux@contoso.com -AppId <client-id>
```

### 5.3 Configure in the app

**Provider = Microsoft 365 (Graph API)**, then:

| Field | Required | Notes |
|---|---|---|
| Tenant ID | ✅ | Directory (tenant) ID |
| Client ID | ✅* | Application (client) ID (*optional only if a central app is configured — see 5.4) |
| Client Secret | ✅* | The secret **Value**; encrypted at rest; leave blank on re-save to keep |
| From address | ✅ | A real **licensed or shared mailbox** in the tenant, allowed by the access policy |
| From name | — | Display name |
| Reply-to | — | Optional |

### 5.4 Credential ownership — two models

- **Each client registers their own app (recommended).** Fill Tenant ID +
  Client ID + Client Secret per license. Best isolation; the client controls and
  rotates their own secret.
- **One central multi-tenant app (you operate it).** Set `GRAPH_CLIENT_ID` and
  `GRAPH_CLIENT_SECRET` in the **server** environment; each license then records
  only its **Tenant ID** (leave Client ID / Secret blank in the UI). The app
  falls back to the global credentials. Each client admin still grants admin
  consent to your app in their tenant. (Set these env vars on **all** instances,
  like `APP_SECRET`.)

### Microsoft 365 gotchas

- **No admin consent** → token lacks `Mail.Send` → `403` on send. Re-grant
  consent and confirm a green check on the permission.
- **Client secret expiry.** Azure secrets expire (max 24 months). When expired,
  sends fail with a token error (`AADSTS7000222` / `invalid_client`). Set a
  rotation reminder; update the secret in the UI (or server env for the central
  app) before it lapses.
- **From address isn't a real mailbox.** Sending as an address with no mailbox
  (or a distribution list) errors (`ErrorInvalidUser` /
  `MailboxNotEnabledForRESTAPI`). Use a licensed user mailbox or a **shared
  mailbox** (shared mailboxes need no license — common for `no-reply@`).
- **Access policy excludes the sender** → `403 ErrorAccessDenied`. Add the
  mailbox to the policy's group (5.2).
- **`saveToSentItems` is on.** Sent copies land in the mailbox's *Sent Items*
  (uses mailbox storage). This is intentional for auditability.
- **Attachments > ~3 MB** are rejected with a clear error — Graph caps a single
  inline `sendMail` request. (Large-attachment upload sessions are not
  implemented.)
- **Egress / TLS inspection** — see §2.4.

---

## 6. Testing & verifying (the Outbox)

After saving settings, **Send Test Email** queues a `system_test` message to the
signed-in admin's address. Then open **Group Pricing → Email → Outbox**:

| Status | Meaning |
|---|---|
| `pending` | Queued, awaiting the worker (or waiting on a retry backoff) |
| `sending` | Claimed by a worker, in flight |
| `sent` | Accepted by the provider (relay 250 / Graph 202) |
| `failed` | Gave up after `max_attempts`, or a terminal error |

Each row shows **Attempts (n/max)** and **Last error**. Use **Retry** to requeue
a `failed`/`pending` row.

**Terminal failures (no retry):** no email settings for the license, secret
**decrypt** failure (usually a wrong/changed `APP_SECRET`), provider build
errors, and attachment-resolution errors. Everything else (auth, connection,
transient provider errors) retries on the backoff schedule.

> `sent` means the **provider accepted** the message — not that it reached the
> inbox. If `sent` but not received, investigate spam/SPF/DKIM/DMARC (SMTP) or
> the recipient side, not this app.

---

## 7. Troubleshooting matrix

| Symptom (in Last error) | Likely cause | Fix |
|---|---|---|
| `crypto: APP_SECRET environment variable is not set` | `APP_SECRET` missing | Set it (§2.1) and restart |
| `decrypt SMTP password` / `decrypt graph client secret` | `APP_SECRET` changed or differs across instances | Re-enter the secret; make `APP_SECRET` identical on all instances (§2.1) |
| Test button: *"template with code system_test configured?"* | `system_test` missing/draft | Seed/activate templates (§2.3) |
| `504 5.7.4 Unrecognized authentication type` | Office 365 / Outlook over **SMTP** (basic auth retired) | Use the **Microsoft 365 (Graph)** provider (§5) |
| `535` authentication failed | Wrong SMTP user/password/key, or relay needs an app password | Re-check credentials; for API-key relays the password is the key (§4) |
| `550` / "sender not allowed" | From address not verified at relay | Verify sender/domain (§4) |
| `unencrypted connection` error | `tls_mode = none` with auth | Use `starttls`/`tls` (§4) |
| `dial tcp … timeout` | Firewall blocks SMTP egress | Open 587/465 outbound (§2.4) |
| Graph `401` / `AADSTS7000215` | Invalid client secret | Re-enter the secret value |
| Graph token `AADSTS7000222` / `invalid_client` | Client secret expired | Rotate the Azure secret, update settings (§5) |
| Graph token `AADSTS700016` | App not found in tenant | Wrong Client ID or wrong Tenant ID |
| Graph token `AADSTS90002` | Tenant not found | Wrong Tenant ID |
| Graph `403 ErrorAccessDenied` | Missing `Mail.Send` consent, or access policy excludes the mailbox | Grant admin consent (§5.1); add mailbox to policy (§5.2) |
| Graph `ErrorInvalidUser` / `MailboxNotEnabledForRESTAPI` | From address isn't a real mailbox | Use a licensed/shared mailbox (§5) |
| `attachments total … exceed inline limit` | Attachment > ~3 MB via Graph | Reduce size / use a relay for large files |
| Settings won't save: *host is required* / *graph_tenant_id is required* | Required field for the selected provider is blank | Fill the field for that provider (§4/§5) |

---

## 8. Multi-instance deployment notes

Production = three systemd services (`app1`/`app2`/`app3`) on one host, each in
its own `/home/aart/apiN` directory with its own `.env`, all pointing at the
**same database**.

- **`APP_SECRET` must match on all three** (§2.1). This is the single most
  common email deployment failure.
- **`GRAPH_CLIENT_ID` / `GRAPH_CLIENT_SECRET`** (central-app model only) must
  also be set identically on all three.
- All three run an outbox worker against the shared DB; the atomic claim
  prevents double-sends, and polling is jittered (4–7s) so they don't lock-step.
- Migrations apply against the shared DB on startup; ensure the new SQL file is
  on the server and the services were restarted after deploy.

---

## 9. Security notes

- SMTP passwords and Graph client secrets are stored **encrypted**
  (AES-256-GCM via `APP_SECRET`); the API returns only `has_password` /
  `has_graph_secret` booleans, never the secret.
- For Microsoft 365, always scope the app with an **Application Access Policy**
  (§5.2) so a compromised secret can't send as arbitrary mailboxes.
- Prefer **shared mailboxes** for system senders (`no-reply@`) — no interactive
  login, no license cost, easy to lock down.

---

## 10. Reference

**Server environment variables**

| Var | Required | Purpose |
|---|---|---|
| `APP_SECRET` | ✅ | Encrypts stored email secrets (must match across instances) |
| `GRAPH_CLIENT_ID` | optional | Central multi-tenant Azure app client id (fallback) |
| `GRAPH_CLIENT_SECRET` | optional | Central app client secret (fallback) |
| `DB_HOST`/`DB_USER`/`DB_PORT`/`DB_PWD`/`DB_NAME` | ✅ | Database connection |

**Permissions (per the routes in `routes/routes.go`)**

| Action | Permission |
|---|---|
| View settings | (none — license header only) |
| Save settings / send test | `email:configure` |
| Manage templates | `email:templates:manage` |
| View / retry outbox | `email:outbox:view` |

**Endpoints** (under `/group-pricing/email`, `X-License-Id` header required):
`GET|PUT settings`, `POST settings/test`, `GET|POST templates`,
`GET|PUT|DELETE templates/:code`, `POST templates/:code/preview`,
`GET outbox`, `GET outbox/:id`, `POST outbox/:id/retry`.

**Source map**

| Concern | File |
|---|---|
| Settings model + constants | `models/email.go` |
| Save/validate settings | `services/email_settings.go` |
| Provider factory | `services/email_mailer.go` |
| SMTP transport | `services/email/smtp.go` |
| Microsoft Graph transport | `services/email/graph.go` |
| Outbox worker / retry | `services/email_worker.go` |
| Settings UI | `app/src/renderer/screens/group_pricing/email/EmailSettings.vue` |
| Migration | `migrations/mysql/20260603000000_add_email_graph_provider.sql` |
