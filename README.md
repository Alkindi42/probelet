# Probelet

**Probelet** is a lightweight HTTP service designed to **simulate application and system behaviors** for testing platforms and orchestrators.

It combines the simplicity of tools like `httpbin` with **controlled system stress capabilities** (CPU today, more to come), making it ideal for:

* testing Kubernetes liveness/readiness probes
* validating autoscaling (HPA)
* debugging timeouts and failure scenarios
* exercising observability and monitoring stacks

Probelet runs as a **single, stateless Go binary**, and exposes everything via a simple HTTP API.

---

## 🚀 Quick start (Docker)

```bash
docker run --rm -p 8000:8000 ghcr.io/alkindi42/probelet:latest
```

Verify it’s running:

```bash
curl http://localhost:8000/healthz
# {"ok":true,"message":"healthy"}
```

---

## 📡 HTTP API

Probelet listens on port **8000** by default.

### Health & readiness

| Endpoint   | Method | Description                                         |
| ---------- | ------ | --------------------------------------------------- |
| `/healthz` | GET    | Liveness probe (always healthy if process is alive) |
| `/readyz`  | GET    | Readiness probe                                     |
| `/readyz`  | POST   | Toggle readiness state                              |

Example:

```bash
curl http://localhost:8000/readyz
```

Mark the service as **not ready**:

```bash
curl -X POST http://localhost:8000/readyz \
  -H "Content-Type: application/json" \
  -d '{"ready":false,"reason":"dependency unavailable"}'
```

---

### HTTP status simulator

Return any HTTP status code on demand.

```bash
curl http://localhost:8000/status/503
```

Response:

```json
{"ok":true,"message":"service unavailable","data":{"code":503}}
```

Useful to test:

* retries
* error handling
* load balancers / gateways

---

### CPU stress testing

Trigger a **controlled CPU load** for a given duration.

```bash
curl "http://localhost:8000/stress/cpu?duration=10s"
```

Parameters:

| Query      | Required | Description                                 |
| ---------- | -------- | ------------------------------------------- |
| `duration` | yes      | Stress duration (`100ms`, `5s`, `2m`)       |
| `cores`    | no       | Number of CPU cores (default `1`, or `max`) |

Examples:

```bash
# Stress 1 core for 500ms
curl "http://localhost:8000/stress/cpu?duration=500ms"

# Stress all available cores for 10s
curl "http://localhost:8000/stress/cpu?duration=10s&cores=max"
```

**Safety defaults**

* maximum duration: **2 minutes**
* runs synchronously
* stops immediately if the client disconnects
* non-root container image

---

## ☸️ Kubernetes

A minimal Kubernetes example is available in `examples/kubernetes/`.

```bash
kubectl apply -f examples/kubernetes/
kubectl port-forward svc/probelet 8000:8000
```

Then try:

```bash
curl http://localhost:8000/healthz
curl http://localhost:8000/status/503
curl "http://localhost:8000/stress/cpu?duration=5s"
```

---

## 🧠 Design principles

* **Safe by default**: bounded stress, no background jobs
* **Stateless**: no persistence, easy to reset
* **Orchestrator-agnostic**: Kubernetes, Nomad, Docker, bare metal
* **Observable**: deterministic behaviors, simple APIs

---

## 🤝 Contributing

Contributions and ideas are welcome.
Feel free to open an issue or a pull request.
