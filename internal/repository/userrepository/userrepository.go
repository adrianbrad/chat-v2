package userrepository

import (
	"database/sql"

	"github.com/lib/pq"

	"github.com/adrianbrad/chat-v2/internal/repository"
	"github.com/adrianbrad/chat-v2/internal/user"
)

type UserRepositoryDB struct {
	repository.DBRepository
}

func NewUserRepositoryDB(db *sql.DB) *UserRepositoryDB {
	return &UserRepositoryDB{db}
}

func (r *UserRepositoryDB) GetOne(id string) (u *user.User, err error) {
	var userPermissions pq.StringArray
	u = &user.User{}
	err = r.QueryRow(`
	SELECT user_id, nickname,
		(SELECT array(SELECT permission_id FROM users_permissions WHERE user_id = $1)) AS "Permissions" 
	FROM users
	WHERE user_id=$1 
	`, id).
		Scan(
			&u.ID,
			&u.Nickname,
			&userPermissions,
		)

	if err != nil {
		return
	}
	u.Permissions = make(map[string]struct{})

	for _, permission := range userPermissions {
		u.Permissions[permission] = struct{}{}
	}
	return
}
