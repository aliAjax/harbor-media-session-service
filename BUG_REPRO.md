# Bug Reproduction

- Bug: canceled lease claims and renewals still changed cluster lease state; expired leases could be renewed; a release from the wrong owner could delete an active lease.
- Trigger: exercise the lease coordinator with an already-canceled context, renew a lease after its expiry, or release a lease using a different owner while retaining the token.
- Error: the baseline verification reported failed assertions for canceled claim/renew, expired renewal, and forged-owner release; the underlying cancellation error was lost and stale lease operations were accepted.
