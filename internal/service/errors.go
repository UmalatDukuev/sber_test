package service

import (
	"errors"
	"fmt"
)

// Return errors.
var (
	ErrUnknownProgram          = fmt.Errorf("unknown program")
	ErrChooseProgram           = fmt.Errorf("choose program")
	ErrChooseOnlyOneProgram    = fmt.Errorf("choose only 1 program")
	ErrInitialPaymentLow       = fmt.Errorf("the initial payment should be more")
	ErrFirstPaymentExceedsLoan = errors.New("first payment exceeds loan sum")
)
