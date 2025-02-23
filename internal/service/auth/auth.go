package auth

type AuthUser interface {
	Login(username, password string) (string, error)
	Validate() (bool, error)
	Refresh() (string, error)
}
