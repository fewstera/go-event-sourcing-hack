package user

import "fmt"

type InvalidAgeError struct {
	message string
}

func (e *InvalidAgeError) Error() string {
	return e.message
}

type UnkownCommandError struct {
	message string
}

func (e *UnkownCommandError) Error() string {
	return e.message
}

type UserNotFoundError struct {
	message string
}

func (e *UserNotFoundError) Error() string {
	return e.message
}

// EventNumberSyncError is an error for when
type EventNumberSyncError struct {
	Expected int
	Got      int
}

func (e *EventNumberSyncError) Error() string {
	return fmt.Sprintf("event number sync error: expected %d, got %d", e.Expected, e.Got)
}

type InsufficientFundsError struct {
	Balance   float32
	Requested float32
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("not enough funds in account: balance is %v, requested %v", e.Balance, e.Requested)
}
