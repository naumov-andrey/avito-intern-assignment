package postgres

import (
	"errors"
	"github.com/naumov-andrey/avito-intern-assignment/internal/model"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AccountRepositoryImpl struct {
	db *gorm.DB
}

func NewAccountRepositoryImpl(db *gorm.DB) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{db}
}

func (r *AccountRepositoryImpl) GetBalance(userId int) (decimal.Decimal, error) {
	tx := r.db.Begin()
	defer func() {
		if recover() != nil {
			tx.Rollback()
		}
	}()

	account := model.Account{UserId: userId}
	result := tx.First(&account, userId)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		tx.Select("UserId").Create(&account)
		tx.Commit()
	} else if result.Error != nil {
		tx.Rollback()
		return decimal.Decimal{}, result.Error
	}

	return account.Balance, nil
}
