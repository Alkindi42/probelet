# Probelet

[![CI](https://github.com/Alkindi42/probelet/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/Alkindi42/probelet/actions/workflows/ci.yaml)
[![Release](https://img.shields.io/github/v/release/Alkindi42/probelet?color=blue)](https://github.com/Alkindi42/probelet/releases)
[![Docker](https://img.shields.io/badge/docker-ghcr.io-blue?logo=docker&logoColor=white)](https://ghcr.io/alkindi42/probelet)
[![License](https://img.shields.io/github/license/Alkindi42/probelet?color=green)](LICENSE)

**Probelet** provides minimal failure primitives for platform testing.

It can be used either as an HTTP service or as a local CLI for controlled stress workloads.

It combines the simplicity of tools like **httpbin** with **controlled system stress capabilities**, making it ideal for:

* testing Kubernetes liveness and readiness probes
* validating autoscaling (HPA)
* simulating slow or failing upstream services
* debugging timeouts and retry logic
* inspecting incoming HTTP requests
* exercising observability and monitoring stacks
* running controlled CPU, memory, and disk stress locally or over HTTP

## Why Probelet exists

Modern platforms rely on mechanisms such as retries, autoscaling, circuit breakers, and readiness probes to **react to failure**.

However, the conditions that trigger these mechanisms are rarely exercised deliberately. Most systems are validated on the **happy path**, while their behavior under degradation is often discovered only during incidents.

Probelet was created to make these failure conditions easier to reproduce.

It exposes **minimal primitives** that allow engineers to simulate response delays, transient errors, resource pressure, and readiness changes in a controlled and observable way.

This makes it possible to exercise platform reliability mechanisms — retries, autoscaling, health checks, and monitoring pipelines — without introducing complex test environments.

If you're interested in the ideas behind this tool, you can read the article:

👉 [You don’t test platform reliability until you test failure semantics](https://brainstacks.dev/posts/reliability-and-failure-semantics)

---

## 🚀 Installation

### Docker

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

### Binary

Download a release from: [https://github.com/Alkindi42/probelet/releases](https://github.com/Alkindi42/probelet/releases).

Example:

```bash
tar -xzf probelet-<version>-linux-amd64.tar.gz
./probelet serve
```

#### macOS

When downloading the binary directly, macOS may block execution due to
Gatekeeper quarantine.

You can remove the quarantine attribute with:

```bash
xattr -d com.apple.quarantine ./probelet
```

---

## 📚 API documentation

Probelet provides an interactive API reference powered by OpenAPI and Scalar.

The documentation describes all available endpoints used to simulate
application behavior, failure scenarios, and system stress workloads.

**Public documentation**

[https://alkindi42.github.io/probelet/](https://alkindi42.github.io/probelet/)

**Local documentation**

When running Probelet locally: [http://localhost:8000/docs](http://localhost:8000/docs)

OpenAPI specification: [http://localhost:8000/openapi.yaml](http://localhost:8000/openapi.yaml)

### 🔐 Authentication

Probelet can optionally protect selected endpoints with a static token.

Set the `PROBELET_TOKEN` environment variable when starting the service:

```bash
PROBELET_TOKEN=my-secret-token ./probelet serve
```

Then send the token in the `X-Probelet-Token` header when calling a protected endpoint:

```bash
curl -H "X-Probelet-Token: my-secret-token" "http://localhost:8000/delay?duration=5s"
```

Public endpoints remain accessible without authentication:

* `/`
* `/healthz`
* `/readyz`
* `/docs`
* `/openapi.yaml`

## 🧪 Examples

### API

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

Trigger disk pressure:

```bash
curl "http://localhost:8000/stress/disk?duration=10s&size=256Mi" 
```

Simulate an intermittently degraded upstream:

```bash
probelet stress memory --duration 30s --size 256Mi
```

### CLI

Run CPU stress:

```bash
probelet stress cpu --duration 10s --cores max
```

Run memory stress:

```bash
probelet stress memory --duration 30s --size 256Mi
```

Run disk stress:

```bash
probelet stress disk --duration 30s --size 256Mi
```

Print the version:

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
