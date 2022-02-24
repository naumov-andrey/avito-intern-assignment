package service

import (
	"github.com/naumov-andrey/avito-intern-assignment/internal/repository"
	"github.com/shopspring/decimal"
)

type AccountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo}
}

func (s *AccountService) GetBalance(userId int) (decimal.Decimal, error) {
	return s.repo.GetBalance(userId)
}
