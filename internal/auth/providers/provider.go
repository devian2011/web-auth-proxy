package providers

type ProviderError struct {
	Message string
}

func (pe *ProviderError) Error() string {
	return pe.Message
}

type AuthError struct {
	Message string
}

func (ae *AuthError) Error() string {
	return ae.Message
}

type User struct {
	Id       string `json:"id"` // Provider user ID
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Provider string `json:"provider"` // Provider code (db, oauth2 or other providers)
	Source   string `json:"source"`   // Source which contains this user
	Roles    []Role `json:"roles"`
}

type Role struct {
	Id     string `json:"id"`     // Role ID in source
	Name   string `json:"name"`   // Role name (User, Admin, Moderator, Article's Moderator and etc.)
	Code   string `json:"code"`   // Role code (USER, ADMIN, MODERATOR, USER_ADMIN, ARTICLE_MODERATOR)
	Source string `json:"source"` // Source which contains this role
}

type Session struct {
	User      User   `json:"user"`       // User info
	StartIn   int64  `json:"start_in"`   // Session start time
	ExpiresIn int64  `json:"expires_in"` // Session expires time
	Token     string `json:"token"`      // Session token
}

type Provider interface {
	GetType() string // Type of provider (db, active directory, oauth2 and other)
	GetCode() string // Provider code, for example we can use to db providers, client_db_provider, owner_db_provider
	IsActive() bool
	AuthenticateUser(token string) (Session, error)
	AuthorizeUser(username string, passwordHash string) (Session, error)
	Shutdown()
}

type Config struct {
	Type     string            `json:"type"`
	Code     string            `json:"code"`
	IsActive bool              `json:"is_active"`
	Driver   string            `json:"driver"`   // Database driver name - postgres, mysql, etc.
	Dsn      string            `json:"dsn"`      // Authorization service, db or another dsn
	Lifetime int64             `json:"lifetime"` // Auth provider lifetime in minutes
	Options  map[string]string `json:"options"`
}
