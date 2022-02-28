package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	Id          int             `json:"id"`
	Date        time.Time       `json:"date"`
	AccountId   int             `json:"account_id"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
}
