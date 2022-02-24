package model

import (
	"github.com/shopspring/decimal"
)

type Account struct {
	UserId  int             `json:"user_id"`
	Balance decimal.Decimal `json:"balance"`
}
