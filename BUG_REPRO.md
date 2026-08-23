# Bug Reproduction

## What happened

Room lifecycle transition predicates used inconsistent state names and inverted several decisions. Open-to-closed and draining-to-closed were rejected, while draining-to-open and transitions out of closed were accepted.

## How to trigger

Run the four Sierra transition checks against the original environment. They cover active transitions, closed terminal behavior, draining transitions, and reopening a closed room.

## Observed error

All four commands exited with status 1. The checks failed at `score_test.go` lines 10, 15, 20, and 25, confirming that the allowed and rejected transition sets did not match the lifecycle contract.
