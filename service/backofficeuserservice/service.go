package backofficeuserservice

import "game-app/entity"

type Service struct{}

func New() Service {
	return Service{}
}

func (s Service) ListAllUsers() ([]entity.User, error) {

	// TODO -> implement it
	list := make([]entity.User, 0)

	list = append(list, entity.User{
		ID:          0,
		PhoneNumber: "fake",
		Name:        "fake",
		Password:    "fake",
		Role:        entity.UserRole,
	})

	return list, nil
}
