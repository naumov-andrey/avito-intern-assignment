package model

import (
	"github.com/shopspring/decimal"
)

type Account struct {
	UserId       int             `json:"-"`
	Balance      decimal.Decimal `json:"balance"`
	Transactions []Transaction   `json:"-" gorm:"foreignKey:AccountId;references:UserId"`
}

type UpdateBalanceInput struct {
	Amount      string `json:"amount" binding:"required"`
	Description string `json:"description"`
}
