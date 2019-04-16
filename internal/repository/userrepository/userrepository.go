package userrepository

import (
	"database/sql"
	"fmt"

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

func (r *UserRepositoryDB) Create(user user.User) (err error) {
	res, err := r.Exec(`
	INSERT INTO users(user_id, nickname)
	VALUES($1, $2)
	`, user.ID, user.Nickname)
	if err != nil {
		return err
	}
	if c, _ := res.RowsAffected(); c != 1 {
		err = fmt.Errorf("Invalid number of rows affected by create: %d", c)
		return
	}

	return
}

func (r *UserRepositoryDB) Update(user user.User) (err error) {
	tx, err := r.Begin()
	if err != nil {
		return
	}

	res, err := tx.Exec(`
	UPDATE users
		SET nickname=$2, updated_at=now() 
	WHERE user_id=$1
	`, user.ID, user.Nickname)
	if err != nil {
		return err
	}
	if c, _ := res.RowsAffected(); c != 1 {
		err = fmt.Errorf("Invalid number of rows affected by create: %d", c)
		return
	}

	_, err = tx.Exec(`
	DELETE FROM users_permissions
	WHERE user_id=$1
	`, user.ID)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(pq.CopyIn("users_permissions", "user_id", "permission_id"))
	if err != nil {
		return
	}

	for permissionID, _ := range user.Permissions {
		_, err = stmt.Exec(user.ID, permissionID)
		if err != nil {
			return
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return
	}

	err = stmt.Close()
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (r *UserRepositoryDB) Delete(userID string) (err error) {
	res, err := r.Exec(`
	DELETE FROM users
	WHERE user_id=$1
	`, userID)
	if err != nil {
		return err
	}
	if c, _ := res.RowsAffected(); c != 1 {
		err = fmt.Errorf("Invalid number of rows affected by delete: %d", c)
		return
	}

	return
}
