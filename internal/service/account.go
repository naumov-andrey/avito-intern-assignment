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

type AccountService struct {
	repo      repository.AccountRepository
	accessKey string
}

func NewAccountService(repo repository.AccountRepository, accessKey string) *AccountService {
	return &AccountService{repo, accessKey}
}

func (s *AccountService) GetBalance(userId int) (decimal.Decimal, error) {
	return s.repo.GetBalance(userId)
}

func (s *AccountService) GetConvertedBalance(userId int, currency string) (decimal.Decimal, error) {
	balance, err := s.repo.GetBalance(userId)
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
