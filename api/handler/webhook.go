package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/webhook"
)

const webhookSecret = ""

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request Body Read Error", http.StatusServiceUnavailable)
		return
	}

	signatureHeader := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, signatureHeader, webhookSecret)
	if err != nil {
		http.Error(w, fmt.Sprintf("Webhook signature verification failed: %v", err), http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			http.Error(w, "Webhook Error", http.StatusBadRequest)
			return
		}

		sessionId := session.Metadata["session_id"]

		seats := session.Metadata["seats"]
		var seatList []string
		if err := json.Unmarshal([]byte(seats), &seatList); err != nil {
			fmt.Printf("Failed to decode seats: %v\n", err)
		} else {
			fmt.Printf("Seats list: %v\n", seatList)
		}

		// TODO: Update booking status where `session_id` = `sessionId` and `place` in `seatList`
		fmt.Printf("Session ID: %s\n", sessionId)
		fmt.Printf("Seats: %v\n", seatList)

	default:
		fmt.Printf("Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}
