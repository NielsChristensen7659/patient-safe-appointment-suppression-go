package main

import "testing"

func TestShouldNotify(t *testing.T) {
	tests := []struct {
		name string
		a    Appointment
		want bool
	}{
		{"active appointment", Appointment{PatientEmail: "patient@example.org", PatientName: "Mina", StartsAt: "2026-09-16T09:00:00Z"}, true},
		{"suppressed bounce", Appointment{PatientEmail: "patient@example.org", Suppressed: true}, false},
		{"missing address", Appointment{PatientName: "Mina"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShouldNotify(tt.a); got != tt.want {
				t.Fatalf("ShouldNotify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotificationPayload(t *testing.T) {
	payload := notificationPayload(Appointment{PatientEmail: "chenhua@changba.com", PatientName: "Mina", StartsAt: "2026-09-16T09:00:00Z"})
	if _, ok := payload["text"]; ok {
		t.Fatal("email.send does not accept text")
	}
	if got := payload["body"]; got != "Hello Mina, your appointment starts at 2026-09-16T09:00:00Z." {
		t.Fatalf("body = %v", got)
	}
}
