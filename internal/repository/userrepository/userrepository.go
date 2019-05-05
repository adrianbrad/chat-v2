package userrepository

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/adrianbrad/chat-v2/internal/repository"
	"github.com/adrianbrad/chat-v2/internal/user"
)

type UserRepositoryDB struct {
	repository.DBRepository
	users map[string]*user.User

	addUser    chan *user.User
	updateUser chan *user.User

	repoStopped chan struct{}
	stopRepo    chan struct{}
}

func NewUserRepositoryDB(db *sql.DB) *UserRepositoryDB {
	r := &UserRepositoryDB{
		DBRepository: db,
		users:        make(map[string]*user.User),
	}

	r.addUser = make(chan *user.User)
	r.updateUser = make(chan *user.User)
	r.repoStopped = make(chan struct{}, 1)
	r.stopRepo = make(chan struct{}, 1)

	go r.run()

	return r
}

//if user exists in the map then make an update, otherwise add it to the map
func (r *UserRepositoryDB) run() {
	for {
		select {
		case u := <-r.addUser:
			r.users[u.ID] = u
			continue

		default:
			select {
			case userUpdates := <-r.updateUser:
				user, userExists := r.users[userUpdates.ID]
				if userExists {
					user.Update(*userUpdates)
				}
				continue

			default:
				select {
				case <-r.stopRepo:
					r.repoStopped <- struct{}{}
					return
				default:
				}
			}
		}
	}
}

func (r *UserRepositoryDB) stop() {
	r.stopRepo <- struct{}{}
	<-r.repoStopped
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

	r.addUser <- u
	return
}

func (r *UserRepositoryDB) Create(user user.User) (err error) {
	tx, err := r.Begin()
	if err != nil {
		return
	}

	res, err := tx.Exec(`
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

	stmt, err := tx.Prepare(pq.CopyIn("users_permissions", "user_id", "permission_id"))
	if err != nil {
		return
	}

	for permissionID := range user.Permissions {
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
	if err != nil {
		return
	}

	log.Infof("Created User: %+v", user)
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

	for permissionID := range user.Permissions {
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

	r.updateUser <- &user

	return
}

func (r *UserRepositoryDB) UpdatePermissions(userID string, permissions []string) (err error) {
	tx, err := r.Begin()
	if err != nil {
		return
	}

	_, err = tx.Exec(`
	DELETE FROM users_permissions
	WHERE user_id=$1
	`, userID)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(pq.CopyIn("users_permissions", "user_id", "permission_id"))
	if err != nil {
		return
	}

	permissionsMap := make(map[string]struct{})

	for _, permissionID := range permissions {
		_, err = stmt.Exec(userID, permissionID)
		if err != nil {
			return
		}
		permissionsMap[permissionID] = struct{}{}
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
	if err != nil {
		return
	}
	u := user.User{
		ID:          userID,
		Permissions: permissionsMap,
	}

	r.updateUser <- &u
	return
}

func (r *UserRepositoryDB) UpdateNickname(userID string, updatedNickname string) (err error) {
	res, err := r.Exec(`
	UPDATE users
		SET nickname=$2, updated_at=now() 
	WHERE user_id=$1
	`, userID, updatedNickname)
	if err != nil {
		return err
	}

	if c, _ := res.RowsAffected(); c != 1 {
		err = fmt.Errorf("Invalid number of rows affected by create: %d", c)
		return
	}

	u := user.User{
		ID:       userID,
		Nickname: updatedNickname,
	}

	r.updateUser <- &u
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
