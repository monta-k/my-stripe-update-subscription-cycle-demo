package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/subscription"
	"github.com/stripe/stripe-go/v79/subscriptionschedule"
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

	// 次の1日に設定
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

	// 新しいサブスクリプションスケジュールを作成
	params := &stripe.SubscriptionScheduleParams{
		FromSubscription: stripe.String(subscriptionID),
	}

	schedule, err := subscriptionschedule.New(params)
	if err != nil {
		log.Fatalf("Failed to create subscription schedule: %v", err)
	}

	// サブスクリプションスケジュールに現在のサブスクリプションをリンク
	_, err = subscriptionschedule.Update(schedule.ID, &stripe.SubscriptionScheduleParams{
		Phases: []*stripe.SubscriptionSchedulePhaseParams{
			{
				Items: []*stripe.SubscriptionSchedulePhaseItemParams{
					{
						Price:    stripe.String(oldPriceID),
						Quantity: stripe.Int64(quantity),
					},
				},
				StartDate:         stripe.Int64(sub.CurrentPeriodStart),
				EndDate:           stripe.Int64(nextFirstOfMonth.Unix()),
				ProrationBehavior: stripe.String("create_prorations"),
			},
			{
				Items: []*stripe.SubscriptionSchedulePhaseItemParams{
					{
						Price:    stripe.String(newPriceID),
						Quantity: stripe.Int64(quantity),
					},
				},
				StartDate:         stripe.Int64(nextFirstOfMonth.Unix()),
				ProrationBehavior: stripe.String("create_prorations"),
			},
		},
		EndBehavior: stripe.String("release"),
	})
	if err != nil {
		log.Fatalf("Failed to link subscription to schedule: %v", err)
	}

	fmt.Printf("Created new subscription schedule for subscription %v: %v\n", subscriptionID, schedule.ID)

}
