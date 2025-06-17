package service

import (
	"sber_test/repo/cache"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	c := cache.New()
	s := New(c)

	req := ExecuteRequest{
		ObjectCost:     5000000,
		InitialPayment: 2999999,
		Months:         240,
		Program: map[string]bool{
			"military": true,
		},
	}

	resp, id, err := s.Execute(req)

	assert.Nil(t, err, "Expected no error")

	assert.Greater(t, id, 0, "ID should be greater than 0")

	assert.Equal(t, req.InitialPayment, resp.Params.InitialPayment, "InitialPayment should match")
	assert.Equal(t, req.Months, resp.Params.Months, "Months should match")
	assert.Equal(t, 9, resp.Aggregates.Rate, "Rate should be 9 for military program")
	assert.Equal(t, 2000001.0, resp.Aggregates.LoanSum, "LoanSum should match")
}

func TestExecuteWithInvalidProgram(t *testing.T) {
	c := cache.New()
	s := New(c)

	req := ExecuteRequest{
		ObjectCost:     5000000,
		InitialPayment: 2999999,
		Months:         240,
		Program: map[string]bool{
			"invalid_program": true,
		},
	}

	_, _, err := s.Execute(req)

	assert.NotNil(t, err, "Expected error for invalid program")
	assert.Equal(t, "unknown program: invalid_program", err.Error(), "Error message should match")
}

func TestExecuteWithMultiplePrograms(t *testing.T) {
	c := cache.New()
	s := New(c)

	req := ExecuteRequest{
		ObjectCost:     5000000,
		InitialPayment: 2999999,
		Months:         240,
		Program: map[string]bool{
			"salary":   true,
			"military": true,
		},
	}

	_, _, err := s.Execute(req)

	assert.NotNil(t, err, "Expected error for choosing more than one program")
	assert.Equal(t, "choose only 1 program", err.Error(), "Error message should match")
}

func TestExecuteWithLowInitialPayment(t *testing.T) {
	c := cache.New()
	s := New(c)

	req := ExecuteRequest{
		ObjectCost:     5000000,
		InitialPayment: 99999,
		Months:         240,
		Program: map[string]bool{
			"military": true,
		},
	}

	_, _, err := s.Execute(req)

	assert.NotNil(t, err, "Expected error for initial payment too low")
	assert.Equal(t, "the initial payment should be more", err.Error(), "Error message should match")
}

func TestExecuteWithEmptyProgram(t *testing.T) {
	c := cache.New()
	s := New(c)

	req := ExecuteRequest{
		ObjectCost:     5000000,
		InitialPayment: 2999999,
		Months:         240,
		Program:        map[string]bool{},
	}

	_, _, err := s.Execute(req)

	assert.NotNil(t, err, "Expected error for empty program")
	assert.Equal(t, "choose program", err.Error(), "Error message should match")
}

func TestCacheWithEmptyItems(t *testing.T) {
	c := cache.New()
	s := New(c)

	cacheItems := s.GetAll()
	assert.Empty(t, cacheItems, "Cache should be empty initially")
}

func TestExecuteWithZeroObjectCost(t *testing.T) {
	c := cache.New()
	s := New(c)

	req := ExecuteRequest{
		ObjectCost:     0,
		InitialPayment: 0,
		Months:         240,
		Program: map[string]bool{
			"military": true,
		},
	}

	resp, id, err := s.Execute(req)

	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, 0.0, resp.Aggregates.LoanSum, "Loan sum should be 0 when ObjectCost is 0")
	assert.Equal(t, 0.0, resp.Aggregates.MonthlyPayment, "Monthly payment should be 0 when ObjectCost is 0")
	assert.Greater(t, id, 0, "ID should be greater than 0")
}

func TestExecuteWithLargeObjectCost(t *testing.T) {
	c := cache.New()
	s := New(c)

	req := ExecuteRequest{
		ObjectCost:     1e12,
		InitialPayment: 5e11,
		Months:         240,
		Program: map[string]bool{
			"military": true,
		},
	}

	resp, id, err := s.Execute(req)

	assert.Nil(t, err, "Expected no error")
	assert.Greater(t, resp.Aggregates.LoanSum, 0.0, "Loan sum should be positive")
	assert.Greater(t, resp.Aggregates.MonthlyPayment, 0.0, "Monthly payment should be positive")
	assert.Greater(t, id, 0, "ID should be greater than 0")
}

func TestExecuteWithNoProgram(t *testing.T) {
	c := cache.New()
	s := New(c)

	req := ExecuteRequest{
		ObjectCost:     5000000,
		InitialPayment: 2999999,
		Months:         240,
		Program:        map[string]bool{},
	}

	_, _, err := s.Execute(req)

	assert.NotNil(t, err, "Expected error for no program selected")
	assert.Equal(t, "choose program", err.Error(), "Error message should match")
}

func TestExecuteWithLongTerm(t *testing.T) {
	c := cache.New()
	s := New(c)

	req := ExecuteRequest{
		ObjectCost:     5000000,
		InitialPayment: 2999999,
		Months:         360,
		Program: map[string]bool{
			"military": true,
		},
	}

	resp, id, err := s.Execute(req)

	assert.Nil(t, err, "Expected no error")
	assert.Greater(t, id, 0, "ID should be greater than 0")
	assert.Equal(t, 9, resp.Aggregates.Rate, "Rate should be 9 for military program")
	assert.Equal(t, 2000001.0, resp.Aggregates.LoanSum, "LoanSum should match")
}

func TestCacheWithAddedItems(t *testing.T) {
	c := cache.New()
	s := New(c)

	req := ExecuteRequest{
		ObjectCost:     5000000,
		InitialPayment: 2999999,
		Months:         240,
		Program: map[string]bool{
			"military": true,
		},
	}

	_, id, err := s.Execute(req)
	assert.Nil(t, err)
	assert.Greater(t, id, 0, "ID should be greater than 0")

	cacheItems := s.GetAll()
	assert.NotEmpty(t, cacheItems, "Cache should not be empty")
	assert.Equal(t, id, cacheItems[0].ID, "ID should match the inserted item")
}
