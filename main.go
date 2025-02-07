package main

import (
	"errors"
	"fmt"
)

// Database interface 
type Database interface {
	GetUserBalance(userID int) (float64, error)
	UpdateUserBalance(userID int, newBalance float64) error
}

// PaymentGateway interface 
type PaymentGateway interface {
	ProcessPayment(userID int, amount float64) (bool, error)
}

// DiscountCalculator struct 
type DiscountCalculator struct{}

func (dc *DiscountCalculator) CalculateDiscount(amount float64) float64 {
	if amount > 100 {
		return amount * 0.1 // 10% discount
	}
	return 0
}

// OrderService (Main Business Logic)
type OrderService struct {
	db     Database
	pg     PaymentGateway
	disc   *DiscountCalculator
}

func (s *OrderService) PlaceOrder(userID int, orderAmount float64) (string, error) {
	// Step 1: Get user balance from database
	balance, err := s.db.GetUserBalance(userID)
	if err != nil {
		return "", errors.New("failed to fetch user balance")
	}

	// Step 2: Apply discount
	discount := s.disc.CalculateDiscount(orderAmount)
	finalAmount := orderAmount - discount

	// Step 3: Check if user has enough balance
	if balance < finalAmount {
		return "", errors.New("insufficient balance")
	}

	// Step 4: Process payment via Payment Gateway
	success, err := s.pg.ProcessPayment(userID, finalAmount)
	if err != nil || !success {
		return "", errors.New("payment failed")
	}

	// Step 5: Update user balance in database
	err = s.db.UpdateUserBalance(userID, balance-finalAmount)
	if err != nil {
		return "", errors.New("failed to update balance")
	}

	return fmt.Sprintf("Order placed successfully! Final Amount: %.2f", finalAmount), nil
}
