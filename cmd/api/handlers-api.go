package main

import (
	"encoding/json"
	"net/http"
)

// getting from client
type stripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

// sending as response
type jsonResponse struct {
	OK bool `json:"ok"`

	Message string `json:"mesage"`
	Content string `json:"content"`
	ID      int    `json:"id"`
}

// route handlers must take these two arguments to ve valid to pass to chi
func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {

	j := jsonResponse{
		OK: true,
	}

	marshalledJson, err := json.MarshalIndent(j, "", "		")
	if err != nil {
		app.errorLog.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(marshalledJson)

}
