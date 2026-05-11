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
	if _, ok := au.storage[user]; !ok {
		return false
	}
	if au.globalPassword != "" {
		return pass == au.globalPassword
	}
	return user != "" && pass == user
}

func (au *inMemoryAuthenticator) Users() []string { return au.usernames }

func NewAuthenticator(users []AuthUser, globalPassword string) Authenticator {
	if len(users) == 0 {
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
