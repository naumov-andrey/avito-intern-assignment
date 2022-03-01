package service

import (
	"errors"
	"fmt"
	"github.com/naumov-andrey/avito-intern-assignment/internal/repository"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
)

type BalanceService struct {
	repo      repository.Repository
	accessKey string
}

func NewBalanceService(repo repository.Repository, accessKey string) *BalanceService {
	return &BalanceService{repo, accessKey}
}

func (s *BalanceService) GetBalance(userId int) (decimal.Decimal, error) {
	return s.repo.Account.GetBalance(userId)
}

func (s *BalanceService) GetConvertedBalance(userId int, currency string) (decimal.Decimal, error) {
	balance, err := s.repo.Account.GetBalance(userId)
	if err != nil {
		return balance, err
	}

	apiCall := fmt.Sprintf("https://freecurrencyapi.net/api/v2/latest?apikey=%s&base_currency=RUB", s.accessKey)
	response, err := http.Get(apiCall)
	if err != nil {
		return balance, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return balance, err
	}

	rateStr := gjson.Get(string(responseData), "data."+currency).String()
	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		return balance, errors.New(fmt.Sprintf("Error while getting RUB/%s currecny rate", currency))
	}

	return balance.Mul(decimal.NewFromFloat(rate)), nil
}

func (s *BalanceService) UpdateBalance(userId int, amount decimal.Decimal, description string) (decimal.Decimal, error) {
	if amount.Equal(decimal.Decimal{}) {
		return s.repo.Account.GetBalance(userId)
	}

	out, err := s.repo.Atomic(func(r *repository.Repository) (interface{}, error) {
		var zero decimal.Decimal
		balance, err := r.Account.GetBalance(userId)
		if err != nil {
			return zero, err
		}

		newBalance := balance.Add(amount)
		if newBalance.LessThan(zero) {
			return zero, errors.New("account balance is less than credit amount")
		}

		err = r.Account.UpdateBalance(userId, newBalance)
		if err != nil {
			return zero, err
		}

		if _, err = r.Transaction.CreateTransaction(userId, amount, description); err != nil {
			return zero, err
		}

		balance, err = r.Account.GetBalance(userId)
		if err != nil {
			return zero, err
		}
		return balance, err
	})

	if err != nil {
		return decimal.Decimal{}, err
	}

	return out.(decimal.Decimal), nil
}
