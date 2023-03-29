package auth

import "github.com/gin-gonic/gin"

type User struct {
	Password string
}

type Auth struct {
	users map[string]User
}

func New(users map[string]User) *Auth {
	return &Auth{
		users: users,
	}
}

func (a *Auth) GetHandlerFunc() gin.HandlerFunc {
	if len(a.users) == 0 {
		// empty middleware
		return func(c *gin.Context) { c.Next() }
	}

	accounts := make(gin.Accounts)

	for name, user := range a.users {
		accounts[name] = user.Password
	}

	return gin.BasicAuth(accounts)
}
