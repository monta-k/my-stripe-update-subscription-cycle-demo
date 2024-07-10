package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/price"
	"github.com/stripe/stripe-go/v79/product"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	prodParams := &stripe.ProductParams{
		Name: stripe.String("Update Cycle Test Product"),
		Type: stripe.String("service"),
	}
	prod, err := product.New(prodParams)
	if err != nil {
		log.Fatalf("Failed to create product: %v", err)
	}

	prices := []*stripe.PriceParams{
		{
			Product:    stripe.String(prod.ID),
			UnitAmount: stripe.Int64(1000),
			Currency:   stripe.String("jpy"),
			Recurring: &stripe.PriceRecurringParams{
				Interval:      stripe.String("day"),
				IntervalCount: stripe.Int64(30),
			},
		},
		{
			Product:    stripe.String(prod.ID),
			UnitAmount: stripe.Int64(1500),
			Currency:   stripe.String("jpy"),
			Recurring: &stripe.PriceRecurringParams{
				Interval: stripe.String("month"),
			},
		},
	}

	for _, priceParams := range prices {
		_, err = price.New(priceParams)
		if err != nil {
			log.Fatalf("Failed to create price: %v", err)
		}
	}

	log.Println("Subscription product created successfully")
}
