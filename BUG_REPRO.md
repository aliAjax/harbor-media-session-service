# Bug Reproduction

## Bug

Concurrent packet ingestion and media-session cleanup access four mutable objects without synchronized read/write lifetimes: the SSRC remapper, jitter buffer, RTCP controller snapshot state, and participant session liveness state.

## Trigger

Run the following commands against the buggy baseline:

```bash
go test -race ./internal/xray -run '^TestXrayRemapperConcurrentReset$' -count=1
go test -race ./internal/xray -run '^TestXrayJitterConcurrentAccess$' -count=1
go test -race ./internal/xray -run '^TestXrayRTCPConcurrentSnapshot$' -count=1
go test -race ./internal/xray -run '^TestXrayParticipantSessionConcurrentClose$' -count=1
```

## Error

All four commands exit with status 1 and report `WARNING: DATA RACE`. The RTCP snapshot case can additionally terminate the process with:

```text
fatal error: concurrent map iteration and map write
```
