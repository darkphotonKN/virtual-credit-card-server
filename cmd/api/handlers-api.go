package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/darkphotonKN/virtual-credit-card-server/internal/cards"
)

// getting from client
type stripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

// sending as response
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

// route handlers must take these two arguments to ve valid to pass to chi
func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload

	// get body from request
	body := r.Body

	// decode from json into the stripePayload
	err := json.NewDecoder(body).Decode(&payload) // reference to update payload variable
	// this is required because passing it directly would be a copy and the original would be
	// unchanged

	if err != nil {
		app.errorLog.Println("Error when decoding json payload:", err)
		// TODO: send error response back to user
	}
	log.Println("body of request, decoded into payload:", payload)

	// get the amount passed in from payload
	amount := payload.Amount

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: payload.Currency,
	}

	okay := true

	pi, msg, err := card.Charge(payload.Currency, amount)

	if err != nil {
		okay = false
	}

	w.Header().Set("Content-Type", "application/json")

	if okay {
		out, err := json.Marshal(pi)

		if err != nil {
			app.errorLog.Println("Error when converting amount (string) to amount (int).", err)
			// TODO: Add error response
		}

		// Constructing Response - Returning Payment Intent
		w.Write(out)
	} else {
		errorJ := jsonResponse{
			OK:      false,
			Message: msg,
			Content: "",
		}
		out, err := json.Marshal(errorJ)

		if err != nil {
			app.errorLog.Println("Error when converting error message to json.")
		}
		w.Write(out)
	}

}
