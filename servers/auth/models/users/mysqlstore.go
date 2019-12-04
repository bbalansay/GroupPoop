package users

import (
	"database/sql"
	"time"
)

//MySQLStore represents a store for Users
type MySQLStore struct {
	db *sql.DB
}

//NewMySQLStore does stuff
func NewMySQLStore(connection *sql.DB) *MySQLStore {
	return &MySQLStore{
		db: connection,
	}
}

//GetByID returns the User with the given ID
func (ms *MySQLStore) GetByID(id int64) (*User, error) {
	rows, err := ms.db.Query("select id, email, pass_hash, user_name, "+
		"first_name, last_name, photo_url from users where id = ?", id)

	if err != nil {
		return nil, ErrUserNotFound
	}

	defer rows.Close()

	user := User{}

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
			&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, ErrUserNotFound
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

//GetByEmail returns the User with the given email
func (ms *MySQLStore) GetByEmail(email string) (*User, error) {
	rows, err := ms.db.Query("select id, email, pass_hash, user_name, "+
		"first_name, last_name, photo_url from users where email = ?", email)

	if err != nil {
		return nil, ErrUserNotFound
	}

	defer rows.Close()

	user := User{}

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
			&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, ErrUserNotFound
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

//GetByUserName returns the User with the given Username
func (ms *MySQLStore) GetByUserName(username string) (*User, error) {
	rows, err := ms.db.Query("select id, email, pass_hash, user_name, "+
		"first_name, last_name, photo_url from users where user_name = ?", username)

	if err != nil {
		return nil, ErrUserNotFound
	}

	defer rows.Close()

	user := User{}

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
			&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, ErrUserNotFound
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (ms *MySQLStore) Insert(user *User) (*User, error) {
	insq := "insert into users(email, pass_hash, user_name, first_name, last_name, photo_url) values (?,?,?,?,?,?)"
	res, err := ms.db.Exec(insq, user.Email, user.PassHash, user.UserName, user.FirstName, user.LastName, user.PhotoURL)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (ms *MySQLStore) Update(id int64, updates *Updates) (*User, error) {
	insq := "update users set first_name = ?, last_name = ? where id = ?"
	_, err := ms.db.Exec(insq, updates.FirstName, updates.LastName, id)

	if err != nil {
		return nil, ErrUserNotFound
	}

	user, err := ms.GetByID(id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

//Delete deletes the user with the given ID
func (ms *MySQLStore) Delete(id int64) error {
	insq := "delete from users where id = ?"
	_, err := ms.db.Exec(insq, id)

	if err != nil {
		return ErrUserNotFound
	}

	return nil
}

// Log tracks a user sign in
func (ms *MySQLStore) Log(time time.Time, ipaddr string) error {
	insq := "insert into logs(time, ipaddr) values (?,?)"
	_, err := ms.db.Exec(insq, time, ipaddr)

	if err != nil {
		return err
	}

	return nil
}