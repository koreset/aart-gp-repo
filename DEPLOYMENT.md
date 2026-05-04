# AART Group Risk API — Production Deployment Plan (Ubuntu + Nginx, 2 instances)

This is a step-by-step runbook for deploying the AART Group Risk Go API to a single Ubuntu
host running **two API instances** behind an **Nginx reverse proxy** with **Let's Encrypt
TLS**. The MySQL database lives on the **same shared MySQL server already used by
aart-valuations** — but in a new database (`aart_grouprisk_db`).

The runtime and deployment user is the **`moruogrsadmin`** account on the Ubuntu host.
That user owns the binaries, runs the API processes, and is the SSH account used to
deploy new builds.

This runbook mirrors the structure of the aart-valuations deployment runbook so the muscle
memory carries over. The differences worth flagging up-front:

- Runtime user is **`moruogrsadmin`**, not `moruoadmin`.
- Install root is **`/opt/aart-gr`**, not `/opt/aart`.
- systemd unit is **`aart-gr-api@.service`**, not `aart-api@.service`.
- Nginx config is **`aart-gr-api.conf`** with upstream **`aart_gr_backend`**.
- DB schema is **`aart_grouprisk_db`** on the shared MySQL host. The `aart_app` user
  already exists on that host (scoped to the aart-valuations API IP). MySQL keys users
  by `(user, host)`, so the GRANT in §5 creates a *new* user row scoped to this host's
  IP — not a conflict.

It is written to be executed top-to-bottom. Where a value depends on your environment it
is shown as `<PLACEHOLDER>` — find-and-replace before you run.

---

## 1. Target topology

```
                    ┌──────────────────────────┐
       Internet ─►  │  Nginx (443/80)          │
                    │  TLS termination         │
                    │  upstream: aart_gr_backend
                    └─────────┬─────────┬──────┘
                              │         │
                       ┌──────▼──┐  ┌───▼─────┐
                       │ aart-gr-│  │ aart-gr-│        same host
                       │ api:9091│  │ api:9092│   (run as moruogrsadmin)
                       └────┬────┘  └────┬────┘
                            │            │
                            └─────┬──────┘
                                  │
                       ┌──────────▼──────────┐
                       │  Redis (localhost)  │   WebSocket pub/sub fan-out
                       └─────────────────────┘   between the two instances
                                  │
                       ┌──────────▼──────────┐
                       │  MySQL (shared DB)  │   private network
                       │  aart_grouprisk_db  │
                       └─────────────────────┘
```

Why these choices, given the AART Group Risk codebase:

- **Two instances on different ports**, not different hosts: `api/main.go:345-350` reads
  `APP_PORT` env var with config.json as fallback. One binary, one config file, two
  systemd instance units differing only in port.
- **Redis is required, not optional.** `services.StartRedisWSSubscriber()`
  (`api/services/ws_hub.go:330-346`, called from `api/main.go:312`) is what lets
  WebSocket events emitted on instance A reach clients connected to instance B. With two
  instances and no Redis, real-time updates will silently drop for half your sessions.
- **Nginx terminates TLS, not the Go app.** The binary's built-in `autotls` (Let's
  Encrypt) is gated behind `ENVIRONMENT=production` (`api/main.go:367-374`); **leave
  that env var unset** or the Go app will fight Nginx for ports 80/443.
- **One instance owns each projection job.** Per the project notes, "only one instance
  can claim a queued job" — the worker pool is safe across two instances. No special
  coordination needed.

---

## 2. Pre-deployment checklist

Have these ready before you start. Values already known are filled in; the rest are
placeholders.

| Item | Value |
| --- | --- |
| Public domain | `moruo-gr-api.aart-enterprise.com` |
| Ubuntu server public IP | `102.133.230.104` |
| Ubuntu server private IP (if any) | `<API_HOST_PRIV_IP>` |
| MySQL host | `<DB_HOST>` (same host as aart-valuations) |
| MySQL port | `3306` (or custom) |
| MySQL admin credentials (one-time, for provisioning) | — |
| MySQL app user / password | `aart_app` / `<DB_PASSWORD>` |
| MySQL database name | `aart_grouprisk_db` |
| Redis password | `<REDIS_PASSWORD>` (generate: `openssl rand -hex 24`) |
| Notification email for certbot | `<OPS_EMAIL>` |
| SSH / runtime user on the host | `moruogrsadmin` (must be sudo-capable) |
| SSH password for `moruogrsadmin` | `<SSH_PASSWORD>` (use ssh-copy-id ASAP — see §3) |
| SSH port on the host | `2222` (custom — not the default 22) |

DNS: ensure `moruo-gr-api.aart-enterprise.com` has an `A` record pointing at
`102.133.230.104` **before** running certbot — Let's Encrypt validates over HTTP-01 and
will fail otherwise.

---

## 3. Provision the Ubuntu host

Use **Ubuntu 24.04 LTS** (or 22.04 LTS). All commands assume you SSH in as
`moruogrsadmin` on the **non-standard SSH port 2222** — every `ssh` and `rsync` command
in this runbook passes the port explicitly with `-p 2222` / `-e "ssh -p 2222"`.

If you'd rather not type the port flag every time, drop a Host entry into your local
`~/.ssh/config` once and use the alias instead:

```sshconfig
Host aart-gr-prod
    HostName 102.133.230.104
    User     moruogrsadmin
    Port     2222
```

Then `ssh aart-gr-prod` and `rsync aart-gr-prod:...` Just-Work without flags. The runbook
keeps the explicit `-p 2222` form so the commands are copy-pasteable as-is.

**Switch to key-based auth before anything else.** You were given a password; replace it
with an SSH key on day one — passwords on a sudo-capable account are the worst-case auth
mechanism for this host.

```bash
# From your laptop, push your public key:
ssh-copy-id -p 2222 moruogrsadmin@102.133.230.104

# Then on the server, lock down sshd:
sudo sed -i \
  -e 's/^#\?PasswordAuthentication .*/PasswordAuthentication no/' \
  -e 's/^#\?PermitRootLogin .*/PermitRootLogin no/' \
  /etc/ssh/sshd_config
sudo systemctl restart ssh
# Verify from a NEW terminal (don't close the existing session yet) that key auth still works.
```

Now provision packages:

```bash
sudo apt update && sudo apt -y full-upgrade
sudo apt -y install \
    curl ca-certificates gnupg lsb-release ufw \
    nginx redis-server \
    mysql-client jq rsync logrotate \
    fail2ban
sudo systemctl enable --now nginx redis-server
```

The Go toolchain is **only** needed if you build on this host. If you build on a separate
build host (recommended for production), skip the Go install here.

```bash
# Optional — build on this host
GO_VERSION=1.23.4
curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
  | sudo tar -C /usr/local -xz
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/golang.sh
```

Set the hostname so logs are useful:

```bash
sudo hostnamectl set-hostname aart-gr-api-prod
```

---

## 4. Directory layout, ownership, firewall

The binary, config, and logs live under `/opt/aart-gr`, owned by `moruogrsadmin`. The
layout separates **release artifacts** (immutable, versioned) from **shared state**
(config, logs) so you can swap binaries without touching config.

> **Security trade-off worth knowing.** `moruogrsadmin` is a sudo-capable login user, so
> a compromise of the API process gives an attacker a path to the same shell and SSH keys
> the deploy operator uses. A dedicated `--system --shell /usr/sbin/nologin` user would
> be tighter. We compensate via the systemd hardening directives in §10
> (`NoNewPrivileges`, `ProtectSystem=strict`, `ProtectHome`, `PrivateTmp`, etc.) which
> sandbox the API process so it can't read `/home/moruogrsadmin`, write outside its
> allowed paths, or escalate privileges from inside the process. If you ever decide to
> separate the runtime user from the deploy user, only §4, §6, §7, and §10 change.

Create the directory tree, owned by `moruogrsadmin`:

```bash
sudo install -d -o moruogrsadmin -g moruogrsadmin -m 0750 /opt/aart-gr
sudo install -d -o moruogrsadmin -g moruogrsadmin -m 0750 /opt/aart-gr/releases
sudo install -d -o moruogrsadmin -g moruogrsadmin -m 0750 /opt/aart-gr/shared
sudo install -d -o moruogrsadmin -g moruogrsadmin -m 0750 /opt/aart-gr/shared/config
sudo install -d -o moruogrsadmin -g moruogrsadmin -m 0750 /opt/aart-gr/shared/logs
sudo install -d -o moruogrsadmin -g moruogrsadmin -m 0750 /opt/aart-gr/shared/uploads
```

Final layout, after a release is staged:

```
/opt/aart-gr/                 (owned by moruogrsadmin:moruogrsadmin)
├── current             -> releases/2026-05-03T1830Z   (symlink, atomic swap)
├── releases/
│   └── 2026-05-03T1830Z/
│       ├── aart_api    (the Go binary)
│       └── config.json -> /opt/aart-gr/shared/config/config.json   (symlink)
└── shared/
    ├── config/
    │   └── config.json
    ├── logs/           (ad-hoc; primary logging is journald)
    └── uploads/        (projection upload staging — see §6)
```

After this point, `moruogrsadmin` can write to `/opt/aart-gr/...` directly (no sudo
needed for release directories or staging the binary). Sudo is still needed for systemd,
Nginx, and anything under `/etc`.

Firewall — only 2222 (SSH), 80, and 443 reachable externally. The instance ports
9091/9092 are bound to all interfaces by the Go server, so we explicitly close them at
the firewall:

```bash
sudo ufw default deny incoming
sudo ufw default allow outgoing

# IMPORTANT: allow the custom SSH port BEFORE enabling UFW. The `OpenSSH` app profile
# only opens port 22 — using it here would lock you out, since sshd listens on 2222.
sudo ufw allow 2222/tcp comment 'SSH (custom port)'

sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw --force enable
sudo ufw status verbose
```

**Lockout safety check.** Before you run `ufw --force enable`, confirm in another
terminal that `ufw status` shows `2222/tcp ALLOW`. If it doesn't, fix that first — once
UFW is enabled with default-deny incoming and no rule for 2222, your active SSH session
will be the last one until someone with console access reverses it.

UFW's default `allow outgoing` covers MySQL outbound. No inbound rule for MySQL is needed
— the API connects out.

**Cloud network firewall — separate from UFW.** If the VM is on Azure, AWS, GCP, or any
other IaaS, there is a **second firewall** at the cloud network layer (Azure Network
Security Group, AWS Security Group, GCP VPC firewall rule). Packets are filtered there
**before** they ever reach the VM, so UFW's `allow 80/tcp` is necessary but not
sufficient — the cloud firewall must also allow 80 and 443 inbound, or Let's Encrypt's
HTTP-01 validator will time out at the TCP-connect stage with:

```
Detail: 102.133.230.104: Fetching http://moruo-gr-api.aart-enterprise.com/.well-known/acme-challenge/...:
Timeout during connect (likely firewall problem)
```

Open the right ports at the cloud layer before §13. For Azure (the most common case):

```bash
# Substitute your resource group and NSG name.
az network nsg rule create -g <RG> --nsg-name <NSG> \
  --name AllowHTTP  --priority 1010 --protocol Tcp \
  --access Allow --direction Inbound \
  --destination-port-ranges 80 --source-address-prefixes '*'

az network nsg rule create -g <RG> --nsg-name <NSG> \
  --name AllowHTTPS --priority 1020 --protocol Tcp \
  --access Allow --direction Inbound \
  --destination-port-ranges 443 --source-address-prefixes '*'
```

Equivalent in the Azure portal: VM → Networking → Network Security Group → Inbound
security rules → Add. The custom SSH port (2222) needs the same treatment if it isn't
already open — without it you wouldn't be SSH'd in, but if you're recreating the NSG
from scratch, remember to add `2222/tcp` first or you'll lock yourself out.

Verify from off-box (laptop, phone hotspot — anywhere not on the VM's network):

```bash
curl -v --max-time 10 http://moruo-gr-api.aart-enterprise.com/
# A timeout here means the cloud firewall is still blocking. Anything else (502, 200,
# refused) means packets are getting through and the issue is elsewhere.
```

---

## 5. Provision the database on the shared MySQL host

Run these against the **shared MySQL host** as an admin user.

The `aart_app` user already exists on this server, scoped to the aart-valuations API
host. MySQL keys users by `(user, host)`, so creating `aart_app` at a *new* host scope
(`102.133.230.104`) is a brand-new user row — it does not collide with the existing one.
You can use the same password as the existing `aart_app`@`<aart-valuations IP>` or a
different one; either is valid.

```sql
-- New schema for group risk.
CREATE DATABASE aart_grouprisk_db
    CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- New user row for the group-risk API host. (user, host) is the key — this does not
-- conflict with 'aart_app' scoped to the aart-valuations host.
CREATE USER 'aart_app'@'102.133.230.104' IDENTIFIED BY '<DB_PASSWORD>';
GRANT ALL PRIVILEGES ON aart_grouprisk_db.* TO 'aart_app'@'102.133.230.104';
FLUSH PRIVILEGES;

-- Sanity check (should list both host scopes for aart_app):
SELECT User, Host FROM mysql.user WHERE User = 'aart_app';
```

The migration runner needs DDL grants (`CREATE`, `ALTER`, `DROP`, `INDEX`, `REFERENCES`)
— `ALL PRIVILEGES ON aart_grouprisk_db.*` covers this and is the simplest grant. If your
security policy requires least-privilege, the minimum set is `SELECT, INSERT, UPDATE,
DELETE, CREATE, ALTER, DROP, INDEX, REFERENCES, CREATE ROUTINE, ALTER ROUTINE, EXECUTE,
CREATE VIEW`.

**TLS to the DB.** If your MySQL server enforces TLS (recommended for cross-host
connections), GORM's MySQL driver supports it via the DSN. The stock AART config object
doesn't currently expose a `tls` flag in the DSN, so for the first cut we recommend
either:

- Putting the API and DB on a private network so cleartext is acceptable, or
- Running an `stunnel` / WireGuard tunnel from the API host to the DB host and pointing
  `db_host` at `127.0.0.1:<tunnel_port>`.

A code change to thread `?tls=true` through the DSN is the cleaner long-term fix, but is
out of scope for this deployment.

Validate connectivity from the API host (logged in as `moruogrsadmin`) before going
further:

```bash
mysql -h <DB_HOST> -P 3306 -u aart_app -p aart_grouprisk_db -e 'SELECT VERSION();'
```

If that fails, fix it now — every later step assumes the DB is reachable.

---

## 6. Build the Go binary

Build on a separate build host (CI runner, your laptop, or a dedicated VM). The output is
a single static-ish binary you `rsync` to the server. Cross-compile from any OS:

```bash
cd api
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -trimpath -ldflags="-s -w" -o aart_api .
```

`CGO_ENABLED=0` produces a fully-static binary so you don't have to match libc versions
on the server. `-trimpath` and `-s -w` strip build paths and debug symbols — smaller
binary, no leaked filesystem layout.

Sanity-check the binary on the build host:

```bash
file aart_api
# expect: ELF 64-bit LSB executable, x86-64, ... statically linked
./aart_api --help 2>&1 | head -5  # just to confirm it runs
```

Stage it on the API server in a timestamped release directory. Because `moruogrsadmin`
owns `/opt/aart-gr`, no `sudo` is needed for the release directory itself — `rsync`
writes straight there.

> **`$RELEASE` lives in *this* (build-host) shell only.** It is expanded locally by
> bash before each `ssh` / `rsync` runs — the variable does *not* travel to the VM. If
> you `exit` this session before §11, the variable is gone. §11 handles this by
> rediscovering the release directory on the VM rather than relying on the variable
> surviving the SSH handoff.

```bash
RELEASE=$(date -u +%Y-%m-%dT%H%MZ)
ssh -p 2222 moruogrsadmin@102.133.230.104 "mkdir -p /opt/aart-gr/releases/$RELEASE"

# 1. The binary.
# Note: do NOT use rsync's --chown flag here. macOS ships rsync 2.6.9 which doesn't
# support it, and it isn't needed since we rsync as the owning user (moruogrsadmin).
rsync -av -e "ssh -p 2222" api/aart_api \
  moruogrsadmin@102.133.230.104:/opt/aart-gr/releases/$RELEASE/aart_api

# 2. The migrations tree. The runner reads from "migrations/" relative to the binary's
#    working directory; if absent it logs a Warn and silently applies nothing — schema
#    drift with no error. Always ship the directory with each release.
rsync -av --delete -e "ssh -p 2222" api/migrations/ \
  moruogrsadmin@102.133.230.104:/opt/aart-gr/releases/$RELEASE/migrations/

# 3. config + writable-state symlinks. The API uses relative paths for `logs/app.log`
#    (`api/log/logger.go:60`) and `tmp/uploads` (`api/controllers/projections.go:50`);
#    both would land inside the release dir, which is read-only under
#    ProtectSystem=strict. Symlinking them out into /opt/aart-gr/shared/ — which IS in
#    ReadWritePaths — is the fix that keeps the API code unchanged.
ssh -p 2222 moruogrsadmin@102.133.230.104 "
  chmod 0750 /opt/aart-gr/releases/$RELEASE/aart_api &&
  ln -sf  /opt/aart-gr/shared/config/config.json /opt/aart-gr/releases/$RELEASE/config.json &&
  ln -sfn /opt/aart-gr/shared/logs              /opt/aart-gr/releases/$RELEASE/logs &&
  mkdir -p /opt/aart-gr/releases/$RELEASE/tmp &&
  ln -sfn /opt/aart-gr/shared/uploads           /opt/aart-gr/releases/$RELEASE/tmp/uploads &&
  ls -la /opt/aart-gr/releases/$RELEASE/
"
```

We do **not** flip the `current` symlink yet — that happens in §11 once config and
systemd are in place.

---

## 7. Author `config.json`

The Go binary loads `config.json` from its working directory on Linux. Both instances
will share one file (with `app_port` ignored — overridden by env var per instance).

Create `/opt/aart-gr/shared/config/config.json` on the server. As `moruogrsadmin`:

```bash
tee /opt/aart-gr/shared/config/config.json >/dev/null <<'JSON'
{
  "db_type":        "mysql",
  "db_host":        "<DB_HOST>",
  "db_port":        "3306",
  "db_name":        "aart_grouprisk_db",
  "db_user":        "aart_app",
  "db_password":    "<DB_PASSWORD>",
  "app_host":       "moruo-gr-api.aart-enterprise.com",
  "app_port":       "9091",
  "redis_enabled":  true,
  "redis_host":     "127.0.0.1",
  "redis_port":     "6379",
  "redis_password": "<REDIS_PASSWORD>",
  "redis_db":       0
}
JSON
chmod 0640 /opt/aart-gr/shared/config/config.json
```

Notes:

- **JSON keys MUST be snake_case.** `models.AppConfig` (`api/models/setup.go:3-19`) uses
  explicit JSON tags like `json:"db_type"`, `json:"app_host"`. Go's JSON parser respects
  the tag and ignores the Go field name when one is set, so `"DbType"` in the file would
  parse silently into a zero value — every field empty, no error logged. If the API
  boots with `db_host=` and `redis_enabled=false` even though you set them, this is the
  cause.
- `app_host` is used to render Swagger URLs and is the hostname the app advertises. Set
  it to your public domain even though Nginx terminates TLS — downstream URL generation
  reads this value.
- `app_port` here is a fallback. Each systemd instance sets `APP_PORT` explicitly, which
  the binary uses preferentially (see `api/main.go:345-350`). The value in the file is
  only exercised if the env var is unset, which it never is in our setup.
- `redis_enabled: true` is critical for two-instance WebSocket fan-out.
- File mode `0640` keeps the DB password from being world-readable; group is
  `moruogrsadmin`, which is also the runtime user, so the API can read it.

---

## 8. Configure Redis (local, password-protected)

Redis was installed in §3 with default settings (binds to 127.0.0.1, no password). Add a
password since it's a shared component:

```bash
sudo sed -i 's/^# requirepass .*/requirepass <REDIS_PASSWORD>/' /etc/redis/redis.conf
sudo systemctl restart redis-server
redis-cli -a '<REDIS_PASSWORD>' ping   # expect: PONG
```

Confirm it's bound to localhost only:

```bash
ss -tlnp | grep 6379   # expect 127.0.0.1:6379, not 0.0.0.0:6379
```

If you ever need to move Redis off-box, remember to update `redis_host` in `config.json`
and re-open the Redis port over the private network only.

---

## 9. Run migrations once, before the first boot

The migration runner is **fatal on failure** and runs at every API start
(`services.RunMigrationsOnStartup`, called from `api/main.go:277`). With two instances
starting at once, both will race against the same `migrations` table. The runner does
tolerate idempotent conflicts, but the safer pattern is to apply migrations **once** from
a single process before bringing up the instances.

You have two ways to do this — pick one.

**Option A — first-instance-only, then second.** Start instance 1 alone, watch it apply
migrations cleanly, then start instance 2. This is what we do in §11 (instance 2 is
started with a 5-second `ExecStartPre=sleep`). The migration runner finishes well within
that window for incremental migrations. Good enough for normal deploys.

**Option B — pre-deploy migration job.** For large or risky migrations, run the binary
once with both instances *stopped*, then stop it after migrations finish:

```bash
# As moruogrsadmin on the API host:
cd /opt/aart-gr/releases/<RELEASE>
APP_PORT=9099 timeout 60 ./aart_api &
# watch the logs; ctrl-c after you see "Running database migrations" → "started server"
```

For day-one (this deploy) Option A is fine — there's no schema yet, so the fresh-install
branch in `services.SetupTables` will create everything and call
`MarkAllMigrationsAsApplied` to record every existing `.sql` file as already applied.
Refer to `api/MIGRATIONS.md` if the runner exits non-zero.

---

## 10. systemd template unit

A single template unit file lets us run as many instances as we like, parameterised by
port via the `%i` placeholder. Both instances run as `moruogrsadmin`.

```bash
sudo tee /etc/systemd/system/aart-gr-api@.service >/dev/null <<'UNIT'
[Unit]
Description=AART Group Risk API instance on port %i
After=network-online.target redis-server.service
Wants=network-online.target

[Service]
Type=simple
User=moruogrsadmin
Group=moruogrsadmin
WorkingDirectory=/opt/aart-gr/current
ExecStart=/opt/aart-gr/current/aart_api
Environment=APP_PORT=%i
# Do NOT set ENVIRONMENT=production — that flips the Go binary into AutoTLS mode
# and it will try to bind 80/443 itself (see api/main.go:367-374).

Restart=always
RestartSec=3
TimeoutStopSec=30

# stdout/stderr go to journald; query with: journalctl -u aart-gr-api@9091
StandardOutput=journal
StandardError=journal
SyslogIdentifier=aart-gr-api-%i

# Hardening — these matter more here because moruogrsadmin is a login/sudo user.
# They sandbox the API process so a compromise can't trivially read /home/moruogrsadmin,
# escalate privileges, or write outside its allowed paths.
NoNewPrivileges=yes
ProtectSystem=strict
ProtectHome=yes
PrivateTmp=yes
PrivateDevices=yes
ProtectKernelTunables=yes
ProtectKernelModules=yes
ProtectControlGroups=yes
RestrictSUIDSGID=yes
LockPersonality=yes
RestrictRealtime=yes
# Writable paths the API actually needs at runtime:
#   - shared/logs    : lumberjack writes logs/app.log here (via release-dir symlink)
#   - shared/uploads : projections.go writes tmp/uploads here (via release-dir symlink)
# Everything else is read-only thanks to ProtectSystem=strict above. If you see
# "read-only file system" errors in journald, it's almost always because something
# is trying to write to a path NOT listed here.
ReadWritePaths=/opt/aart-gr/shared/logs /opt/aart-gr/shared/uploads

# File descriptors — the WS hub plus DB connections plus client sockets adds up.
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
UNIT
```

Add a small drop-in for instance 2 only, to delay its start so instance 1 wins the
migration race in normal deploys (§9 Option A):

```bash
sudo install -d /etc/systemd/system/aart-gr-api@9092.service.d
sudo tee /etc/systemd/system/aart-gr-api@9092.service.d/delay.conf >/dev/null <<'EOF'
[Service]
ExecStartPre=/bin/sleep 5
EOF
```

Reload systemd:

```bash
sudo systemctl daemon-reload
```

`ProtectHome=yes` denies the API access to `/home`, including `/home/moruogrsadmin`. The
binary lives in `/opt/aart-gr/...`, so this is fine — but if you ever stash anything in
the home directory and expect the API to read it, you'll need to either move it to
`/opt/aart-gr/shared/...` or relax `ProtectHome`. The first option is correct.

---

## 11. First boot — flip the symlink, start instance 1, then 2

The `$RELEASE` variable from §6 lives in your **build-host** shell — it is not carried
into the VM when you SSH in. Rather than re-typing the timestamp (and risking a typo
that leaves `current` pointing at a dangling target), discover the latest release on
the VM itself and verify the binary is reachable before handing off to systemd.

Run this **on the VM** as `moruogrsadmin`:

```bash
# Pick the most recently created release directory.
RELEASE=$(ls -1t /opt/aart-gr/releases | head -1)

# Hard guard: refuse to flip the symlink unless the binary is actually there and
# executable. This is what prevents the "status=203/EXEC" restart loop — systemd
# fails 203/EXEC if WorkingDirectory or ExecStart can't be resolved or executed.
if [ -z "$RELEASE" ] || [ ! -x "/opt/aart-gr/releases/$RELEASE/aart_api" ]; then
    echo "ABORT: no usable release found. /opt/aart-gr/releases/$RELEASE/aart_api missing or not executable."
    ls -la /opt/aart-gr/releases/ 2>/dev/null
    return 1 2>/dev/null || exit 1
fi

echo "Linking /opt/aart-gr/current -> /opt/aart-gr/releases/$RELEASE"
ln -sfn /opt/aart-gr/releases/$RELEASE /opt/aart-gr/current

# Sanity check the symlink resolves before starting systemd:
ls -la /opt/aart-gr/current/aart_api
file   /opt/aart-gr/current/aart_api
# expect: ELF 64-bit LSB executable, x86-64, ... statically linked
```

Now start the instances:

```bash
# Start instance 1 — this is the one that runs migrations.
sudo systemctl enable --now aart-gr-api@9091
journalctl -u aart-gr-api@9091 -f --since '1 minute ago'
# Watch for: "Running database migrations" → "Application configuration loaded"
#            → "Starting development server"  (development is the LOG label; the binary
#            is in production mode because ENVIRONMENT is unset and we just don't use the
#            autotls branch)
# Ctrl-C out of the follow once you see it listening.

# Verify it answers locally. /health has no auth and no build prerequisites — the
# cheapest, most reliable live-check.
curl -sS -o /dev/null -w 'HTTP %{http_code}\n' http://127.0.0.1:9091/health
# expect: HTTP 200

# Start instance 2.
sudo systemctl enable --now aart-gr-api@9092
journalctl -u aart-gr-api@9092 --since '30 seconds ago'
curl -sS -o /dev/null -w 'HTTP %{http_code}\n' http://127.0.0.1:9092/health
```

If either instance fails to start, the most common culprits are: (a) `status=203/EXEC`
in journald with no application output — the `current` symlink is dangling or
`aart_api` is missing/not executable inside it; re-run the `RELEASE=$(ls -1t ...)`
block above and confirm `ls -la /opt/aart-gr/current/aart_api` shows a real file; (b)
the DB user can't connect from the API server's IP (verify the GRANT scope from §5 is
`'aart_app'@'102.133.230.104'`, not the aart-valuations IP); (c) `config.json` is
unreadable to `moruogrsadmin` (check `ls -l /opt/aart-gr/shared/config/config.json` —
owner should be `moruogrsadmin`); (d) Redis password mismatch.

---

## 12. Nginx reverse proxy — initial HTTP-only config

Nginx setup is two steps to avoid a chicken-and-egg with certbot: §12 lays down a port-80
config (so `nginx -t` passes and the ACME HTTP-01 challenge can succeed), §13 obtains the
cert, and §13 then replaces the config with the full HTTPS version.

Why two steps: certbot's `--nginx` plugin runs `nginx -t` before doing anything. If your
config references `/etc/letsencrypt/options-ssl-nginx.conf`, `ssl-dhparams.pem`, or the
cert files before they exist, `nginx -t` fails — and certbot is the thing that creates
those files. Trying to write the final config first leads to:

```
nginx: [emerg] open() "/etc/letsencrypt/options-ssl-nginx.conf" failed (2: No such file or directory)
```

So we start with an HTTP-only config:

```bash
sudo tee /etc/nginx/sites-available/aart-gr-api.conf >/dev/null <<'NGINX'
upstream aart_gr_backend {
    ip_hash;                       # WebSocket-friendly affinity
    server 127.0.0.1:9091 max_fails=3 fail_timeout=10s;
    server 127.0.0.1:9092 max_fails=3 fail_timeout=10s;
    keepalive 32;
}

# Required for the Upgrade header to pass through cleanly.
map $http_upgrade $connection_upgrade {
    default upgrade;
    ''      close;
}

server {
    listen 80;
    listen [::]:80;
    server_name moruo-gr-api.aart-enterprise.com;

    # ACME HTTP-01 challenge path (certbot writes here)
    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    # WebSocket endpoint (/ws). Mirrors the §13 location block so we can validate
    # WS routing before TLS is in place. Same long timeouts + buffering off; the
    # only cleartext exposure is for the brief window between §12 and §13.
    location /ws {
        proxy_pass http://aart_gr_backend;
        proxy_http_version 1.1;
        proxy_set_header Host              $host;
        proxy_set_header X-Real-IP         $remote_addr;
        proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Upgrade           $http_upgrade;
        proxy_set_header Connection        $connection_upgrade;

        proxy_buffering         off;
        proxy_connect_timeout   10s;
        proxy_send_timeout      3600s;
        proxy_read_timeout      3600s;
    }

    # Proxy on port 80 too, so we can smoke-test the upstream before TLS is up.
    location / {
        proxy_pass http://aart_gr_backend;
        proxy_http_version 1.1;
        proxy_set_header Host              $host;
        proxy_set_header X-Real-IP         $remote_addr;
        proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Upgrade           $http_upgrade;
        proxy_set_header Connection        $connection_upgrade;
    }
}
NGINX

sudo ln -sf /etc/nginx/sites-available/aart-gr-api.conf /etc/nginx/sites-enabled/aart-gr-api.conf
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t && sudo systemctl reload nginx
```

Smoke test on port 80 — proves the upstream works before adding TLS:

```bash
# 1. /health is the canonical canary — registered at routes.go:16, no auth, no build
#    prerequisites. Any 200 here proves nginx → upstream → app end-to-end.
curl -sSI http://moruo-gr-api.aart-enterprise.com/health | head -3
# expect: HTTP/1.1 200 OK

# 1b. (Optional) Swagger UI through location /. Registered at routes.go:632 as
#     /swagger/*any. A 404 on this path doesn't mean the proxy is broken — confirm
#     by hitting :9091/:9092 directly. If both upstreams 404, the build was likely
#     run without `swag init` first (see CLAUDE.md). Skip if you don't need the
#     Swagger UI in production.
curl -sSI http://moruo-gr-api.aart-enterprise.com/swagger/index.html | head -3

# 2. WebSocket upgrade through location /ws. Sends a synthetic WS handshake to
#    confirm nginx routes /ws to the backend with the Upgrade headers intact.
#    A 101 means the upgrade succeeded. A 400/426 from the Go side ("expected
#    websocket") still proves the proxy is doing its job — the request reached
#    the upstream with the right headers; only the app-level WS auth/handshake
#    failed (likely because we didn't supply the query-param auth token).
#    A 502/504/timeout means the location/upstream wiring is wrong — fix
#    before moving on.
curl -sS -i --http1.1 --max-time 5 \
  -H 'Connection: Upgrade' \
  -H 'Upgrade: websocket' \
  -H 'Sec-WebSocket-Version: 13' \
  -H "Sec-WebSocket-Key: $(openssl rand -base64 16)" \
  http://moruo-gr-api.aart-enterprise.com/ws | head -5
# expect: HTTP/1.1 101 Switching Protocols   (ideal — upgrade fully went through)
#     or: HTTP/1.1 400/426 ...               (acceptable — proxy OK, app rejected)
# NOT:    HTTP/1.1 502/504, or hang/timeout  (proxy not reaching backend)
```

---

## 13. Issue the Let's Encrypt certificate, then write the final HTTPS config

Use snap-installed certbot (the only currently-supported install path):

```bash
sudo snap install --classic certbot
sudo ln -sf /snap/bin/certbot /usr/bin/certbot

# Now that nginx -t passes, certbot --nginx can validate, obtain the cert, and create
# its helper files (options-ssl-nginx.conf, ssl-dhparams.pem) under /etc/letsencrypt/.
sudo certbot --nginx \
  -d moruo-gr-api.aart-enterprise.com \
  -m <OPS_EMAIL> \
  --agree-tos --no-eff-email --redirect

# Auto-renew is installed as a systemd timer.
sudo systemctl list-timers | grep -i certbot
sudo certbot renew --dry-run

# Confirm everything we'll reference in the final config exists.
sudo ls /etc/letsencrypt/live/moruo-gr-api.aart-enterprise.com/fullchain.pem
sudo ls /etc/letsencrypt/live/moruo-gr-api.aart-enterprise.com/privkey.pem
sudo ls /etc/letsencrypt/options-ssl-nginx.conf
sudo ls /etc/letsencrypt/ssl-dhparams.pem
```

Certbot will have edited `aart-gr-api.conf` to add a basic HTTPS block and an HTTP→HTTPS
redirect. We now overwrite that with the **full** AART Group Risk config — same
structure but with the AART-specific tuning (timeouts, body size, gzip) plus a
dedicated `location /ws` block for the WebSocket endpoint.

The `/ws` block is the proactive fix for an issue first observed on the **aart-valuations**
deployment, where clients silently dropped and reconnected on a regular cadence — the
classic symptom of nginx killing idle WebSocket connections when no traffic flowed for
longer than `proxy_read_timeout`. Splitting WS into its own location lets us give it a
much longer read/send timeout and turn off proxy buffering, while keeping the 600s
timeouts for normal HTTP requests (projections still need them).

```bash
sudo tee /etc/nginx/sites-available/aart-gr-api.conf >/dev/null <<'NGINX'
upstream aart_gr_backend {
    ip_hash;                       # WebSocket-friendly affinity
    server 127.0.0.1:9091 max_fails=3 fail_timeout=10s;
    server 127.0.0.1:9092 max_fails=3 fail_timeout=10s;
    keepalive 32;
}

# Required for the Upgrade header to pass through cleanly.
map $http_upgrade $connection_upgrade {
    default upgrade;
    ''      close;
}

# HTTP -> HTTPS redirect, with the ACME path left open for renewals.
server {
    listen 80;
    listen [::]:80;
    server_name moruo-gr-api.aart-enterprise.com;

    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    location / {
        return 301 https://$host$request_uri;
    }
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name moruo-gr-api.aart-enterprise.com;

    ssl_certificate     /etc/letsencrypt/live/moruo-gr-api.aart-enterprise.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/moruo-gr-api.aart-enterprise.com/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;

    # AART-specific: bordereaux / excel uploads.
    client_max_body_size 100m;

    # Projections can run long. These cover the longest single HTTP request the app makes.
    proxy_connect_timeout    30s;
    proxy_send_timeout       600s;
    proxy_read_timeout       600s;
    send_timeout             600s;

    # Standard proxy headers.
    proxy_http_version 1.1;
    proxy_set_header Host              $host;
    proxy_set_header X-Real-IP         $remote_addr;
    proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;

    # WebSocket upgrade.
    proxy_set_header Upgrade    $http_upgrade;
    proxy_set_header Connection $connection_upgrade;

    # gzip for JSON / Swagger assets.
    gzip on;
    gzip_proxied any;
    gzip_types application/json application/javascript text/css text/plain application/xml;
    gzip_min_length 1024;

    # Dedicated WebSocket endpoint (/ws — registered in api/routes/routes.go:23).
    # Declared before `location /` because nginx prefix-matches in declaration order
    # for non-regex locations. Auth is handled inside the controller via a
    # query-param token, so no special header rewrite is needed here.
    location /ws {
        proxy_pass http://aart_gr_backend;
        proxy_http_version 1.1;
        proxy_set_header Host              $host;
        proxy_set_header X-Real-IP         $remote_addr;
        proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Upgrade           $http_upgrade;
        proxy_set_header Connection        $connection_upgrade;

        # WebSocket-specific tuning. The server-block timeouts above are sized for
        # projection HTTP requests (600s); for an idle WS that's still too short —
        # any quiet client gets the connection axed by nginx and reconnects.
        # 1 hour matches what worked for aart-valuations after the same issue.
        proxy_buffering         off;
        proxy_connect_timeout   10s;
        proxy_send_timeout      3600s;   # 1h cap on idle client→server silence
        proxy_read_timeout      3600s;   # 1h cap on idle server→client silence;
                                         # must exceed the Go WS hub's ping cadence
                                         # so periodic pings keep the connection alive.
    }

    location / {
        proxy_pass http://aart_gr_backend;
    }

    access_log /var/log/nginx/aart-gr-api.access.log;
    error_log  /var/log/nginx/aart-gr-api.error.log warn;
}
NGINX

sudo nginx -t && sudo systemctl reload nginx
```

Smoke test over HTTPS:

```bash
# Canonical canary — /health is route-registered with no auth and no build deps,
# so a 200 here proves DNS → TLS → nginx → upstream → app end-to-end.
curl -sSI https://moruo-gr-api.aart-enterprise.com/health | head -3
# expect: HTTP/2 200, valid TLS, no warnings.
```

If you ever need to rebuild Nginx from scratch later, the workflow stays the same: §12
HTTP-only first, certbot, §13 full HTTPS. Once the cert and helper files exist, you can
of course skip §12 — the full config validates fine on its own from then on.

---

## 14. Smoke tests

Run all of these from off-box (your laptop) using the public domain — they exercise the
whole path: DNS → Nginx → upstream selection → API → DB.

```bash
# 1. TLS is valid, HTTP/2 negotiates, and the full path works end-to-end.
#    /health is the canonical canary — no auth, no build deps.
curl -sSI https://moruo-gr-api.aart-enterprise.com/health | head -3
# expect: HTTP/2 200, valid TLS

# 2. Same endpoint with body — controllers.CheckHealth returns a small JSON payload.
curl -sS https://moruo-gr-api.aart-enterprise.com/health
# expect: 200 OK with whatever payload controllers.CheckHealth returns.

# 3. Both backends are reachable. Repeat ~10 times — ip_hash will pin you to one,
#    so to confirm both are live, run from two different source IPs (laptop + phone hotspot)
#    or temporarily switch the upstream to round-robin and re-test.
for i in $(seq 1 5); do
  curl -sS https://moruo-gr-api.aart-enterprise.com/api/v1/...your-cheapest-endpoint... | head -1
done

# 4. Health of each backend, locally on the server.
ssh -p 2222 moruogrsadmin@102.133.230.104 'curl -s -o /dev/null -w "9091: %{http_code}\n" http://127.0.0.1:9091/health'
ssh -p 2222 moruogrsadmin@102.133.230.104 'curl -s -o /dev/null -w "9092: %{http_code}\n" http://127.0.0.1:9092/health'

# 5. Logs for the last few minutes — should be quiet, no panics.
ssh -p 2222 moruogrsadmin@102.133.230.104 \
  'journalctl -u aart-gr-api@9091 -u aart-gr-api@9092 --since "5 minutes ago" | tail -50'

# 6. From the Electron app, point its API base URL at
#    https://moruo-gr-api.aart-enterprise.com/api/v1/ and walk a real flow: log in, list
#    valuations, run a small projection. Watch logs while you do.
```

If a request times out at exactly 60s, you missed the `proxy_*_timeout` settings — fix
and reload Nginx. If WebSocket events stop arriving after a minute, same root cause for
the WS read timeout.

---

## 15. Logging & rotation

The API logs in **two places**, by design:

1. **journald** — every line the binary writes to stdout/stderr is captured by systemd
   (because of `StandardOutput=journal` / `StandardError=journal` in §10). This is the
   primary, always-on log stream. Query with `journalctl -u 'aart-gr-api@*'`.
2. **`logs/app.log` on disk** — `api/log/logger.go:60` configures `lumberjack` with
   `Filename: "logs/app.log"` and a `MultiWriter(file, stdout)`. The file path is
   **relative** to the working directory (`/opt/aart-gr/current`), so on disk it lands
   at `/opt/aart-gr/current/logs/app.log`. The §6 deploy step symlinks
   `/opt/aart-gr/current/logs → /opt/aart-gr/shared/logs`, and §10's `ReadWritePaths`
   includes that target — so writes succeed.

If the symlink isn't there (or the symlink target isn't in `ReadWritePaths`), lumberjack
fails with errors like:

```
Failed to write to log, can't rename log file: rename logs/app.log
  logs/app-2026-05-01T08-19-25.026.log: read-only file system
```

The API keeps running because the MultiWriter still has stdout, but the file half is
dead — **only journald has the logs from that point on.** This is the most common
log-capture surprise. Diagnose it by:

```bash
# Confirm the release dir's logs/ is a symlink to /opt/aart-gr/shared/logs:
ls -la /opt/aart-gr/current/logs
# expect: lrwxrwxrwx ... logs -> /opt/aart-gr/shared/logs

# Confirm the target exists and is owned by moruogrsadmin:
ls -la /opt/aart-gr/shared/logs/
# expect: drwxr-x--- moruogrsadmin moruogrsadmin ... .

# Confirm the systemd unit has it in ReadWritePaths:
systemctl show aart-gr-api@9091 -p ReadWritePaths
# expect: ReadWritePaths=/opt/aart-gr/shared/logs /opt/aart-gr/shared/uploads

# Look for the read-only error in recent journal:
sudo journalctl -u 'aart-gr-api@*' --since '10 minutes ago' --no-pager | grep -iE 'read-only|rename log'
```

If `app.log` exists but is empty/stale, it's because the API has been running but
writing to a nowhere-symlink. To repair on a live deploy without changing release dirs:

```bash
# As moruogrsadmin on the server:
ln -sfn /opt/aart-gr/shared/logs /opt/aart-gr/current/logs
mkdir -p /opt/aart-gr/shared/uploads
ln -sfn /opt/aart-gr/shared/uploads /opt/aart-gr/current/tmp/uploads 2>/dev/null || \
  { mkdir -p /opt/aart-gr/current/tmp 2>/dev/null; ln -sfn /opt/aart-gr/shared/uploads /opt/aart-gr/current/tmp/uploads; }

# Reload the systemd unit if you also changed ReadWritePaths:
sudo systemctl daemon-reload
sudo systemctl restart aart-gr-api@9091
sleep 5
ls -la /opt/aart-gr/shared/logs/
# expect: app.log appearing within ~5s of restart
```

**Nginx access/error logs** need rotation — Ubuntu's default `/etc/logrotate.d/nginx`
covers the standard log paths. Confirm:

```bash
sudo logrotate -d /etc/logrotate.d/nginx 2>&1 | head -30
```

**Application file logs** — once `app.log` is being written, add a logrotate stanza so
lumberjack's own rotation isn't your only line of defense:

```bash
sudo tee /etc/logrotate.d/aart-gr-api >/dev/null <<'EOF'
/opt/aart-gr/shared/logs/*.log {
    daily
    rotate 14
    compress
    missingok
    notifempty
    copytruncate
    su moruogrsadmin moruogrsadmin
}
EOF
```

`copytruncate` matters: lumberjack keeps the file descriptor open, so any rotation by
logrotate that renames-then-creates would leave the API writing to the renamed file.
`copytruncate` copies the contents and truncates in place, preserving the FD.

---

## 16. Deployment / update workflow (rolling restart)

This is the steady-state procedure for shipping a new build, with zero externally-visible
downtime as long as one instance stays up.

```bash
# On the build host — produce a fresh binary.
cd api
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o aart_api .

# Stage to a new release directory on the server (no sudo — moruogrsadmin owns /opt/aart-gr).
RELEASE=$(date -u +%Y-%m-%dT%H%MZ)
ssh -p 2222 moruogrsadmin@102.133.230.104 "mkdir -p /opt/aart-gr/releases/$RELEASE"

# 1. Binary.
rsync -av -e "ssh -p 2222" \
  api/aart_api \
  moruogrsadmin@102.133.230.104:/opt/aart-gr/releases/$RELEASE/aart_api

# 2. Migrations. Without this, the runner finds no migrations directory, logs a Warn,
#    and applies nothing — schema drift with no error. Always ship.
rsync -av --delete -e "ssh -p 2222" \
  api/migrations/ \
  moruogrsadmin@102.133.230.104:/opt/aart-gr/releases/$RELEASE/migrations/

# 3. Permissions, config symlink, writable-state symlinks (logs, tmp/uploads).
ssh -p 2222 moruogrsadmin@102.133.230.104 "
  chmod 0750 /opt/aart-gr/releases/$RELEASE/aart_api &&
  ln -sf  /opt/aart-gr/shared/config/config.json /opt/aart-gr/releases/$RELEASE/config.json &&
  ln -sfn /opt/aart-gr/shared/logs              /opt/aart-gr/releases/$RELEASE/logs &&
  mkdir -p /opt/aart-gr/releases/$RELEASE/tmp &&
  ln -sfn /opt/aart-gr/shared/uploads           /opt/aart-gr/releases/$RELEASE/tmp/uploads
"

# 4. Flip the symlink atomically. The new release is now live for the next process start,
#    but currently-running processes still hold the old binary in memory.
#    Guard: if $RELEASE is empty (lost variable, wrong shell), the link would resolve
#    to /opt/aart-gr/releases/ — a directory with no aart_api in it — and the next
#    service restart would fail status=203/EXEC. Abort before that can happen.
[ -n "$RELEASE" ] || { echo "ABORT: \$RELEASE not set in this shell; re-export it from §16 step 1"; exit 1; }
ssh -p 2222 moruogrsadmin@102.133.230.104 "
  [ -x /opt/aart-gr/releases/$RELEASE/aart_api ] || {
    echo 'ABORT: /opt/aart-gr/releases/$RELEASE/aart_api missing or not executable'; exit 1;
  } &&
  ln -sfn /opt/aart-gr/releases/$RELEASE /opt/aart-gr/current
"

# 5. Rolling restart. Instance 1 first (it applies any new migrations from step 2).
ssh -p 2222 moruogrsadmin@102.133.230.104 '
  sudo systemctl restart aart-gr-api@9091 &&
  for i in $(seq 1 30); do
    code=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:9091/health);
    [ "$code" = "200" ] && break;
    sleep 1;
  done &&
  echo "instance 9091: HTTP $code (after ${i}s)" &&
  [ "$code" = "200" ] || { echo "ABORT: 9091 not healthy, not restarting 9092"; exit 1; } &&
  sudo systemctl restart aart-gr-api@9092 &&
  for i in $(seq 1 30); do
    code=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:9092/health);
    [ "$code" = "200" ] && break;
    sleep 1;
  done &&
  echo "instance 9092: HTTP $code (after ${i}s)"
'

# 6. Verify migrations actually applied (only meaningful if your release contained any).
ssh -p 2222 moruogrsadmin@102.133.230.104 \
  "journalctl -u aart-gr-api@9091 --since '2 minutes ago' --no-pager | grep -iE 'migration|migrate' | tail -10"
# expect at least: "Running migrations on startup" and "Migrations completed successfully"
# expect for new ones: "Applied migration: <timestamp>_<name>.sql"

# Prune old releases (keep last 5). No sudo — moruogrsadmin owns the dirs.
ssh -p 2222 moruogrsadmin@102.133.230.104 '
  cd /opt/aart-gr/releases &&
  ls -1tr | head -n -5 | xargs --no-run-if-empty rm -rf
'
```

Migration risk during a rolling restart: if the new release contains a migration that's
incompatible with the **old** binary (e.g. drops a column the old binary still reads),
you'll have a window between instance 1 (new) and instance 2 (old) where instance 2
crashes. Mitigations: prefer additive-only migrations (the AART generator defaults to
this — see `api/MIGRATIONS.md` "additive-only by default"); or take a brief full outage
by stopping both instances, applying, then starting both.

---

## 17. Rollback

The release-symlink layout makes rollback a one-liner if the issue is in the binary
itself:

```bash
# List recent releases.
ssh -p 2222 moruogrsadmin@102.133.230.104 'ls -1t /opt/aart-gr/releases'

# Flip back to the previous good release and restart both.
PREV=<previous-release-dir>
ssh -p 2222 moruogrsadmin@102.133.230.104 "
  ln -sfn /opt/aart-gr/releases/$PREV /opt/aart-gr/current &&
  sudo systemctl restart aart-gr-api@9091 &&
  sudo systemctl restart aart-gr-api@9092
"
```

If the issue was a migration that ran successfully but is application-incompatible, the
binary rollback alone may not be enough — you may need a forward-fix migration. The
runner does not roll migrations back; that's a deliberate design choice (see
`api/MIGRATIONS.md` "Safety notes").

For the DB, take a logical backup before any deploy that includes new migrations:

```bash
ssh -p 2222 moruogrsadmin@102.133.230.104 '
  mysqldump -h <DB_HOST> -u aart_app -p<DB_PASSWORD> \
    --single-transaction --routines --triggers aart_grouprisk_db \
    | gzip > /opt/aart-gr/shared/logs/aart-gr-$(date +%Y%m%d-%H%M).sql.gz
'
```

(Better still: take the backup on the DB host, not via the API host — and now that the
DB host is shared with aart-valuations, coordinate with that team's backup schedule so
you're not double-running mysqldump on the same instance during peak hours.)

---

## 18. Recommended follow-ups (not blocking the first deploy)

These would meaningfully improve operability but aren't required for the deployment to
work. Rank them by your own risk tolerance.

1. **Wire Nginx upstream health checks to `/health`.** The API already exposes
   `GET /health` (`api/routes/routes.go:16`, `controllers.CheckHealth`) and is exempt
   from auth (the route is registered before the `apiv1 := router.Group("",
   GetActiveUser())` middleware applies). Nginx OSS doesn't ship active health checks,
   but you can use `proxy_next_upstream error timeout http_502 http_503;` plus
   `max_fails`/`fail_timeout` (already in §13's upstream block) to take an unhealthy
   backend out of rotation passively. For active health checks, `/health` is also what
   you point Azure load balancer probes / external monitors at.

2. **Plumb MySQL TLS into the GORM DSN.** Add a `DbTLS bool` (or `DbTLSMode string`)
   field to `models.AppConfig` (`api/models/setup.go`) and append the appropriate DSN
   params in `services.SetupTables`. Required if the DB ever moves outside a private
   network. Especially relevant here because the DB is shared with aart-valuations —
   a TLS rollout would affect both deployments.

3. **CORS lockdown.** `api/main.go:336` currently sets `AllowOrigins = []string{"*"}`.
   For a production deployment with a known frontend origin, restrict this to the
   Electron app's origin (or remove CORS entirely if the API isn't browser-fetched).

4. **Observability.** Forward journald to a log aggregator (Loki, Datadog, Cloudwatch);
   put a Prometheus exporter (or just `node_exporter`) on the host. Per-instance metrics
   are valuable for confirming load is actually balanced.

5. **fail2ban for Nginx.** The base `fail2ban` install in §3 covers SSH; add an Nginx
   filter if the API is publicly exposed and you see brute-force attempts on Swagger or
   auth endpoints.

6. **Backup automation for MySQL.** Daily `mysqldump` to S3 (or equivalent) with a
   retention policy. Run it from the DB host or a dedicated backup host, not the API
   host. Coordinate with the aart-valuations team on the shared host.

7. **Per-environment config files.** If you eventually want staging + prod on different
   hosts, name configs `config.<env>.json` and wire the systemd unit to pass
   `-config=` (the binary doesn't currently accept that flag — would need a small
   `main.go` change to honour `CONFIG_PATH` env or a CLI flag).

8. **Separate deploy user from runtime user.** When you have time, create a system user
   (`aart-gr` with `--shell /usr/sbin/nologin`) that owns `/opt/aart-gr` and runs the
   systemd units, and keep `moruogrsadmin` only as the SSH/deploy user (with
   `sudo chown` steps in the deploy flow). Reduces the blast radius of an API process
   compromise.

---

## 19. Quick reference

| Action | Command |
| --- | --- |
| Status of both instances | `systemctl status 'aart-gr-api@*'` |
| Live logs | `journalctl -u 'aart-gr-api@*' -f` |
| Logs since last hour | `journalctl -u 'aart-gr-api@*' --since '1 hour ago'` |
| Restart one instance | `sudo systemctl restart aart-gr-api@9091` |
| Reload Nginx | `sudo nginx -t && sudo systemctl reload nginx` |
| Test renewal | `sudo certbot renew --dry-run` |
| Current symlink | `ls -l /opt/aart-gr/current` |
| Release list | `ls -1t /opt/aart-gr/releases` |
| MySQL connectivity from API host | `mysql -h <DB_HOST> -u aart_app -p aart_grouprisk_db -e 'SELECT 1'` |
| Redis ping | `redis-cli -a '<REDIS_PASSWORD>' ping` |
| Health (local) | `curl -s http://127.0.0.1:9091/health; curl -s http://127.0.0.1:9092/health` |
| Health (public) | `curl -sI https://moruo-gr-api.aart-enterprise.com/health` |

---

## 20. Risks / things I want you to know before you run this

- **The migrations runner is fatal-on-failure.** A bad migration takes the API down at
  boot. Read `api/MIGRATIONS.md` end-to-end before your first risky migration.
- **`ENVIRONMENT=production` is a footgun.** It activates Go's built-in Let's Encrypt
  (`api/main.go:367-374`) and tries to bind 80/443 from the Go process. The systemd unit
  deliberately omits it.
- **`/health` is the API's own health endpoint** (`controllers.CheckHealth`, registered
  in `api/routes/routes.go:16`, exempt from `GetActiveUser()` auth). Use it for Nginx
  passive health, Azure LB probes, and external uptime monitors.
- **WebSockets need three things: `ip_hash`, Redis, and a long `/ws` read timeout.**
  `ip_hash` pins clients to a backend, Redis fans events across instances
  (`services.StartRedisWSSubscriber()` in `api/services/ws_hub.go:330-346`), and the
  dedicated `location /ws` block in §13 gives idle connections an hour before nginx
  considers them dead. Drop any one of the three and you'll see the
  silently-disconnecting-and-reconnecting pattern that bit aart-valuations.
- **MySQL grants are scoped to the API host's IP (`102.133.230.104`).** If you re-IP the
  API server, the grants will break and the API will fail to start. Update the grant or
  use a wildcard host in MySQL with caution.
- **DB host is shared with aart-valuations.** A misconfigured GRANT, a runaway
  migration, or a heavy mysqldump on this host affects *both* deployments. Coordinate.
- **CORS is wildcard `*`** (`api/main.go:336`). Restrict before the API is reachable
  from a browser other than the Electron app.
- **Runtime user is `moruogrsadmin`, a sudo-capable login user.** A process compromise
  has more reach than it would under a dedicated `--system --shell /usr/sbin/nologin`
  user. The systemd hardening in §10 reduces the blast radius substantially, but
  consider follow-up #8 above when the deploy stabilises.
- **Password-based SSH was the bootstrap auth.** §3 walks you through swapping to
  key-only auth. Don't skip it — a sudo-capable account with a password on a public
  port 2222 will be brute-forced eventually.
