package repository

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AccountRepository interface {
	WithTx(tx *gorm.DB) AccountRepository
	GetBalance(userId int) (decimal.Decimal, error)
	UpdateBalance(userId int, balance decimal.Decimal) error
}
