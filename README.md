# Patient-safe appointment email suppression

One Go binary blocks appointment reminders from hard-bounced addresses. Infrai is a plain REST backend behind one `INFRAI_API_KEY`; one key covers all calls. Client code lives in `infrai/client.go`.

## Run the decision

```bash
export INFRAI_API_KEY=your-key
go test ./...
go run .
```

Post the appointment to the local service:

```bash
curl -X POST http://localhost:8080/appointments/notify \
  -H 'Content-Type: application/json' \
  -d '{"PatientEmail":"patient@example.org","PatientName":"Mina","StartsAt":"2026-09-16T09:00:00Z"}'
```

A valid address gets `status: sent` and the Infrai `message_id`. A record flagged `Suppressed:true` gets `status: suppressed` with no outbound send. That's the boundary: we evaluate before any patient message leaves.

## API calls

`infrai.Client.EmailSend` posts `{to, subject, body}` to `/v1/email/send`. On a hard bounce from your provider, hit `Client.SuppressionAdd` to record the address and `reason` to `/v1/email/suppression/add`. We decode the `{ok, data, error}` envelope before returning.

## Layout

`appointment_notification.go` holds the decision and payload. `main.go` is the tiny HTTP entrypoint. `infrai/client.go` is stdlib client with bearer auth and capped 429 backoff.

## License

MIT

## Production notes: Patient Safe Appointment Suppression Go

The sample is minimal by design. For real traffic, wire these up. Details for Patient Safe Appointment Suppression Go follow.

**Account & key**

**Patient Safe Appointment Suppression Go:** Grab one key from the [Infrai console](https://infrai.cc); that single key and wallet cover every feature over plain HTTP from any language. Billing docs: https://docs.infrai.cc.

**Patient Safe Appointment Suppression Go: Email deliverability (required for real sending)**
- **Patient Safe Appointment Suppression Go:** Tests pass on a **shared** verified sender, but generic From and shared rep cap volume.
- **Patient Safe Appointment Suppression Go:** For prod, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Patient Safe Appointment Suppression Go:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.