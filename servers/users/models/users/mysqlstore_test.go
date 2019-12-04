package users

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ms := NewMySQLStore(db)

	user := User{
		ID:        1,
		Email:     "my@email.com",
		PassHash:  []byte("password"),
		UserName:  "myUsername",
		FirstName: "First",
		LastName:  "Last",
		PhotoURL:  "myphoto.com",
	}

	mock.ExpectExec("insert into users").
		WithArgs(&user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = ms.Insert(&user)

	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ms := NewMySQLStore(db)

	user := User{
		ID:        1,
		Email:     "my@email.com",
		PassHash:  []byte("password"),
		UserName:  "myUsername",
		FirstName: "First",
		LastName:  "Last",
		PhotoURL:  "myphoto.com",
	}

	columns := []string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"}

	mock.ExpectQuery("select (.+) from users").
		WithArgs(&user.ID).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL))

	_, err = ms.GetByID(1)

	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("select (.+) from users").
		WithArgs(99).
		WillReturnError(ErrUserNotFound)

	_, err = ms.GetByID(99)

	if err == nil {
		t.Errorf("error was expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ms := NewMySQLStore(db)

	user := User{
		ID:        1,
		Email:     "my@email.com",
		PassHash:  []byte("password"),
		UserName:  "myUsername",
		FirstName: "First",
		LastName:  "Last",
		PhotoURL:  "myphoto.com",
	}

	columns := []string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"}

	mock.ExpectQuery("select (.+) from users").
		WithArgs("my@email.com").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL))

	_, err = ms.GetByEmail("my@email.com")

	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("select (.+) from users").
		WithArgs("notcorrect@email.com").
		WillReturnError(ErrUserNotFound)

	_, err = ms.GetByEmail("notcorrect@email.com")

	if err == nil {
		t.Errorf("error was expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByUserName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ms := NewMySQLStore(db)

	user := User{
		ID:        1,
		Email:     "my@email.com",
		PassHash:  []byte("password"),
		UserName:  "myUsername",
		FirstName: "First",
		LastName:  "Last",
		PhotoURL:  "myphoto.com",
	}

	columns := []string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"}

	mock.ExpectQuery("select (.+) from users").
		WithArgs("myUsername").
		WillReturnRows(sqlmock.NewRows(columns).AddRow(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL))

	_, err = ms.GetByUserName("myUsername")

	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectQuery("select (.+) from users").
		WithArgs("notMyUsername").
		WillReturnError(ErrUserNotFound)

	_, err = ms.GetByUserName("notMyUsername")

	if err == nil {
		t.Errorf("error was expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ms := NewMySQLStore(db)

	updates := Updates{
		FirstName: "NewFirstName",
		LastName:  "NewLastName",
	}

	mock.ExpectExec("update users").
		WithArgs(updates.FirstName, updates.LastName, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = ms.Update(1, &updates)

	if err == nil {
		t.Errorf("error was expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectExec("update users").
		WithArgs(updates.FirstName, updates.LastName, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	user := User{
		ID:        1,
		Email:     "my@email.com",
		PassHash:  []byte("password"),
		UserName:  "myUsername",
		FirstName: "First",
		LastName:  "Last",
		PhotoURL:  "myphoto.com",
	}

	columns := []string{"id", "email", "pass_hash", "user_name", "first_name", "last_name", "photo_url"}

	mock.ExpectQuery("select (.+) from users").
		WithArgs(&user.ID).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(&user.ID, &user.Email, &user.PassHash, &user.UserName, &user.FirstName, &user.LastName, &user.PhotoURL))

	_, err = ms.Update(1, &updates)

	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	mock.ExpectExec("update users").
		WithArgs(updates.FirstName, updates.LastName, 99).
		WillReturnError(ErrUserNotFound)

	_, err = ms.Update(99, &updates)

	if err == nil {
		t.Errorf("error was expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ms := NewMySQLStore(db)

	// expect to delete user id = 1
	mock.ExpectExec("delete from users").
		WithArgs(1).
		WillDelayFor(time.Second).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = ms.Delete(1)
	if err != nil {
		t.Errorf("unexpected error when deleting: %s", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("unfulfilled expectations on Delete(0): %s", err)
	}

	mock.ExpectExec("delete from users").
		WithArgs(99).
		WillReturnError(ErrUserNotFound)

	err = ms.Delete(99)
	if err == nil {
		t.Errorf("expected error when deleting: %s", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("unfulfilled expectations on Delete(0): %s", err)
	}
}
