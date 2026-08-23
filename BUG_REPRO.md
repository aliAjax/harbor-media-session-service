# Bug Reproduction

## What happened

Snapshot helpers returned aliased slices instead of independent copies. Later caller or queue mutations changed stored manifests, room layers, signaling payloads, and drained audit details.

## How to trigger

Run the four Uniform copy checks against the original environment. Each test takes or stores a snapshot and then mutates the original or reuses the backing storage.

## Observed error

All four commands exited with status 1. Stored manifest data became `caller-mutated`, a room layer changed to `200`, a queued payload became `caller-mutated`, and a drained feedback slice was overwritten by subsequent queue pushes.
