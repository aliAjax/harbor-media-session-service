# Bug Reproduction

## Bug

After a lease conflict, duplicate session, quota overflow, or duplicate subscriber error, the affected in-memory manager remains blocked. A later valid release or create operation on the same manager never completes.

## Trigger

Run these checks against the buggy baseline:

```bash
go test ./internal/zulu -run '^TestZuluLeaseConflictUnlocksStore$' -count=1
go test ./internal/zulu -run '^TestZuluDuplicateSessionUnlocksManager$' -count=1
go test ./internal/zulu -run '^TestZuluQuotaFailureUnlocksManager$' -count=1
go test ./internal/zulu -run '^TestZuluDuplicateSubscriberUnlocksForwarder$' -count=1
```

## Error

All four commands exit with status 1. The observed failures report that the follow-up operation remained blocked after the preceding error:

```text
lease release remained blocked after an error
next session open remained blocked after an error
quota release remained blocked after an error
next subscriber add remained blocked after an error
```
