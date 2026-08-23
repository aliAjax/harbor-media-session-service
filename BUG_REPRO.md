# Bug Reproduction

## What happened

The Quebec state wrappers accessed shared values without synchronization. Concurrent reads and writes to changed, current, hysteresis, and metric state produced data races.

## How to trigger

Run the four dedicated Quebec checks with `go test -race` against the original environment. Each check concurrently reads and writes one state wrapper.

## Observed error

All four commands exited with status 1 and reported `DATA RACE`. The reports identified conflicting accesses in `changed_state_lock.go`, `current_state_lock.go`, `hysteresis_state_lock.go`, and `metric_snapshot_lock.go`.
