package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/subscription"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	customerID := "dummy"
	priceID := "dummy"

	subParams := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(2),
			},
		},
		Metadata: map[string]string{
			"workspace_id": "ws_1234567890",
		},
	}
	_, err = subscription.New(subParams)
	if err != nil {
		log.Fatalf("Failed to create subscription: %v", err)
	}

	log.Println("Subscription product and subscription created successfully")
}
