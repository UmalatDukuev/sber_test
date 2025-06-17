package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sber_test/internal/repo/cache"
	"time"
)

type Service struct {
	cache *cache.Cache
}

func New(c *cache.Cache) *Service {
	return &Service{cache: c}
}

type ExecuteRequest struct {
	ObjectCost     float64         `json:"object_cost"`
	InitialPayment float64         `json:"initial_payment"`
	Months         int             `json:"months"`
	Program        map[string]bool `json:"program"`
}

type Aggregates struct {
	Rate            int     `json:"rate"`
	LoanSum         float64 `json:"loan_sum"`
	MonthlyPayment  float64 `json:"monthly_payment"`
	Overpayment     float64 `json:"overpayment"`
	LastPaymentDate string  `json:"last_payment_date"`
}

type ExecuteResponse struct {
	Params struct {
		ObjectCost     float64 `json:"object_cost"`
		InitialPayment float64 `json:"initial_payment"`
		Months         int     `json:"months"`
	} `json:"params"`
	Program    map[string]bool `json:"program"`
	Aggregates Aggregates      `json:"aggregates"`
}

type CacheItem struct {
	ID int `json:"id"`
	ExecuteResponse
}

func loadProgramRates() (map[string]int, error) {
	absPath, err := filepath.Abs("programs.json")
	if err != nil {
		return nil, fmt.Errorf("unable to get absolute path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read programs.json: %w", err)
	}

	var result map[string]map[string]int
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return result["program_rates"], nil
}

func (s *Service) Execute(req ExecuteRequest) (ExecuteResponse, int, error) {
	programRates, err := loadProgramRates()
	if err != nil {
		return ExecuteResponse{}, 0, err
	}

	chosen := 0
	var annualRate int
	validPrograms := map[string]struct{}{
		"salary":   {},
		"military": {},
		"base":     {},
	}

	for k, v := range req.Program {
		if _, ok := validPrograms[k]; !ok {
			return ExecuteResponse{}, 0, errors.New("unknown program: " + k)
		}
		if v {
			chosen++
			if k == "salary" {
				annualRate = programRates["salary"]
			} else if k == "military" {
				annualRate = programRates["military"]
			} else if k == "base" {
				annualRate = programRates["base"]
			}
		}
	}

	if chosen == 0 {
		return ExecuteResponse{}, 0, errors.New("choose program")
	}
	if chosen > 1 {
		return ExecuteResponse{}, 0, errors.New("choose only 1 program")
	}

	if req.InitialPayment < 0.2*req.ObjectCost {
		return ExecuteResponse{}, 0, errors.New("the initial payment should be more")
	}

	loanSum := req.ObjectCost - req.InitialPayment
	r := float64(annualRate) / 12.0 / 100.0
	n := float64(req.Months)
	payment := loanSum * (r * math.Pow(1+r, n)) / (math.Pow(1+r, n) - 1)

	overpayment := (payment * n) - loanSum

	payment = math.Round(payment*100) / 100.0
	overpayment = math.Round(overpayment*100) / 100.0

	lastDate := time.Now().AddDate(0, req.Months, 0).Format("2006-01-02")

	var resp ExecuteResponse
	resp.Params.ObjectCost = req.ObjectCost
	resp.Params.InitialPayment = req.InitialPayment
	resp.Params.Months = req.Months
	resp.Program = req.Program
	resp.Aggregates = Aggregates{
		Rate:            annualRate,
		LoanSum:         loanSum,
		MonthlyPayment:  payment,
		Overpayment:     overpayment,
		LastPaymentDate: lastDate,
	}

	item := CacheItem{
		ID:              len(s.cache.GetAll()),
		ExecuteResponse: resp,
	}
	id := s.cache.Add(item)
	item.ID = id

	return resp, item.ID, nil
}

func (s *Service) GetAll() []CacheItem {
	raw := s.cache.GetAll()
	out := make([]CacheItem, 0, len(raw))
	for _, x := range raw {
		if ci, ok := x.(CacheItem); ok {
			out = append(out, ci)
		}
	}
	return out
}
