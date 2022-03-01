package repository

import (
	"github.com/naumov-andrey/avito-intern-assignment/internal/model"
	"github.com/shopspring/decimal"
)

type TransactionRepository interface {
	CreateTransaction(userId int, amount decimal.Decimal, description string) (model.Transaction, error)
	GetHistory(accountId, limit int, sortBy, orderBy string) ([]model.Transaction, error)
	GetHistoryWithCursor(accountId, limit, cursor int, sortBy, orderBy string) ([]model.Transaction, error)
}
