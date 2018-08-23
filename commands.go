package main

type Command interface {
	ImplementsCommand()
}

type CreateUserCommand struct {
	id   string
	name string
	age  int
}

func (c CreateUserCommand) ImplementsCommand() {}

type ChangeUsersNameCommand struct {
	id      string
	newName string
}

func (c ChangeUsersNameCommand) ImplementsCommand() {}

type AgeUserCommand struct {
	id string
}

func (c AgeUserCommand) ImplementsCommand() {}
