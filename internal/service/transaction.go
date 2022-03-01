package service

import (
	"errors"
	"github.com/naumov-andrey/avito-intern-assignment/internal/model"
	"github.com/naumov-andrey/avito-intern-assignment/internal/repository"
)

type TransactionService struct {
	repo repository.Repository
}

func NewTransactionService(repo repository.Repository) *TransactionService {
	return &TransactionService{repo}
}

func (s *TransactionService) GetHistory(
	userId int,
	limit int,
	cursor int,
	sortBy string,
	orderBy string,
) (model.HistoryOutput, error) {
	if sortBy != "date" && sortBy != "amount" {
		return model.HistoryOutput{}, errors.New("sort must be by date or amount")
	}

	if orderBy != "asc" && orderBy != "desc" {
		return model.HistoryOutput{}, errors.New("order must be asc or desc")
	}

	var (
		data []model.Transaction
		err  error
	)
	if cursor != -1 {
		data, err = s.repo.Transaction.GetHistoryWithCursor(userId, limit, cursor, sortBy, orderBy)
	} else {
		data, err = s.repo.Transaction.GetHistory(userId, limit, sortBy, orderBy)
	}

	if err != nil {
		return model.HistoryOutput{}, err
	}

	if len(data) != 0 && limit == len(data) {
		cursor = data[len(data)-1].Id
	} else {
		cursor = -1
	}

	return model.HistoryOutput{Data: data, Cursor: cursor}, nil
}
