package postgres

import (
	"errors"
	"github.com/naumov-andrey/avito-intern-assignment/internal/model"
	"github.com/naumov-andrey/avito-intern-assignment/internal/repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AccountRepositoryImpl struct {
	db *gorm.DB
}

func NewAccountRepositoryImpl(db *gorm.DB) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{db}
}

func (r *AccountRepositoryImpl) WithTx(tx *gorm.DB) repository.AccountRepository {
	newRepo := *r
	newRepo.db = tx
	return &newRepo
}

func (r *AccountRepositoryImpl) GetBalance(userId int) (decimal.Decimal, error) {
	account := model.Account{UserId: userId}
	result := r.db.First(&account, userId)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		r.db.Select("UserId").Create(&account)
	} else if result.Error != nil {
		return decimal.Decimal{}, result.Error
	}

	return account.Balance, nil
}

func (r *AccountRepositoryImpl) UpdateBalance(userId int, balance decimal.Decimal) error {
	return r.db.Model(&model.Account{}).Where("user_id = ?", userId).Update("balance", balance).Error
}
