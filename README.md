# Probelet

[![CI](https://github.com/Alkindi42/probelet/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/Alkindi42/probelet/actions/workflows/ci.yaml)

**Probelet** is a lightweight tool for simulating application and system behaviors in test, staging, and platform environments.

It can be used either as an HTTP service or as a local CLI for controlled stress workloads.

It combines the simplicity of tools like `httpbin` with **controlled system stress capabilities**, making it ideal for:

* testing Kubernetes liveness and readiness probes
* validating autoscaling (HPA)
* simulating slow or failing upstream services
* debugging timeouts and retry logic
* inspecting incoming HTTP requests
* exercising observability and monitoring stacks
* running controlled CPU and memory stress locally or over HTTP

---

## 🚀 Quick start (Docker)

```bash
docker run --rm -p 8000:8000 ghcr.io/alkindi42/probelet:latest
```

Verify it’s running:

```bash
curl http://localhost:8000/healthz
```

Expected response:

```json
{"ok":true,"message":"healthy"}
```

---

## 📚 API documentation

Probelet exposes an interactive API reference powered by OpenAPI.

Once the service is running, open: [http://localhost:8000/docs](http://localhost:8000/docs)

The OpenAPI specification is also available at: [http://localhost:8000/openapi.yaml](http://localhost:8000/openapi.yaml)

## 🧪 Examples

Simulate a slow service:

```bash
curl "http://localhost:8000/delay?duration=3s"
```

Trigger CPU pressure:

```bash
curl "http://localhost:8000/stress/cpu?duration=10s&cores=max"
```

Return a specific HTTP status:

```bash
curl http://localhost:8000/status/503
```

Inspect an incoming request:

```bash
curl http://localhost:8000/echo
```

Run CPU stress locally from the CLI:

```bash
probelet stress cpu --duration 10s --cores max
```

Run memory stress locally from the CLI:

```bash
probelet stress memory --duration 30s --size 256Mi
```

Print the running build version:

```bash
probelet version
```

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

## 📄 License

This project is licensed under the [MIT License](LICENSE).
