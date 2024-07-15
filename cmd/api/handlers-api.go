package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/darkphotonKN/virtual-credit-card-server/internal/cards"
	"github.com/darkphotonKN/virtual-credit-card-server/internal/models"
	"github.com/darkphotonKN/virtual-credit-card-server/internal/product"
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
	err := json.NewDecoder(body).Decode(&payload) // reference to update payload variable this is required because passing it directly would be a copy and the original would be
	// unchanged

	if err != nil {
		app.errorLog.Println("Error when decoding json payload:", err)
		// TODO: send error response back to user
	}
	log.Println("body of request, decoded into payload:", payload)

	// get the amount passed in from payload
	amount := payload.Amount

	log.Println("Amount:", amount)

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

type ProductPurchaseResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    models.Product
}

// product handler
func (app *application) ProductPurchase(w http.ResponseWriter, r *http.Request) {

	var payload stripePayload

	// get body from request
	body := r.Body

	// decode from json into the stripePayload
	err := json.NewDecoder(body).Decode(&payload) // reference to update payload variable this is required because passing it directly would be a copy and the original would be
	// unchanged

	if err != nil {
		app.errorLog.Println("Error when decoding json payload:", err)
	}
	log.Println("body of request, decoded into payload:", payload)

	// get the amount passed in from payload
	amount := payload.Amount

	log.Println("Amount:", amount)

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

	if okay {
		newProduct := models.Product{
			Name:           "Default Product",
			InventoryLevel: "Default",
			Price:          2000, // default price
		}

		// create product purchase record
		product.CreateProductRecord(&app.DB, newProduct)

		createdProdRes := ProductPurchaseResponse{
			Status:  200,
			Message: "Successfully created product.",
			Data:    newProduct,
		}

		fmt.Println("Created product res:", createdProdRes)

		// return payment intent to user
		out, err := json.Marshal(pi)

		if err != nil {
			app.errorLog.Println("Error when converting amount (string) to amount (int).", err)
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

// get all product records
func (app *application) GetProductRecords(w http.ResponseWriter, r *http.Request) {

	// get the data from database
	products := product.GetProductRecords(&app.DB)

	// encode JSON
	productsJSON, err := json.Marshal(&products)

	if err != nil {
		fmt.Println("Error occured when attempting to encode json:", err)
	}

	w.Write(productsJSON)
}
