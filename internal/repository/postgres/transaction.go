package postgres

import (
	"fmt"
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

func (r *TransactionRepositoryImpl) GetHistory(
	accountId int,
	limit int,
	sortBy string,
	orderBy string,
) ([]model.Transaction, error) {
	tx := r.db.Begin()
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	data := make([]model.Transaction, 0)
	result := tx.
		Limit(limit).
		Where("account_id = ?", accountId).
		Order(sortBy + " " + orderBy).
		Order("id " + orderBy).
		Find(&data)

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	return data, tx.Commit().Error
}

func (r *TransactionRepositoryImpl) GetHistoryWithCursor(
	accountId int,
	limit int,
	cursor int,
	sortBy string,
	orderBy string,
) ([]model.Transaction, error) {
	tx := r.db.Begin()
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	var cursorTransaction model.Transaction
	tx.Select(sortBy).First(&cursorTransaction, cursor)

	if tx.Error != nil {
		tx.Rollback()
		return nil, tx.Error
	}

	sign := "<"
	if orderBy == "asc" {
		sign = ">"
	}
	whereCondition := fmt.Sprintf("(%s, id) %s (?, ?)", sortBy, sign)

	var cursorData interface{}
	if sortBy == "date" {
		cursorData = cursorTransaction.Date
	} else {
		// sort by amount
		cursorData = cursorTransaction.Amount
	}

	data := make([]model.Transaction, 0)
	result := tx.
		Limit(limit).
		Where("account_id = ?", accountId).
		Where(whereCondition, cursorData, cursor).
		Order(sortBy + " " + orderBy).
		Order("id " + orderBy).
		Find(&data)

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	return data, tx.Commit().Error
}
