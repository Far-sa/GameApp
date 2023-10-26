package userservice

import (
	"game-app/param"
	"game-app/pkg/richerror"
)

func (s Service) Profile(req param.ProfileRequest) (param.ProfileResponse, error) {
	const op = "userservice.Profile"

	user, err := s.repo.GetUserById(uint64(req.UserID))
	if err != nil {
		//!  - use rich error to develope error handling from different layers
		return param.ProfileResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})

		//return ProfileResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	return param.ProfileResponse{Name: user.Name}, nil
}
