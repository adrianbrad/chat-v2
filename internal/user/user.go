package user

type User struct {
	ID          string
	Nickname    string
	Permissions map[string]struct{}
}

func (u *User) Update(userUpdates User) {
	if userUpdates.Nickname != "" {
		u.Nickname = userUpdates.Nickname
	}
	if userUpdates.Permissions != nil && len(userUpdates.Permissions) != 0 {
		u.Permissions = userUpdates.Permissions
	}
}
