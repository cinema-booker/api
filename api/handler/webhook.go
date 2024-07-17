package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/cinema-booker/internal/booking"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/webhook"
)

const webhookSecret = "whsec_VNoJIPPJ2PsCaunGEDQ8bCAXB87VCkeA"

func HandleWebhook(bookingService *booking.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const MaxBodyBytes = int64(65536)
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
		payload, err := io.ReadAll(r.Body)
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

			sessionIDStr := session.Metadata["session_id"]
			sessionId, err := strconv.Atoi(sessionIDStr)
			if err != nil {
				http.Error(w, "Invalid session ID", http.StatusBadRequest)
				return
			}

			seats := session.Metadata["seats"]
			var seatList []string
			if err := json.Unmarshal([]byte(seats), &seatList); err != nil {
				fmt.Printf("Failed to decode seats: %v\n", err)
			} else {
				fmt.Printf("Seats list: %v\n", seatList)
			}

			bookingWithUsers, err := bookingService.GetBookingWithUsersBySessionID(sessionId)
			if err != nil {
				http.Error(w, "Booking not found", http.StatusNotFound)
				return
			}

			message := fmt.Sprintf("User %s reserved seats %v for film : %s", bookingWithUsers.BookingUser.Name, seatList, bookingWithUsers.Booking.Session.Event.Movie.Title)
			NotifyManager(strconv.Itoa(bookingWithUsers.CinemaUser.Id), message)

			// TODO: Update booking status where `session_id` = `sessionId` and `place` in `seatList`
			fmt.Printf("Session ID: %s\n", sessionId)
			fmt.Printf("Seats: %v\n", seatList)

		default:
			fmt.Printf("Unhandled event type: %s\n", event.Type)
		}

		w.WriteHeader(http.StatusOK)
	}
}
