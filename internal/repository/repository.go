package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	Account     AccountRepository
	Transaction TransactionRepository
	db          *gorm.DB
}

func (r *Repository) Atomic(fn func(r *Repository) (interface{}, error)) (out interface{}, err error) {
	tx := r.db.Begin()

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	out, err = fn(r)
	return
}

func NewRepository(account AccountRepository, transaction TransactionRepository, db *gorm.DB) Repository {
	return Repository{account, transaction, db}
}
