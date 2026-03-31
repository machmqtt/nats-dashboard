# Configuration Reference

The dashboard is configured via a YAML file, specified with the `-config` flag (default: `config.yaml`).

## Full Example

```yaml
listen: ":8080"
poll_interval: 5s
session_secret: "change-me-to-a-random-string"
data_dir: "./data"

# A default admin user (admin/admin) is created automatically on first startup.

environments:
  - name: production
    servers:
      - url: "http://nats-1.prod:8222"
      - url: "http://nats-2.prod:8222"
      - url: "http://nats-3.prod:8222"
    tls:
      ca_file: "/etc/ssl/certs/nats-ca.pem"
      insecure: false

  - name: staging
    servers:
      - url: "http://nats-staging:8222"
```

## Default Admin User

On first startup, if the user database is empty, a default administrator account is created automatically:

- **Username:** `admin`
- **Password:** `admin`
- **Role:** `admin`

Change this password immediately after first login. The default admin is only created when no users exist in the database. Subsequent startups skip this step.

Admin users can create additional users via the User Management page in the UI or via the admin API endpoints.

## Fields

### Top-Level

| Field            | Type     | Default      | Required | Description |
|------------------|----------|--------------|----------|-------------|
| `listen`         | string   | `":8080"`    | No       | HTTP listen address (`host:port` or `:port`) |
| `poll_interval`  | duration | `5s`         | No       | How often to poll NATS monitoring endpoints |
| `session_secret` | string   | —            | **Yes**  | Secret key for signing JWT session tokens. Must be changed from the default. |
| `data_dir`       | string   | `"./data"`   | No       | Directory for the SQLite database file |
| `environments`   | list     | —            | **Yes**  | At least one environment must be defined |

### `environments[]`

| Field     | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `name`    | string | **Yes**  | Display name for the environment (must be unique) |
| `servers` | list   | **Yes**  | At least one NATS server monitoring URL |
| `tls`     | object | No       | TLS settings for connecting to NATS monitoring endpoints |

### `environments[].servers[]`

| Field | Type   | Required | Description |
|-------|--------|----------|-------------|
| `url` | string | **Yes**  | NATS server monitoring URL (e.g., `http://nats:8222`) |

The URL should point to the NATS HTTP monitoring port (default 8222), not the client port (4222).

### `environments[].tls`

| Field      | Type   | Default | Description |
|------------|--------|---------|-------------|
| `ca_file`  | string | —       | Path to CA certificate PEM file for verifying the NATS server's TLS certificate |
| `insecure` | bool   | `false` | Skip TLS certificate verification. Use only for testing. |

## Polling Behavior

The collector uses a two-tier polling strategy:

- **Fast poll** (every `poll_interval`): `/varz`, `/routez`, `/gatewayz`, `/leafz`, `/healthz` — lightweight endpoints needed for topology and overview
- **Slow poll** (every 3x `poll_interval`): `/connz`, `/subsz`, `/jsz`, `/accountz` — heavier endpoints that return more data

With the default 5s interval, fast data updates every 5 seconds and slow data every 15 seconds.

Each poll fetches all servers in an environment concurrently.

## Session Tokens

Sessions use HMAC-SHA256 signed JWTs stored in an `httpOnly` cookie with `SameSite=Strict`. Tokens expire after 24 hours.

The `session_secret` must be:
- At least a few characters long (a random 32+ character string is recommended)
- Kept secret and consistent across restarts (changing it invalidates all sessions)

Generate a secret:

```bash
openssl rand -base64 32
```

## Data Directory

The `data_dir` contains a single SQLite database file (`dashboard.db`) that stores user accounts. The database is created automatically on first run with WAL journaling mode enabled.

The directory is created if it doesn't exist.

## Environment Variables

The dashboard does not read environment variables directly. All configuration is file-based. To inject secrets in containerized environments, use volume mounts or templated config files.

## Validation

The config loader validates:
1. `session_secret` is not empty
2. At least one environment is defined
3. Each environment has a non-empty `name`
4. Each environment has at least one server
