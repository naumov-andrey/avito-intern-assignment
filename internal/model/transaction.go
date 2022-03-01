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

type HistoryOutput struct {
	Data   []Transaction `json:"data"`
	Cursor int           `json:"cursor"`
}

type TransferInput struct {
	DebitUserId  int    `json:"debit_user_id" binding:"required"`
	CreditUserId int    `json:"credit_user_id" binding:"required"`
	Amount       string `json:"amount" binding:"required"`
	Description  string `json:"description"`
}

type TransferData struct {
	DebitUserId  int
	CreditUserId int
	Amount       decimal.Decimal
	Description  string
}

type TransferOutput struct {
	Debit  Transaction `json:"debit"`
	Credit Transaction `json:"credit"`
}
