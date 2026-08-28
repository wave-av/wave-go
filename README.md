# wave-go

The official WAVE SDK for Go. Video infrastructure for people and AI agents: one API for live and
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
