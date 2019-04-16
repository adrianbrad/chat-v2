package user

type User struct {
	ID          string
	Nickname    string
	Permissions map[string]struct{}
}
