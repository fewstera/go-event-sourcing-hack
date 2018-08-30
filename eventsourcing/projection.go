package eventsourcing

import (
	"fmt"
)

type Projection struct {
	users map[string]*User
}

func NewProjection() *Projection {
	projection := new(Projection)
	projection.users = make(map[string]*User)
	return projection
}

func (projection *Projection) Apply(event Event) error {
	fmt.Println(event)
	return nil
}

func (projection *Projection) GetUser(userId string) (*User, error) {
	_, userExistsInRepo := projection.users[userId]
	if !userExistsInRepo {
		return nil, &UserNotFoundError{fmt.Sprintf("User id %v not found", userId)}
	}

	user := projection.users[userId]
	return user, nil

}
