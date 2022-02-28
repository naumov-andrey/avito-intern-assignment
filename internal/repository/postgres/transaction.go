package postgres

import (
	"github.com/naumov-andrey/avito-intern-assignment/internal/model"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type TransactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepositoryImpl(db *gorm.DB) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{db}
}

func (r *TransactionRepositoryImpl) CreateTransaction(
	userId int,
	amount decimal.Decimal,
	description string,
) (model.Transaction, error) {
	t := model.Transaction{
		Date:        time.Now(),
		AccountId:   userId,
		Amount:      amount,
		Description: description,
	}

	result := r.db.Create(&t)
	if result.Error != nil {
		return model.Transaction{}, result.Error
	}

	return t, nil
}
