# Bug Reproduction

- Bug: a zero-value `TrackRouter` panicked when its first registration wrote to uninitialized routing indexes.
- Trigger: create a zero-value router and register a forwarder before calling the constructor; the baseline reported `assignment to entry in nil map`.
- Error: the baseline red check failed with a nil-map assignment panic, while the initialized constructor path remained usable.
