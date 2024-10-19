package models

import (
	"github.com/shopspring/decimal"
)

type PaymentInput struct {
	Amount decimal.Decimal `json:"amount" binding:"required"`
	Type   string          `json:"type"`
}
