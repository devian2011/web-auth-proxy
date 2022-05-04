package db

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	_ "github.com/lib/pq"
	"strings"
	"time"
)

func initDb(driver string, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(8)
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	return db, nil
}

type userRepository struct {
	db *sqlx.DB
}

func initUserRepo(db *sqlx.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) findById(id string) (*User, error) {
	user := &User{}
	err := ur.db.Get(user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) findByName(name string) (*User, error) {
	user := &User{}
	err := ur.db.Get(user, "SELECT * FROM users WHERE username = $1", name)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) create(user *User) (*User, error) {
	foundUser, userFindByNameError := ur.findByName(user.Username)
	if userFindByNameError != nil {
		return nil, userFindByNameError
	}
	if foundUser != nil {
		return nil, errors.New(fmt.Sprintf("User with same username: %s is already exists", user.Username))
	}

	_, err := ur.db.NamedExec("INSERT INTO users (username, password_hash, email, phone) VALUES (:username, :password_hash, :email, :phone)", user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) update(user *User) (*User, error) {
	foundUser, userFindByIdError := ur.findById(user.Id)
	if userFindByIdError != nil {
		return nil, userFindByIdError
	}
	if foundUser == nil {
		return nil, errors.New(fmt.Sprintf("Unknown user with id: %s", user.Id))
	}

	_, err := ur.db.NamedExec("UPDATE users SET password_hash=:password_hash, email = :email, phone=:phone WHERE id = :id", user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) remove(id string) error {
	_, err := ur.db.NamedExec("DELETE FROM users WHERE id=$1", id)
	return err
}

type roleRepository struct {
	db *sqlx.DB
}

func initRoleRepo(db *sqlx.DB) *roleRepository {
	return &roleRepository{
		db: db,
	}
}

func (rr *roleRepository) findAll() (*[]Role, error) {
	userRoles := &[]Role{}
	err := rr.db.Get(userRoles, "SELECT * FROM roles")
	if err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (rr *roleRepository) findForUser(userId string) (*[]Role, error) {
	userRoles := &[]Role{}
	err := rr.db.Select(userRoles, "SELECT roles.* FROM roles JOIN roles_for_user rr on rr.role_id = roles.id WHERE rr.user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (rr *roleRepository) findById(roleId string) (*Role, error) {
	role := &Role{}
	err := rr.db.Get(role, "SELECT * FROM roles WHERE id=$1", roleId)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (rr *roleRepository) create(role *Role) (*Role, error) {
	found, errRoleFound := rr.findById(role.Id)
	if errRoleFound != nil {
		return nil, errRoleFound
	}
	if found != nil {
		return nil, errors.New(fmt.Sprintf("Role with id: %s not found", role.Id))
	}
	_, err := rr.db.NamedExec("INSERT INTO roles (name, code) VALUES (:name, :code)", role)
	if err != nil {
		return nil, err
	}

	return role, nil

}

func (rr *roleRepository) update(role *Role) (*Role, error) {
	found, errRoleFound := rr.findById(role.Id)
	if errRoleFound != nil {
		return nil, errRoleFound
	}
	if found == nil {
		return nil, errors.New(fmt.Sprintf("Unknown role with id: %s not found", role.Id))
	}
	_, err := rr.db.NamedExec("UPDATE roles SET name=:name, code=:code, name=:name WHERE id=:id", role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

type sessionRepository struct {
	db *sqlx.DB
}

func initSessionRepo(db *sqlx.DB) *sessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (sr *sessionRepository) save(session Session) error {
	_, err := sr.db.NamedQuery(
		"INSERT INTO sessions (user_id, start_in, expires_in, token) VALUES (:user_id, :start_in, :expired_in, :token)",
		session)
	return err
}

func (sr *sessionRepository) findByToken(token string) (*Session, error) {
	session := &Session{}
	err := sr.db.Get(session, "SELECT * FROM sessions WHERE token = $1", token)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (sr *sessionRepository) findAll() (*[]Session, error) {
	sessions := &[]Session{}
	err := sr.db.Select(sessions, "SELECT * FROM sessions")
	if err != nil {
		return nil, err
	}
	return sessions, err
}

func (sr *sessionRepository) findAllActive() (*[]Session, error) {
	sessions := &[]Session{}
	now := time.Now().UnixNano()
	err := sr.db.Select(sessions, "SELECT * FROM sessions WHERE expires_in > $1", now)
	if err != nil {
		return nil, err
	}
	return sessions, err
}

func (sr *sessionRepository) findAllInactive() (*[]Session, error) {
	sessions := &[]Session{}
	now := time.Now().UnixNano()
	err := sr.db.Select(sessions, "SELECT * FROM sessions WHERE expires_in <= $1", now)
	if err != nil {
		return nil, err
	}
	return sessions, err
}

func (sr *sessionRepository) removeAllInactive() error {
	now := time.Now().UnixNano()
	_, err := sr.db.NamedExec("DELETE FROM sessions WHERE expires_in < $1", now)
	return err
}
