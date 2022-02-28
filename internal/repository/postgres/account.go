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
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
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

func (r *AccountRepositoryImpl) UpdateBalance(userId int, balance decimal.Decimal) error {
	return r.db.Model(&model.Account{}).Where("user_id = ?", userId).Update("balance", balance).Error
}
