# AART API — Installation & Setup

## Prerequisites

- A running database server. Supported backends: **MySQL**, **PostgreSQL**, **SQL Server**.
- A database user with read/write (and DDL) permissions on the target database. The database itself should exist — AART will create tables but does not create the database.
- (Optional) A Redis server, if you plan to enable Redis-backed caching and pub/sub.
- Network connectivity from the machine running the API to the database (and Redis, if enabled).

## 1. Obtain the binary

Download the release for your platform from
[the releases page](https://github.com/koreset/aart-api/releases) and place the
binary somewhere writable — e.g. `/opt/aart/` on Linux/macOS or `C:\aart\` on
Windows. The binary has no GUI; all interaction is via a terminal.

Below, we use `aart_api` as the binary name; substitute whatever your download
is called (e.g. `api-1.2.0-linux`).

## 2. Run the setup wizard

From the directory containing the binary, run:

```shell
./aart_api -service setup
```

The wizard walks through the configuration interactively and writes the result
to `config.json` in the current directory (mode `0600` — readable only by the
owner, since it contains the database password).

If `config.json` already exists, setup refuses to overwrite it. Use
[`-service reconfigure`](#reconfiguring-an-existing-install) instead.

### What the wizard asks

1. **Confirmation** — "Installing a new instance of AART. Have your database
   credentials ready. Continue?" Answer `Yes` to proceed.

2. **Database type** — choose `MySQL`, `PostgreSQL`, or `SQL Server`.

3. **Database connection details**:
   - `Database name` — the existing database AART should use.
   - `Database user`
   - `Database password` (masked)
   - `Database host (IP or hostname)` — whitespace and URL schemes are
     rejected; use a bare host like `db.internal` or `10.0.0.5`.
   - `Database port` — must be an integer between 1 and 65535. Common
     defaults: MySQL `3306`, PostgreSQL `5432`, SQL Server `1433`.

   After collecting these, the wizard **opens a test connection and pings
   the database**. If it fails, the error is printed and you're asked whether
   you want to re-enter credentials. Nothing is written to disk until the
   ping succeeds.

4. **Application host / port** — the hostname and port the API server
   should listen on. `9090` is the conventional default.

5. **Redis (optional)** — "Enable Redis?". Answer `No` to skip Redis
   entirely; the application will run fine without it (caching falls back to
   in-memory and pub/sub features are disabled). If you answer `Yes`, you
   will be asked for:
   - Redis host
   - Redis port
   - Redis password (may be left blank)
   - Redis DB number (0–15; defaults to 0)

6. **Write & run** — the wizard writes `config.json`, runs the initial
   schema setup against the database (creates tables on first boot), and
   then starts the API server.

## 3. Subsequent starts

Once `config.json` exists, you no longer need any flags. Just run:

```shell
./aart_api
```

The binary reads `config.json`, runs any pending migrations, starts the
background workers, and begins listening on the configured host/port.

## Reconfiguring an existing install

To change any setting (rotate a DB password, point at a different host,
toggle Redis, etc.) without hand-editing `config.json`:

```shell
./aart_api -service reconfigure
```

This loads the current `config.json`, runs the same wizard with **every
value pre-filled as the default** (press Enter to keep), and rewrites the
file. The DB password prompt accepts a blank value to keep the existing
password. The connection is still tested before anything is written.

Reconfigure does **not** restart the running service. After it completes,
restart the API process manually for the new settings to take effect.

## Running as a system service (Windows)

On Windows you can register AART as a Windows Service:

```shell
aart_api.exe -service install
```

You'll be prompted for a service name and a display name. The service is
installed and started automatically. To remove it:

```shell
aart_api.exe -service uninstall
```

`-service install` is **Windows-only**. On macOS/Linux, run the binary
directly or wrap it in a `launchd` / `systemd` unit of your own.

## The config file

`config.json` is a plain JSON document. A typical file looks like:

```json
{
  "db_type": "postgresql",
  "db_name": "aart",
  "db_host": "db.internal",
  "db_user": "aart",
  "db_password": "********",
  "db_port": "5432",
  "app_port": "9090",
  "app_host": "0.0.0.0",
  "redis_enabled": false,
  "redis_host": "",
  "redis_port": "",
  "redis_password": "",
  "redis_db": 0
}
```

Field reference:

| Field            | Description                                                   |
|------------------|---------------------------------------------------------------|
| `db_type`        | `mysql`, `postgresql`, or `mssql`                             |
| `db_name`        | Target database (must already exist)                          |
| `db_host`        | DB hostname or IP                                             |
| `db_user`        | DB username                                                   |
| `db_password`    | DB password (stored in plaintext — protect this file)         |
| `db_port`        | DB port                                                       |
| `app_host`       | Host the API binds to                                         |
| `app_port`       | Port the API listens on                                       |
| `redis_enabled`  | `true` to use Redis; `false` to disable (default)             |
| `redis_host`     | Redis hostname (only used when `redis_enabled`)               |
| `redis_port`     | Redis port                                                    |
| `redis_password` | Redis password (empty for no auth)                            |
| `redis_db`       | Redis logical DB index, 0–15                                  |

You can edit `config.json` by hand in an emergency, but
`-service reconfigure` is preferred because it also verifies the database
connection before writing.

## Environment overrides

A few fields can be overridden at runtime without editing the config:

- `APP_PORT` — overrides `app_port`.
- `ENVIRONMENT=production` — switches the HTTP listener to AutoTLS (Let's
  Encrypt) on `app_host`. Use this only on a publicly reachable host with a
  valid DNS name and ports 80/443 open.

## Troubleshooting

- **"config.json already exists"** when running `-service setup`: delete
  the file if you really want a clean setup, or use `-service reconfigure`.
- **"✗ Database connection failed"** during setup: the wizard prints the
  underlying driver error. Check host/port reachability, credentials, and
  that the database named in the prompt actually exists on the server.
- **Redis connection errors at boot** (when `redis_enabled: true`): the
  app logs a warning and continues without Redis — it does not crash. If
  you want to disable Redis entirely, set `redis_enabled: false` or run
  reconfigure and answer `No` to the Redis prompt.
- **Permission denied reading `config.json`**: the file is written as
  `0600`. Run the binary as the same user that ran `-service setup`, or
  `chown` the file to the runtime user.
