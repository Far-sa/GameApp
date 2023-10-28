package authorizationservice

import "game-app/entity"

type Repository interface {
	GetUserPermissionTitles(userID uint) ([]entity.PermissionTitle, error)
}

type Service struct{}

func (s Service) CheckAccess(userID uint, permissions ...entity.PermissionTitle) (bool, error) {
	// get user role
	// get all permissions for the given role

	// merge ACLs

	// check the access

	return false, nil
}
