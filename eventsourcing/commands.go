package eventsourcing

type Command interface {
	ImplementsCommand()
}

type CreateUserCommand struct {
	Id   string
	Name string
	Age  int
}

func (c CreateUserCommand) ImplementsCommand() {}

type ChangeUsersNameCommand struct {
	Id      string
	NewName string
}

func (c ChangeUsersNameCommand) ImplementsCommand() {}

type IncreaseUsersAgeCommand struct {
	Id string
}

func (c IncreaseUsersAgeCommand) ImplementsCommand() {}
