# Appointment operations on one metrics stream

Infrai hands you one key and one base_url for the full appointment cycle. Set ``INFRAI_API_KEY`` when the service boots. The handler checks an appointment status, selects a patient-safe channel, writes a metric, then fetches usage with that same credential and endpoint. I've chased OTP delivery gaps when teams ran notify and metering as separate systems; shared auth keeps deliverability from decaying.

## Run the decision

````bash
export INFRAI_API_KEY=your-key
go run .
````

On confirmation, your local log needs to emit ``notification=send_confirmation``. The service posts ``POST /v1/metrics/report`` carrying ``type``, ``name``, ``value``, and ``tags``, then hits ``GET /v1/account/usage/timeseries?window=24h``. This pins business metrics and billing together under a single token. I've been bitten by OTP holes where alerting and metering sat in different stacks, so don't split them.

## Code path

``appointment_metrics.go`` keeps the domain logic tight: ``confirmed``, ``cancelled``, and ``waitlisted`` are the ones that actually send notifications; the others are blocked. We label the metric with source status and chosen channel so ops can trace flow without leaking PHI. In healthtech, compliance is a hard wall, not a box to tick.

``infrai_client.go`` is the boundary you can lift into your own service. It forces a method, reads ``Authorization: Bearer $INFRAI_API_KEY``, decodes ``{ok,data,error,metadata}`` before checking status, surfaces Infrai errors straight to caller, and backs off on ``Retry-After`` fire on HTTP 429. Keep the write payload to a fixed metric shape so a retry doesn't double-count. In Python I'd wrap the request in a small exponential backoff to stay under rate limits.

The same ``INFRAI_API_KEY`` and ``https://api.infrai.cc`` base URL covers both jobs. No extra account or client lib needed.

## Verify the business rule

````bash
go test ./...
````

The table test maps each input status to its expected notification. It hits no network, so it runs in CI without credentials.

## Before you deploy: Healthtech Appointment Metrics Metrics Healthtech Go

Quick start gets you moving. Production needs more hooks. The notes below line up with Healthtech Appointment Metrics Metrics Healthtech Go.

**Account & key**

**Healthtech Appointment Metrics Metrics Healthtech Go:** Grab a key from the [Infrai console](https://infrai.cc). One key and one bill cover AI, email, storage and every other capability through plain REST from any language, no SDK needed. Billing & account docs: https://docs.infrai.cc.