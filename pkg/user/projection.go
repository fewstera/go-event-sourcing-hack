package user

import (
	"fmt"

	"github.com/fewstera/go-event-sourcing-hack/pkg/eventstore"
)

type Projection struct {
	Users     map[string]*User
	eventChan chan eventstore.Event
}

func NewProjection() *Projection {
	p := new(Projection)
	p.Users = make(map[string]*User)
	p.eventChan = make(chan eventstore.Event, 100)

	go p.receiveEvent()

	return p
}

func (p *Projection) EventChan() chan<- eventstore.Event {
	return p.eventChan
}

func (p *Projection) receiveEvent() {
	// Wait for next event, this line will block until an event is sent
	event := <-p.eventChan
	p.Apply(event)

	// Get next event
	go p.receiveEvent()
}

func (p *Projection) Apply(event eventstore.Event) {
	userID := event.GetStreamID()
	user, err := p.GetUser(userID)
	if err != nil {
		switch err.(type) {
		case *UserNotFoundError:
			// If user doesn't already exist - create it
			p.Users[userID] = new(User)
			user = p.Users[userID]
		default:
			fmt.Println(err)
		}
	}

	user.Apply(event)
}

func (p *Projection) GetUser(userID string) (*User, error) {
	_, userExistsInRepo := p.Users[userID]
	if !userExistsInRepo {
		return nil, &UserNotFoundError{fmt.Sprintf("User id %v not found", userID)}
	}

	user := p.Users[userID]
	return user, nil

}

func (p *Projection) GetAllUsers() []*User {
	var users []*User
	for _, user := range p.Users {
		users = append(users, user)
	}

	return users
}
