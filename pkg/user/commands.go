package user

// A Command is a command that can be handled by a CommandHandler
type Command interface {
	ImplementsCommand()
}

// BaseCommand can be embedded by implementers of the Command interface to gain
// implementations for the ImplementsCommand method.
type BaseCommand struct {
	StreamID          string
	ClientEventNumber int
}

// ImplementsCommand is a noop to make this an implementation of Command
func (c BaseCommand) ImplementsCommand() {}

// A CreateUserCommand is used to create a new intel
type CreateUserCommand struct {
	BaseCommand
	Name string
	Age  int
}

// NewCreateUserCommand creates a new create intel command
func NewCreateUserCommand(streamID string, name string, age int) CreateUserCommand {
	c := CreateUserCommand{}
	c.BaseCommand = BaseCommand{streamID, 1}
	c.Name = name
	c.Age = age
	return c
}
