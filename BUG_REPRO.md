# Bug Reproduction

## What happened

Error wrapping in the Tango transport and signaling paths did not preserve the underlying causes required by callers.

## How to trigger

Run the four dedicated Tango checks against the original environment. They exercise RTP parsing, frame reads, signaling dispatch, and UDP reads with concrete sentinel or network causes.

## Observed error

All four commands exited with status 1. The failures exposed `rtp packet too short`, `unexpected EOF`, `stale signaling sequence`, and `use of closed network connection` as causes that the returned errors did not preserve correctly.
