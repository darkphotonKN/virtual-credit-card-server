// --- Handling all Stripe Card Actions
package cards

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

type Card struct {
	Secret   string
	Key      string
	Currency string
}

type Transaction struct {
	TransactionStatusID int
	Amount              int
	Currency            string
	LastFour            string
	BankReturnCode      string
}

func (c *Card) Charge(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.CreatePaymentIntent(currency, amount)
}

// Creating Payment Intent for Stripe
func (c *Card) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Secret

	// create payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}

	// add extra meta data for the transaction
	params.AddMetadata("key", "value")

	pi, err := paymentintent.New(params)

	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}

		// returns no payment intent, the error, and the error message
		return nil, msg, err
	}

	// returns payment intent, no error, empty message
	return pi, "", err
}

func cardErrorMessage(code stripe.ErrorCode) string {
	var msg = ""

	switch code {
	case stripe.ErrorCodeCardDeclined:
		msg = "Your card was declined."
	case stripe.ErrorCodeExpiredCard:
		msg = "Your card has expired."
	default:
		msg = "Your card was declined."
	}

	return msg

}
