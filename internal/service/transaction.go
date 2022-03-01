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

func (s *TransactionService) Transfer(transferData model.TransferData) (model.TransferOutput, error) {
	out, err := s.repo.Atomic(func(r *repository.Repository) (interface{}, error) {
		creditUserBalance, err := s.repo.Account.GetBalance(transferData.CreditUserId)
		if err != nil {
			return model.TransferOutput{}, err
		}

		if creditUserBalance.LessThan(transferData.Amount) {
			return model.TransferOutput{}, errors.New("credit account balance is less than amount")
		}

		credit, err := r.Transaction.CreateTransaction(
			transferData.CreditUserId,
			transferData.Amount.Neg(),
			transferData.Description,
		)
		if err != nil {
			return model.TransferOutput{}, err
		}

		newCreditBalance := creditUserBalance.Sub(transferData.Amount)
		if err = r.Account.UpdateBalance(transferData.CreditUserId, newCreditBalance); err != nil {
			return model.TransferOutput{}, err
		}

		debit, err := r.Transaction.CreateTransaction(
			transferData.DebitUserId,
			transferData.Amount,
			transferData.Description,
		)
		if err != nil {
			return model.TransferOutput{}, err
		}

		debitUserBalance, err := s.repo.Account.GetBalance(transferData.DebitUserId)
		if err != nil {
			return model.TransferOutput{}, err
		}

		newDebitBalance := debitUserBalance.Add(transferData.Amount)
		if err = r.Account.UpdateBalance(transferData.DebitUserId, newDebitBalance); err != nil {
			return model.TransferOutput{}, err
		}

		return model.TransferOutput{Debit: debit, Credit: credit}, nil
	})

	if err != nil {
		return model.TransferOutput{}, err
	}

	return out.(model.TransferOutput), nil
}
