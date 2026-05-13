package auth

type Authenticator interface {
	Verify(user string, pass string) bool
	Users() []string
}

type AuthStore interface {
	Authenticator() Authenticator
	SetAuthenticator(Authenticator)
}

type AuthUser struct {
	User string
	Pass string
}

type inMemoryAuthenticator struct {
	storage        map[string]string
	usernames      []string
	globalPassword string
}

func (au *inMemoryAuthenticator) Verify(user string, pass string) bool {
	if user == "" {
		return false
	}
	// Per-user password (only when the user was declared with a non-empty pass).
	if expected, ok := au.storage[user]; ok && expected != "" && pass == expected {
		return true
	}
	// Global password lets any non-empty user authenticate, even if not
	// pre-declared in `authentication:`. This is the primary path; IN-USER
	// rules then key off the username from the request, not pre-registration.
	if au.globalPassword != "" && pass == au.globalPassword {
		return true
	}
	return false
}

func (au *inMemoryAuthenticator) Users() []string { return au.usernames }

func NewAuthenticator(users []AuthUser, globalPassword string) Authenticator {
	if len(users) == 0 && globalPassword == "" {
		return nil
	}
	au := &inMemoryAuthenticator{
		storage:        make(map[string]string),
		usernames:      make([]string, 0, len(users)),
		globalPassword: globalPassword,
	}
	for _, user := range users {
		au.storage[user.User] = user.Pass
		au.usernames = append(au.usernames, user.User)
	}
	return au
}

var AlwaysValid Authenticator = alwaysValid{}

type alwaysValid struct{}

func (alwaysValid) Verify(string, string) bool { return true }

func (alwaysValid) Users() []string { return nil }
