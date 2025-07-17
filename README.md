# Distributed Rate Limiter

## Introduction

This project is a distributed rate limiter for coarse-grained flow control. It is based on local rate limiting and cannot achieve full precision.

---

## Principle

- Each instance maintains a local rate limiter (using `golang.org/x/time/rate`).
- The global QPS is divided among all instances according to their weights.
- Node info (ID, weight) is pushed to all instances via HTTP API.
- Each instance dynamically adjusts its local rate according to the latest node list and weights.

---

## Design

- **NodeManager**: Manages node list and weights.
- **Limiter**: Local rate limiter, supports dynamic rate adjustment.
- **API**: Exposes HTTP API for external node updates.
- **main.go**: Entry point, initializes and starts the service.

---

## Rate Limiting Strategy

- The global QPS is divided among all instances by weight:
  - Each instance's rate = QPS × (instance weight / total weight)
- As long as node info is synchronized in time, the effect is close to global limiting, but not fully precise.

---

## Suitable Scenarios

- Suitable for coarse-grained flow control, such as API gateways, service clusters, etc.
- Not suitable for scenarios requiring strong consistency or millisecond-level precision per request.

---

## Usage

### 1. Start Service

```sh
go run main.go --id=inst1 --qps=1000 --weight=1 --port=8080
```

- `--id`: Unique ID for this instance
- `--qps`: Global QPS limit
- `--weight`: Weight for this instance
- `--port`: Listen port

### 2. Update Node List

POST all node IDs and weights to `/update_nodes` via HTTP POST:

```json
[
  {"id": "inst1", "weight": 1},
  {"id": "inst2", "weight": 2}
]
```

---

## Note

- Suitable for coarse-grained control, not fully precise. 