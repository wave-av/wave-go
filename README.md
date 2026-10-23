# wave-go

The official WAVE SDK for Go. Video infrastructure for people and AI agents: one API for live and <!-- # guard:allow plain-English word, not a repo reference -->
on-demand video, agent-payable over x402.

## Install

```bash
go get github.com/wave-av/wave-go
```

## Use

```go
client := waveav.NewClient("wav_live_...") // from https://console.wave.online

models, err := client.ListModels(ctx)
completion, err := client.Complete(ctx, waveav.CompletionRequest{
    Model:    "deepseek-v4",
    Messages: []waveav.Message{{Role: "user", Content: "hello"}},
})
```

Every method maps to one documented gateway route (https://dev.wave.online/reference). Errors are
`*waveav.APIError` carrying the gateway's typed codes (`AUTH_MISSING`, `SCOPE_REQUIRED`,
`PAYMENT_REQUIRED`).

## The surfaces

- `ListModels` / `Complete`: the inference funnel (measured routing, 13 providers, per-token metering)
- `Usage`: the platform's live usage snapshot

More modules follow the same one-method-one-route pattern.
## Buildkite runbook

`guard-content-policy` refuses a build when the normalized repository
identities differ: it compares host/owner/name (lowercased, scheme and
`.git` suffix agnostic), so SSH-vs-HTTPS spellings of the same
repository proceed while any other mismatch \u2014 including malformed or
empty identities \u2014 refuses. Secrets-queue steps consume host
credentials, so refusal is fail-closed (exit 2).
Trusted re-run procedure for a fork PR: run the pipeline from the
trusted base ref (a base-branch build) — never by re-running the fork
build itself. A re-run of the fork build keeps fork metadata
(`BUILDKITE_PULL_REQUEST_REPO` still points at the fork), so the guard
still refuses it; only a build whose head is the trusted base proceeds.
The host withholds secret injection from any fork build, and a manual
approval gate is required before any trusted re-run. Never approve a
re-run that executes PR-head code with host credentials.
