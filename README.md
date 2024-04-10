# Learn Concurrent Programming With Go

## Troubleshootings

When running `go tool compile`, there is `could not import fmt (file not found)` error.
It can be fixed by running `GODEBUG=installgoroot=all go install std` before running `go tool compile` again.