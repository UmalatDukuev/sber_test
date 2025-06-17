// Package service contains the business logic for handling loan calculations,
package service

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sber_test/internal/repo/cache"
	"time"
)

// Service handles loan calculations and caching.
type Service struct {
	cache *cache.Cache
}

// New creates a new Service instance.
func New(c *cache.Cache) *Service {
	return &Service{cache: c}
}

// ExecuteRequest contains parameters for loan calculation.
type ExecuteRequest struct {
	Program        map[string]bool `json:"program"`
	ObjectCost     float64         `json:"object_cost"`
	InitialPayment float64         `json:"initial_payment"`
	Months         int             `json:"months"`
}

// Aggregates holds the results of loan calculations.
type Aggregates struct {
	LastPaymentDate string  `json:"last_payment_date"`
	Rate            int     `json:"rate"`
	LoanSum         float64 `json:"loan_sum"`
	MonthlyPayment  float64 `json:"monthly_payment"`
	Overpayment     float64 `json:"overpayment"`
}

// ExecuteResponse contains the result of loan calculation.
type ExecuteResponse struct {
	Program    map[string]bool `json:"program"`
	Aggregates Aggregates      `json:"aggregates"`
	Params     struct {
		ObjectCost     float64 `json:"object_cost"`
		InitialPayment float64 `json:"initial_payment"`
		Months         int     `json:"months"`
	} `json:"params"`
}

// CacheItem stores the loan calculation result and its ID.
type CacheItem struct {
	ExecuteResponse
	ID int `json:"id"`
}

func loadProgramRates() (map[string]int, error) {
	absPath, err := filepath.Abs("programs.json")
	if err != nil {
		return nil, fmt.Errorf("unable to get absolute path: %w", err)
	}

	cleanPath := filepath.Clean(absPath)

	data, err := os.ReadFile(cleanPath)
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

// Execute - adding and calculating new credit.
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
			return ExecuteResponse{}, 0, fmt.Errorf("%w: %s", ErrUnknownProgram, k)
		}
		if v {
			chosen++
			switch k {
			case "salary":
				annualRate = programRates["salary"]
			case "military":
				annualRate = programRates["military"]
			case "base":
				annualRate = programRates["base"]
			}
		}
	}
	if req.InitialPayment >= req.ObjectCost {
		return ExecuteResponse{}, 0, fmt.Errorf("%w: %f >= %f", ErrFirstPaymentExceedsLoan, req.InitialPayment, req.ObjectCost)
	}
	if chosen == 0 {
		return ExecuteResponse{}, 0, ErrChooseProgram
	}
	if chosen > 1 {
		return ExecuteResponse{}, 0, ErrChooseOnlyOneProgram
	}
	if req.InitialPayment < 0.2*req.ObjectCost {
		return ExecuteResponse{}, 0, ErrInitialPaymentLow
	}
	loanSum, payment, overpayment, lastDate := calculateCredit(req, annualRate)
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
	return resp, id, nil
}

func calculateCredit(req ExecuteRequest, annualRate int) (loanSum, payment, overpayment float64, lastDate string) {
	loanSum = req.ObjectCost - req.InitialPayment
	r := float64(annualRate) / 12.0 / 100.0
	n := float64(req.Months)
	payment = loanSum * (r * math.Pow(1+r, n)) / (math.Pow(1+r, n) - 1)

	overpayment = (payment * n) - loanSum

	payment = math.Round(payment*100) / 100.0
	overpayment = math.Round(overpayment*100) / 100.0

	lastDate = time.Now().AddDate(0, req.Months, 0).Format("2006-01-02")

	return
}

// GetAll Cache Items.
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
