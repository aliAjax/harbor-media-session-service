# Bug Reproduction

## Bug

Four media input paths read a slice before checking that enough data is present. Empty recording manifests and quality-layer lists panic instead of returning an error or fallback decision. Short RTP and RTCP packets panic instead of returning parsing errors.

## Trigger

Run the following commands against the buggy baseline:

```bash
go test ./internal/cinder -run '^TestCinderEmptyManifestReturnsError$' -count=1
go test ./internal/cinder -run '^TestCinderNoLayerUsesFallback$' -count=1
go test ./internal/cinder -run '^TestCinderShortRTPReturnsError$' -count=1
go test ./internal/cinder -run '^TestCinderShortRTCPReturnsError$' -count=1
```

## Error

All four commands exit with status 1. The observed failures are:

```text
empty manifest validation panicked: runtime error: index out of range [0] with length 0
empty quality decision panicked: runtime error: index out of range [0] with length 0
short RTP parsing panicked: runtime error: slice bounds out of range [:4] with capacity 1
short RTCP parsing panicked: runtime error: slice bounds out of range [:8] with capacity 4
```
