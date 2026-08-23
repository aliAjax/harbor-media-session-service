# harbor-sfu

`harbor-sfu` is a pure-Go multi-party real-time media SFU control/data plane. It exposes a REST control API, WebSocket signaling, a UDP RTP ingress, bounded fan-out forwarding, RTCP quality feedback, simulcast selection, recording jobs, room leases with fencing tokens, tenant quotas, and Prometheus metrics.

## Quick start

```sh
go run ./cmd/harbor-sfu
curl -s http://127.0.0.1:8093/healthz
curl -s http://127.0.0.1:8093/readyz
curl -s -X POST http://127.0.0.1:8093/v1/rooms -H 'content-type: application/json' -d '{"id":"demo","tenant_id":"acme","name":"Demo"}'
curl -s -X POST http://127.0.0.1:8093/v1/rooms/demo/join -H 'content-type: application/json' -d '{"id":"alice","identity":"Alice","role":"publisher"}'
curl -s http://127.0.0.1:8093/metrics
```

The signaling endpoint is `ws://127.0.0.1:8093/v1/rooms/demo/signal`. Frames are JSON and use monotonic `seq` values. `POST /publish` registers an RTP track; UDP packets sent to port 10000 are parsed with bounds checks and delivered to bounded subscribers. The media transport adapter is intentionally replaceable so a deployment can inject a Pion PeerConnection factory for ICE/DTLS/SRTP without coupling domain logic to a vendor implementation.

The static operations console is in `web/`; run `npm ci && npm run build` there when packaging a release. The Go service remains the API and serves the health/readiness endpoints used by deployment checks.

## Architecture

Domain, application, adapter and infrastructure layers are split under `internal/`. The in-memory repository is used for a self-contained local run; PostgreSQL migrations in `migrations/001_init.sql` define durable control-plane tables for production adapters. Context cancellation, wrapped errors, structured JSON logging, graceful shutdown and health/readiness/metrics endpoints are included.

## Verification

Run `gofmt -w .`, `go vet ./...`, `go test -race ./...`, and `go build ./...`. The repository deliberately has no generated code or frontend assets; all non-test behavior is Go source.
