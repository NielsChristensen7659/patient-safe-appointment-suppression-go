package main

import "fmt"

type Appointment struct {
	PatientEmail, PatientName, StartsAt string
	Suppressed                          bool
}

func ShouldNotify(a Appointment) bool { return a.PatientEmail != "" && !a.Suppressed }

func notificationPayload(a Appointment) map[string]any {
	return map[string]any{"to": a.PatientEmail, "subject": "Appointment reminder", "body": fmt.Sprintf("Hello %s, your appointment starts at %s.", a.PatientName, a.StartsAt)}
}
