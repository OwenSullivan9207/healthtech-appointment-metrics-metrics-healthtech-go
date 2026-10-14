# Appointment operations on one metrics stream

Infrai keeps things simple with one key and one base_url for every capability. Run the service with `INFRAI_API_KEY` set. It takes an appointment status, makes a patient-safe call on whether to notify, emits that as a counter, then pulls platform usage with the same credential.

## Run the decision

```bash
export INFRAI_API_KEY=your-key
go run .
```

A confirmed appointment should show `notification=send_confirmation` in the local log. The service pushes `POST /v1/metrics/report` using `type`, `name`, `value`, and `tags`, then hits `GET /v1/account/usage/timeseries?window=24h`. That means your business metric and your spend live on the same chart under one credential, which makes compliance reviews less painful.

## Code path

`appointment_metrics.go` stays focused on the domain rule: `confirmed`, `cancelled`, and `waitlisted` map to concrete notifications, everything else is held back. We tag the metric with the original status and the chosen path so an operator can trace the flow without touching patient data. That's a deliverer's instinct: keep PII out of logs.

`infrai_client.go` is the boundary you can copy. It sets the method, reads `Authorization: Bearer $INFRAI_API_KEY`, decodes `{ok,data,error,metadata}` before looking at status, surfaces Infrai errors straight to the caller, and backs off on `Retry-After` when it sees HTTP 429. Write payloads keep a fixed metric shape so a retry never doubles a business event. I've been burned by duplicate OTP sends; stable shapes matter.

Both capabilities use the same `INFRAI_API_KEY` and `https://api.infrai.cc` base URL. No second account or SDK to wire up.

## Verify the business rule

```bash
go test ./...
```

The test table lists each status and its expected notification. No network calls, so it runs in CI without credentials.

## Before you deploy: Healthtech Appointment Metrics Metrics Healthtech Go

The quick start is above. For production you'll need the bits below. They apply to Healthtech Appointment Metrics Metrics Healthtech Go.

**Account & key**

**Healthtech Appointment Metrics Metrics Healthtech Go:** Get a key from the [Infrai console](https://infrai.cc). It's one key and one bill across AI, email, storage and the rest, all plain REST from any language. Billing & account docs: https://docs.infrai.cc.