package db

import (
	"github.com/jmoiron/sqlx"
	"lproxy/internal/auth/providers"
)

type dataManager struct {
	db          *sqlx.DB
	userRepo    *userRepository
	roleRepo    *roleRepository
	sessionRepo *sessionRepository
}

func InitData(db *sqlx.DB, userRepo *userRepository, roleRepo *roleRepository, sessionRepo *sessionRepository) *dataManager {
	return &dataManager{
		db:          db,
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		sessionRepo: sessionRepo,
	}
}

func (dm *dataManager) findSessionByToken(token string) (*providers.Session, error) {
	sessData, sessFindErr := dm.sessionRepo.findByToken(token)
	if sessFindErr != nil {
		return nil, sessFindErr
	}
	userData, userFindErr := dm.userRepo.findById(sessData.UserId)
	if userFindErr != nil {
		return nil, userFindErr
	}
	rolesData, rolesFindErr := dm.roleRepo.findForUser(userData.Id)
	if rolesFindErr != nil {
		return nil, rolesFindErr
	}

	var roles []providers.Role
	for _, rd := range *rolesData {
		roles = append(roles, providers.Role{
			Id:     rd.Id,
			Name:   rd.Name,
			Code:   rd.Code,
			Source: ProviderType,
		})
	}

	return &providers.Session{
		User: providers.User{
			Id:       userData.Id,
			Username: userData.Username,
			Email:    userData.Email,
			Phone:    userData.Phone,
			Provider: ProviderType,
			Source:   "proxy",
			Roles:    roles,
		},
		StartIn:   sessData.StartIn,
		ExpiresIn: sessData.ExpiresIn,
		Token:     sessData.Token,
	}, nil

}
