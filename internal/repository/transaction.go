package repository

import (
	"github.com/naumov-andrey/avito-intern-assignment/internal/model"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	WithTx(tx *gorm.DB) TransactionRepository
	CreateTransaction(userId int, amount decimal.Decimal, description string) (model.Transaction, error)
	GetHistory(accountId, limit int, sortBy, orderBy string) ([]model.Transaction, error)
	GetHistoryWithCursor(accountId, limit, cursor int, sortBy, orderBy string) ([]model.Transaction, error)
}
