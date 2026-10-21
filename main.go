package main

import (
	"encoding/json"
	"log"
	"net/http"

	"appointment-suppression/infrai"
)

func main() {
	c, err := infrai.New()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/appointments/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var a Appointment
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			http.Error(w, "invalid appointment", http.StatusBadRequest)
			return
		}
		if !ShouldNotify(a) {
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "suppressed"})
			return
		}
		data, err := c.EmailSend(notificationPayload(a))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "sent", "message_id": data["message_id"]})
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
