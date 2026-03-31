# Deployment

## Docker (Recommended)

### Build the Image

```bash
docker build -t nats-dashboard .

# With version tag
docker build --build-arg VERSION=v1.0.0 -t nats-dashboard:v1.0.0 .
```

### Run

```bash
docker run -d \
  --name nats-dashboard \
  -p 8080:8080 \
  -v /path/to/config.yaml:/etc/nats-dashboard/config.yaml:ro \
  -v dashboard-data:/data \
  nats-dashboard
```

The container:
- Runs as non-root user `app` (uid 1000)
- Expects config at `/etc/nats-dashboard/config.yaml`
- Stores the SQLite database in `/data`
- Listens on the port defined in the config (`listen` field)

### Docker Compose

For a complete local development stack with a 3-node NATS cluster:

```bash
docker compose up -d
```

This starts:
- `nats-1`, `nats-2`, `nats-3` — Clustered NATS servers with JetStream
- `dashboard` — The dashboard, configured to poll all three servers

Ports:
- `8080` — Dashboard UI
- `4222-4224` — NATS client connections
- `8222-8224` — NATS monitoring endpoints

## Binary Deployment

### Build

```bash
make build
# Produces: bin/nats-dashboard
```

The binary is statically linked (`CGO_ENABLED=0`) and self-contained. Copy it and the config file to your server.

### Systemd Service

Create `/etc/systemd/system/nats-dashboard.service`:

```ini
[Unit]
Description=NATS Dashboard
After=network.target

[Service]
Type=simple
User=nats-dashboard
ExecStart=/usr/local/bin/nats-dashboard -config /etc/nats-dashboard/config.yaml
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
# Copy files
sudo cp bin/nats-dashboard /usr/local/bin/
sudo mkdir -p /etc/nats-dashboard /var/lib/nats-dashboard
sudo cp config.yaml /etc/nats-dashboard/

# Create user
sudo useradd -r -s /usr/sbin/nologin nats-dashboard
sudo chown -R nats-dashboard: /var/lib/nats-dashboard

# Set data_dir in config.yaml to /var/lib/nats-dashboard

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable nats-dashboard
sudo systemctl start nats-dashboard
```

## Kubernetes

Example deployment manifest:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nats-dashboard
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nats-dashboard
  template:
    metadata:
      labels:
        app: nats-dashboard
    spec:
      containers:
        - name: dashboard
          image: nats-dashboard:latest
          ports:
            - containerPort: 8080
          volumeMounts:
            - name: config
              mountPath: /etc/nats-dashboard
              readOnly: true
            - name: data
              mountPath: /data
      volumes:
        - name: config
          configMap:
            name: nats-dashboard-config
        - name: data
          persistentVolumeClaim:
            claimName: nats-dashboard-data
---
apiVersion: v1
kind: Service
metadata:
  name: nats-dashboard
spec:
  selector:
    app: nats-dashboard
  ports:
    - port: 8080
      targetPort: 8080
```

Create the ConfigMap from your config file:

```bash
kubectl create configmap nats-dashboard-config --from-file=config.yaml
```

## Reverse Proxy

### Nginx

```nginx
server {
    listen 443 ssl;
    server_name dashboard.example.com;

    ssl_certificate     /etc/ssl/certs/dashboard.pem;
    ssl_certificate_key /etc/ssl/private/dashboard.key;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /api/ws {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_read_timeout 86400;
    }
}
```

The WebSocket endpoint (`/api/ws`) requires the `Upgrade` and `Connection` headers to be forwarded.

## Network Requirements

The dashboard needs HTTP access to the NATS monitoring port (default 8222) on each configured server. Ensure firewall rules allow:

- Dashboard -> NATS servers: TCP port 8222 (or custom monitoring port)
- Browsers -> Dashboard: TCP port 8080 (or custom listen port)

The dashboard does **not** connect to the NATS client port (4222). It only uses the HTTP monitoring API.

## Security Considerations

- Change the `session_secret` from the default value
- Change the default admin password after first login
- Use TLS termination via a reverse proxy for production
- If NATS monitoring endpoints use HTTPS, configure the `tls` section in the environment config
- The SQLite database contains only bcrypt-hashed passwords
- Session cookies are `httpOnly` and `SameSite=Strict`
