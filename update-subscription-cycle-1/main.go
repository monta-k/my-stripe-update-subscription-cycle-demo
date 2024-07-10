package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

	subscriptionID := os.Getenv("STRIPE_SUBSCRIPTION_ID")
	oldPriceID := os.Getenv("STRIPE_OLD_PRICE_ID")
	newPriceID := os.Getenv("STRIPE_NEW_PRICE_ID")

	now := time.Now()
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("Failed to load JST location: %v", err)
	}
	nextFirstOfMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, jst)

	sub, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		log.Fatalf("Failed to retrieve subscription: %v", err)
	}

	var itemID string
	var quantity int64
	for _, item := range sub.Items.Data {
		if item.Price.ID == oldPriceID {
			itemID = item.ID
			quantity = item.Quantity
			break
		}
	}
	if itemID == "" {
		log.Fatalf("Failed to find item ID: %v", err)
	}

	_, err = subscription.Cancel(subscriptionID, &stripe.SubscriptionCancelParams{
		InvoiceNow: stripe.Bool(true),
		Prorate:    stripe.Bool(true),
	})
	if err != nil {
		log.Fatalf("Failed to cancel subscription: %v", err)
	}

	params := &stripe.SubscriptionParams{
		Customer:           stripe.String(sub.Customer.ID),
		BillingCycleAnchor: stripe.Int64(nextFirstOfMonth.Unix()),
		ProrationBehavior:  stripe.String("create_prorations"),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price:    stripe.String(newPriceID),
				Quantity: &quantity,
			},
		},
	}

	newSubscription, err := subscription.New(params)
	if err != nil {
		log.Fatalf("Failed to create new subscription: %v", err)
	}

	fmt.Printf("Created new subscription: %v\n", newSubscription.ID)
}
