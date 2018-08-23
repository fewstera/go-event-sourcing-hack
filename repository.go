package main

import "fmt"

type Repository struct {
	users map[string]*User
}

func NewRepository() *Repository {
	repository := new(Repository)
	repository.users = make(map[string]*User)
	return repository
}

func (repository *Repository) SaveUser(user *User) {
	userId := user.GetId()
	_, userExistsInRepo := repository.users[userId]
	if !userExistsInRepo {
		repository.users[userId] = user
	}

	user.MarkChangesAsCommitted()
}

func (repository *Repository) GetUser(userId string) (*User, error) {
	_, userExistsInRepo := repository.users[userId]
	if !userExistsInRepo {
		return nil, &UserNotFoundError{fmt.Sprintf("User id %v not found", userId)}
	}

	user := repository.users[userId]
	return user, nil

}
