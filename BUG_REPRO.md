# Bug Reproduction

## What happened

The room service and its in-memory repository accepted already-canceled contexts. A canceled create could leave a room behind, reads still returned stored rooms, and a canceled transition could change room state.

## How to trigger

Create an already-canceled context and pass it through the repository and service create, get, and transition paths. The red evidence runs the four dedicated checks against the original environment; the repaired environment passes the same checks.

## Observed error

The red checks reported `context.Canceled` was expected but got `<nil>` for repository create and service transition. Repository get and service create returned a room and `<nil>` instead of `(nil, context.Canceled)`.
