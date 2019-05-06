package user

type User struct {
	ID          string              `json:"id"`
	Nickname    string              `json:"nickname"`
	Permissions map[string]struct{} `json:",omitempty"`
}

func (u *User) Update(userUpdates User) {
	if userUpdates.Nickname != "" {
		u.Nickname = userUpdates.Nickname
	}
	if userUpdates.Permissions != nil && len(userUpdates.Permissions) != 0 {
		u.Permissions = userUpdates.Permissions
	}
}
