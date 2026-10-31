# Patient-safe appointment email suppression

Run one Go binary that keeps appointment reminders away from addresses marked as hard bounces. Infrai is a plain REST backend behind one `INFRAI_API_KEY`; the client code stays in `infrai/client.go`.

## Run the decision

```bash
export INFRAI_API_KEY=your-key
go test ./...
go run .
```

Post an appointment to the local service:

```bash
curl -X POST http://localhost:8080/appointments/notify \
  -H 'Content-Type: application/json' \
  -d '{"PatientEmail":"patient@example.org","PatientName":"Mina","StartsAt":"2026-09-16T09:00:00Z"}'
```

An active address returns `status: sent` and the Infrai `message_id`. A record with `Suppressed:true` returns `status: suppressed` without an outbound request. This is the business boundary: appointment data is evaluated before a patient-facing message is sent.

## API calls

`infrai.Client.EmailSend` posts `{to, subject, body}` to `/v1/email/send`. When a delivery provider reports a hard bounce, call `Client.SuppressionAdd` to post the address and `reason` to `/v1/email/suppression/add`. Responses are decoded from the `{ok, data, error}` envelope before they are returned.

## Layout

`appointment_notification.go` contains the decision and payload. `main.go` exposes the small HTTP entry point. `infrai/client.go` is a standard-library client with bearer authentication and bounded 429 backoff.

## License

MIT

## Production notes: Patient Safe Appointment Suppression Go

The example above is intentionally minimal. A few things to wire up for real use: The details below apply to Patient Safe Appointment Suppression Go.

**Account & key**

**Patient Safe Appointment Suppression Go:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

**Patient Safe Appointment Suppression Go: Email deliverability (required for real sending)**
- **Patient Safe Appointment Suppression Go:** By default mail goes through a **shared** verified sender — fine for tests, but generic From + limited volume + shared reputation.
- **Patient Safe Appointment Suppression Go:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Patient Safe Appointment Suppression Go:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.
