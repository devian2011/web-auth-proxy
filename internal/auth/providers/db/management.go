package db

import "github.com/jmoiron/sqlx"

type controlManager struct {
	db          *sqlx.DB
	userRepo    *userRepository
	roleRepo    *roleRepository
	sessionRepo *sessionRepository
}

func NewControlManager(db *sqlx.DB, userRepo *userRepository, roleRepo *roleRepository, sessionRepo *sessionRepository) *controlManager {
	return &controlManager{
		db:          db,
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		sessionRepo: sessionRepo,
	}
}




