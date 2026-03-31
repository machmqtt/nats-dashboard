# API Reference

All endpoints are served from the dashboard's HTTP server (default `:8080`).

## Authentication & Authorization

Authentication uses JWT tokens stored in an `httpOnly` cookie named `session`. The cookie is set on successful login and cleared on logout.

All `/api/*` endpoints except `POST /api/login` require authentication. Unauthenticated requests receive a `401 Unauthorized` response.

### Roles

Users have one of two roles:

- **admin** — Full access. Can create and delete users via `/api/admin/*` endpoints.
- **viewer** — Read-only access to all monitoring data. Cannot manage users.

Admin endpoints (`/api/admin/*`) return `403 Forbidden` for non-admin users.

---

## Auth Endpoints

### POST /api/login

Authenticate and receive a session cookie.

**Request body:**
```json
{
  "username": "admin",
  "password": "admin"
}
```

**Response (200):**
```json
{
  "id": 1,
  "username": "admin",
  "role": "admin",
  "created_at": "2026-03-20T10:00:00Z"
}
```

Sets a `session` httpOnly cookie.

**Response (401):** Invalid credentials.

### POST /api/logout

Clear the session cookie.

**Response (200):**
```json
{ "ok": true }
```

### GET /api/me

Get the current authenticated user.

**Response (200):**
```json
{
  "id": 1,
  "username": "admin",
  "role": "admin",
  "created_at": "2026-03-20T10:00:00Z"
}
```

### PUT /api/users/{id}/password

Change a user's password. Users can only change their own password.

**Request body:**
```json
{
  "old_password": "current",
  "new_password": "new-password"
}
```

**Response (200):**
```json
{ "ok": true }
```

**Response (403):** Trying to change another user's password.

---

## Admin Endpoints

These endpoints require the `admin` role. Non-admin users receive `403 Forbidden`.

### GET /api/admin/users

List all users.

**Response (200):**
```json
{
  "users": [
    {
      "id": 1,
      "username": "admin",
      "role": "admin",
      "created_at": "2026-03-20T10:00:00Z"
    },
    {
      "id": 2,
      "username": "viewer1",
      "role": "viewer",
      "created_at": "2026-03-20T12:00:00Z"
    }
  ]
}
```

### POST /api/admin/users

Create a new user. Admin role required.

**Request body:**
```json
{
  "username": "newuser",
  "password": "secure-password",
  "role": "viewer"
}
```

`role` is optional, defaults to `viewer`. Valid values: `admin`, `viewer`.

**Response (201):**
```json
{
  "id": 2,
  "username": "newuser",
  "role": "viewer",
  "created_at": "2026-03-20T12:00:00Z"
}
```

### DELETE /api/admin/users/{id}

Delete a user. Cannot delete your own account.

**Response (200):**
```json
{ "ok": true }
```

**Response (400):** Attempting to delete your own account.

**Response (404):** User not found.

---

## Environment Endpoints

All environment data endpoints are under `/api/environments/{env}/` where `{env}` is the environment name from the config file.

### GET /api/environments

List all configured environments.

**Response:**
```json
{
  "environments": ["production", "staging"]
}
```

### GET /api/environments/{env}/overview

Aggregated overview with server summaries.

**Response:**
```json
{
  "server_count": 3,
  "healthy_count": 3,
  "connection_count": 150,
  "in_msgs_rate": 1500.5,
  "out_msgs_rate": 1200.3,
  "in_bytes_rate": 50000,
  "out_bytes_rate": 40000,
  "subscriptions": 800,
  "js_streams": 5,
  "js_consumers": 12,
  "js_messages": 100000,
  "js_bytes": 5000000,
  "servers": [
    {
      "id": "NABC123",
      "name": "nats-1",
      "version": "2.11.0",
      "connections": 50,
      "cpu": 12.5,
      "mem": 1048576,
      "in_msgs_rate": 500.0,
      "out_msgs_rate": 400.0,
      "healthy": true,
      "uptime": "24h5m"
    }
  ]
}
```

### GET /api/environments/{env}/topology

Force-graph data for the cluster topology visualization.

**Response:**
```json
{
  "nodes": [
    {
      "id": "NABC123",
      "name": "nats-1",
      "type": "server",
      "connections": 50,
      "healthy": true,
      "in_msgs_rate": 500.0,
      "out_msgs_rate": 400.0,
      "cluster": "dc1"
    }
  ],
  "links": [
    {
      "source": "NABC123",
      "target": "NDEF456",
      "type": "route",
      "in_msgs_rate": 100.0,
      "out_msgs_rate": 80.0
    }
  ]
}
```

Node types: `server`, `gateway`, `leaf`.
Link types: `route`, `gateway`, `leaf`.

### GET /api/environments/{env}/varz

Per-server variable data, keyed by server ID.

**Response:**
```json
{
  "NABC123": {
    "server_id": "NABC123",
    "server_name": "nats-1",
    "version": "2.11.0",
    "host": "0.0.0.0",
    "port": 4222,
    "connections": 50,
    "in_msgs": 100000,
    "out_msgs": 80000,
    "in_bytes": 5000000,
    "out_bytes": 4000000,
    "mem": 1048576,
    "cpu": 12.5,
    "cores": 4,
    "subscriptions": 300,
    "uptime": "24h5m"
  }
}
```

### GET /api/environments/{env}/connz

Paginated connections list, fetched live from all servers.

**Query parameters:**

| Parameter        | Type   | Default | Description |
|------------------|--------|---------|-------------|
| `limit`          | int    | 256     | Max connections to return per server |
| `offset`         | int    | 0       | Pagination offset |
| `sort`           | string | —       | Sort field: `cid`, `start`, `subs`, `pending`, `msgs_to`, `msgs_from`, `bytes_to`, `bytes_from`, `idle`, `last` |
| `acc`            | string | —       | Filter by account name |
| `state`          | string | —       | Filter by state: `open`, `closed` |
| `filter_subject` | string | —       | Filter by subscription subject |

**Response:**
```json
{
  "connections": [
    {
      "cid": 5,
      "ip": "192.168.1.10",
      "port": 54321,
      "name": "my-service",
      "account": "$G",
      "authorized_user": "admin",
      "rtt": "1.5ms",
      "in_msgs": 1000,
      "out_msgs": 800,
      "in_bytes": 50000,
      "out_bytes": 40000,
      "subscriptions": 5,
      "uptime": "1h30m",
      "lang": "go",
      "version": "1.36.0"
    }
  ],
  "total": 150,
  "limit": 256,
  "offset": 0
}
```

### GET /api/environments/{env}/connz/{cid}

Single connection detail by CID (from cached snapshot).

**Response:** A single connection object (same fields as above).

### GET /api/environments/{env}/routez

Cluster route information, keyed by server ID.

### GET /api/environments/{env}/gatewayz

Gateway connections, keyed by server ID.

### GET /api/environments/{env}/leafz

Leaf node connections, keyed by server ID.

### GET /api/environments/{env}/subsz

Subscription statistics per server.

**Response:**
```json
{
  "NABC123": {
    "server_id": "NABC123",
    "num_subscriptions": 300,
    "num_cache": 100,
    "num_inserts": 5000,
    "num_removes": 4700,
    "num_matching": 10000,
    "cache_hit_rate": 85,
    "max_fanout": 10,
    "avg_fanout": 2.5
  }
}
```

### GET /api/environments/{env}/jsz

JetStream information with stream and consumer details, keyed by server ID.

### GET /api/environments/{env}/accountz

Account list, keyed by server ID.

### GET /api/environments/{env}/accountz/{acc}

Detailed information for a single account, fetched live.

**Response:**
```json
{
  "account_name": "$G",
  "is_system": false,
  "expired": false,
  "jetstream_enabled": true,
  "leafnode_connections": 0,
  "client_connections": 50,
  "subscriptions": 200
}
```

---

## WebSocket

### GET /api/ws

Upgrade to WebSocket for real-time data updates.

**Client -> Server:** Subscribe to an environment:
```json
{ "subscribe": "production" }
```

**Server -> Client:** The server pushes messages on each poll cycle:

```json
{ "type": "overview", "env": "production", "data": { ... } }
{ "type": "topology", "env": "production", "data": { ... } }
{ "type": "health",   "env": "production", "data": { ... } }
```

Message types:
- `overview` — Summary numbers (same shape as `/api/environments/{env}/overview`)
- `topology` — Full graph (same shape as `/api/environments/{env}/topology`)
- `health` — Per-server health status map

The WebSocket connection supports ping/pong keepalive (54s interval, 60s timeout).

To switch environments, send a new subscribe message. Only one environment subscription is active per connection.

---

## Static Assets

All non-`/api/` paths serve the embedded React SPA. Unknown paths fall back to `index.html` for client-side routing.
