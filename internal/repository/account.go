package repository

import "github.com/shopspring/decimal"

type AccountRepository interface {
	GetBalance(userId int) (decimal.Decimal, error)
}
