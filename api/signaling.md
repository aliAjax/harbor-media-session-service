# WebSocket signaling

Connect to `/v1/rooms/{room}/signal` and send JSON frames containing `type`, `id`, and monotonically increasing `seq`. Supported types are `join`, `offer`, `answer`, `candidate`, `publish`, `subscribe`, `unsubscribe`, `renegotiate`, and `leave`. Every accepted frame receives an `ack`; frames are broadcast to peers in the room. Frames larger than 1 MiB and stale sequence numbers are rejected.
