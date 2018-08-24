package eventsourcing

type Event interface {
	ImplementsEvent()
}

type UserCreatedEvent struct {
	Id   string
	Name string
	Age  int
}

func (e UserCreatedEvent) ImplementsEvent() {}

type UsersNameChangedEvent struct {
	Id      string
	NewName string
}

func (e UsersNameChangedEvent) ImplementsEvent() {}

type UserGotOlderEvent struct {
	Id string
}

func (e UserGotOlderEvent) ImplementsEvent() {}
