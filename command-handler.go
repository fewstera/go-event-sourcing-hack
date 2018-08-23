package main

type CommandHandler struct {
	repository *Repository
}

func NewCommandHandler(repository *Repository) *CommandHandler {
	ch := new(CommandHandler)
	ch.repository = repository
	return ch
}

func (commandHandler *CommandHandler) handle(command Command) error {
	switch c := command.(type) {
	case CreateUserCommand:
		return commandHandler.handleCreateUserCommand(c)
	default:
		return &UnkownCommandError{}
	}
}

func (commandHandler *CommandHandler) handleCreateUserCommand(c CreateUserCommand) error {
	user, err := NewUser(c.id, c.name, c.age)
	if err == nil {
		commandHandler.repository.SaveUser(user)
	}
	return err
}
