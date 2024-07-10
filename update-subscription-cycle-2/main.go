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

	subscriptionID := "dummy"
	oldPriceID := "dummy"
	newPriceID := "dummy"

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

	params := &stripe.SubscriptionParams{
		TrialEnd:          stripe.Int64(nextFirstOfMonth.Unix()),
		ProrationBehavior: stripe.String("create_prorations"),
		Items: []*stripe.SubscriptionItemsParams{
			{
				ID:       stripe.String(itemID),
				Price:    stripe.String(newPriceID),
				Quantity: &quantity,
			},
		},
	}

	newSubscription, err := subscription.Update(subscriptionID, params)
	if err != nil {
		log.Fatalf("Failed to update subscription: %v", err)
	}

	fmt.Printf("Update subscription: %v\n", newSubscription.ID)
}
