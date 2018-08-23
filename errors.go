package main

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
