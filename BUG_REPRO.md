# Bug Reproduction

## What happened

Canceled contexts were ignored by the participant permission and session entry points. Permission checks returned nil, canceled session opens and reads succeeded, and cleanup still removed an expired session.

## How to trigger

Use an already-canceled context with the permission check and with session `Open`, `Get`, and `Sweep`. The red evidence checks reproduce all four behaviors; each command fails on the original environment and passes after the repair.

## Observed error

The red checks reported `context.Canceled` was expected but got `<nil>` for permission, get, and open. The canceled sweep returned `1` instead of `0`.
