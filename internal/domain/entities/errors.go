package entities

import "errors"

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrInvalidCurrency     = errors.New("invalid currency")
	ErrAccountNotFound     = errors.New("account not found")
	ErrUserNotFound        = errors.New("user not found")
)
