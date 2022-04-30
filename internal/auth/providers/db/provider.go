package db

import (
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"lproxy/internal/auth/providers"
	"time"
)

const ProviderType = "db"
const defaultLifeTime = 3600

type Provider struct {
	code        string
	isActive    bool
	dsn         string
	lifeTime    int64
	builder     builder
	db          *sqlx.DB
	userRepo    *userRepository
	roleRepo    *roleRepository
	sessionRepo *sessionRepository
}

func NewDbProvider(config *providers.Config) *Provider {
	db := initDb(config.Driver, config.Dsn)
	userRepo := initUserRepo(db)
	roleRepo := initRoleRepo(db)
	sessionRepo := initSessionRepo(db)

	return &Provider{
		code:        config.Code,
		isActive:    config.IsActive,
		lifeTime:    getSessionLifeTime(config.Lifetime),
		dsn:         config.Dsn,
		db:          db,
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		sessionRepo: sessionRepo,
		builder: builder{
			userRepo:    userRepo,
			sessionRepo: sessionRepo,
			roleRepo:    roleRepo,
		},
	}
}

func getSessionLifeTime(time int64) int64 {
	if time <= 0{
		return int64(defaultLifeTime)
	}
	return 0
}

func (p *Provider) GetType() string {
	return ProviderType
}

func (p *Provider) GetCode() string {
	return p.code
}

func (p *Provider) IsActive() bool {
	return p.isActive
}

func (p *Provider) AuthenticateUser(token string) (providers.Session, error) {
	session, sessionErr := p.sessionRepo.findByToken(token)
	if sessionErr != nil {
		return providers.Session{}, sessionErr
	}
	if session.ExpiresIn <= time.Now().Unix() {
		return providers.Session{}, &providers.AuthError{Message: "Token has been expired"}
	}

	user, userErr := p.builder.userEntityByUserId(session.UserId)
	if userErr != nil {
		return providers.Session{}, userErr
	}

	return providers.Session{
		User:      user,
		StartIn:   session.StartIn,
		ExpiresIn: session.ExpiresIn,
		Token:     session.Token,
	}, nil
}

func (p *Provider) AuthorizeUser(username string, passwordHash string) (providers.Session, error) {
	userDbData, userFoundErr := p.userRepo.findByName(username)
	if userFoundErr != nil {
		return providers.Session{}, userFoundErr
	}
	if userDbData.PasswordHash != passwordHash {
		return providers.Session{}, &providers.AuthError{Message: "Invalid password"}
	}
	user, userErr := p.builder.userEntityFromDbEntity(*userDbData)
	if userErr != nil {
		return providers.Session{}, userErr
	}
	sessionDbData := Session{
		UserId:    user.Id,
		StartIn:   time.Now().Unix(),
		ExpiresIn: time.Now().Unix() + p.lifeTime,
		Token:     uuid.New().String(),
	}

	saveErr := p.sessionRepo.save(sessionDbData)
	if saveErr != nil {
		return providers.Session{}, saveErr
	}

	return providers.Session{
		User:      user,
		StartIn:   sessionDbData.StartIn,
		ExpiresIn: sessionDbData.ExpiresIn,
		Token:     sessionDbData.Token,
	}, nil
}

func (p *Provider) Shutdown() {
	err := p.db.Close()
	if err != nil {
		glog.Exitf("Db connection cannot close: %s", err.Error())
	}
}
