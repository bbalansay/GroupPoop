package users

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name           string
		expectedOutput string
		input          *NewUser
	}{
		{
			"Invalid email address",
			"This is not a valid email address",
			&NewUser{
				Email:        "realEmail",
				Password:     "123456",
				PasswordConf: "123456",
				UserName:     "myUsername",
				FirstName:    "My",
				LastName:     "name",
			},
		},
		{
			"Password is too short",
			"Password must be at least 6 characters",
			&NewUser{
				Email:        "real@Email.com",
				Password:     "1234",
				PasswordConf: "1234",
				UserName:     "myUsername",
				FirstName:    "My",
				LastName:     "name",
			},
		},
		{
			"Mismatch Passwords",
			"Passwords do not match",
			&NewUser{
				Email:        "real@Email.com",
				Password:     "1234567",
				PasswordConf: "123456",
				UserName:     "myUsername",
				FirstName:    "My",
				LastName:     "name",
			},
		},
		{
			"No username",
			"UserName must be greater than 0 characters",
			&NewUser{
				Email:        "real@Email.com",
				Password:     "123456",
				PasswordConf: "123456",
				UserName:     "",
				FirstName:    "My",
				LastName:     "name",
			},
		},
		{
			"Spaces in username",
			"UserName may not contain spaces",
			&NewUser{
				Email:        "real@Email.com",
				Password:     "123456",
				PasswordConf: "123456",
				UserName:     "my Username",
				FirstName:    "My",
				LastName:     "name",
			},
		},
		{
			"Multiple Errors",
			"This is not a valid email address",
			&NewUser{
				Email:        "real Email.com",
				Password:     "1234",
				PasswordConf: "1234",
				UserName:     "my Username",
				FirstName:    "My",
				LastName:     "name",
			},
		},
		{
			"No Errors",
			"",
			&NewUser{
				Email:        "real@Email.com",
				Password:     "123456",
				PasswordConf: "123456",
				UserName:     "myUsername",
				FirstName:    "My",
				LastName:     "name",
			},
		},
	}

	for _, c := range cases {
		err := c.input.Validate()
		if err != nil {
			if err.Error() != c.expectedOutput {
				t.Errorf("Case: %s\n value: %s did not match expected output:\n %v\n", c.name, err, c.expectedOutput)
			}
		}
	}
}

func TestFullName(t *testing.T) {
	cases := []struct {
		name           string
		expectedOutput string
		input          *User
	}{
		{
			"Has First name and Last name",
			"Alex Wong",
			&User{
				ID:        0,
				Email:     "my@email.com",
				PassHash:  []byte("password"),
				UserName:  "myUsername",
				FirstName: "Alex",
				LastName:  "Wong",
				PhotoURL:  "myphoto.com",
			},
		},
		{
			"Only First name",
			"Alex",
			&User{
				ID:        0,
				Email:     "my@email.com",
				PassHash:  []byte("password"),
				UserName:  "myUsername",
				FirstName: "Alex",
				LastName:  "",
				PhotoURL:  "myphoto.com",
			},
		},
		{
			"Only Last name",
			"Wong",
			&User{
				ID:        0,
				Email:     "my@email.com",
				PassHash:  []byte("password"),
				UserName:  "myUsername",
				FirstName: "",
				LastName:  "Wong",
				PhotoURL:  "myphoto.com",
			},
		},
	}

	for _, c := range cases {
		res := c.input.FullName()

		if res != c.expectedOutput {
			t.Errorf("Case: %s\n value: %s did not match expected output:\n %v\n", c.name, res, c.expectedOutput)
		}
	}
}

func TestSetPassword(t *testing.T) {
	cases := []struct {
		name           string
		expectedOutput string
		current        *User
		password       string
	}{
		{
			"Hash a password",
			"password",
			&User{
				ID:        0,
				Email:     "my@email.com",
				PassHash:  nil,
				UserName:  "myUsername",
				FirstName: "Alex",
				LastName:  "Wong",
				PhotoURL:  "myphoto.com",
			},
			"password",
		},
		{
			"Empty password",
			"error generating bcrypt hash:",
			&User{
				ID:        0,
				Email:     "my@email.com",
				PassHash:  nil,
				UserName:  "myUsername",
				FirstName: "Alex",
				LastName:  "Wong",
				PhotoURL:  "myphoto.com",
			},
			"",
		},
		{
			"Error Hashing",
			"error generating bcrypt hash:",
			&User{
				ID:        0,
				Email:     "my@email.com",
				PassHash:  nil,
				UserName:  "myUsername",
				FirstName: "Alex",
				LastName:  "Wong",
				PhotoURL:  "myphoto.com",
			},
			"This is going to be a really really really long password that I am trying to make in order to make the hash fail",
		},
	}

	for _, c := range cases {
		err := c.current.SetPassword(c.password)
		if err != nil {
			if err.Error() != c.expectedOutput {
				t.Errorf("Case: %s\n value: %s did not match expected output\n %v\n", c.name, err, c.expectedOutput)
			}
		} else {
			res := c.current.PassHash

			if bcrypt.CompareHashAndPassword(res, []byte(c.password)) != nil {
				t.Errorf("Case: %s\n value: %s did not match expected output\n %v\n", c.name, res, c.expectedOutput)
			}
		}
	}
}

func TestAuthenticate(t *testing.T) {
	cases := []struct {
		name           string
		expectedOutput string
		input          *User
		password       string
	}{
		{
			"Authenticate valid password",
			"$2y$12$rvdG8Ta75p4uOm/1LTWQsOJ9U3.RNAiGp2faRL.fLR..A1aze4JRW",
			&User{
				ID:        0,
				Email:     "my@email.com",
				PassHash:  []byte("$2y$12$rvdG8Ta75p4uOm/1LTWQsOJ9U3.RNAiGp2faRL.fLR..A1aze4JRW"),
				UserName:  "myUsername",
				FirstName: "Alex",
				LastName:  "Wong",
				PhotoURL:  "myphoto.com",
			},
			"password",
		},
		{
			"Authenticate wrong password",
			"password doesn't match stored hash",
			&User{
				ID:        0,
				Email:     "my@email.com",
				PassHash:  []byte("$2y$12$rvdG8Ta75p4uOm/1LTWQsOJ9U3.RNAiGp2faRL.fLR..A1aze4JRW"),
				UserName:  "myUsername",
				FirstName: "Alex",
				LastName:  "Wong",
				PhotoURL:  "myphoto.com",
			},
			"diffpassword",
		},
	}

	for _, c := range cases {
		err := c.input.Authenticate(c.password)
		if err != nil {
			if err.Error() != c.expectedOutput {
				t.Errorf("Case: %s\n value: %s did not match expected output\n %v\n", c.name, err, c.expectedOutput)
			}
		} else {
			res := string(c.input.PassHash)

			if res != c.expectedOutput {
				t.Errorf("Case: %s\n value: %s did not match expected output\n %v\n", c.name, res, c.expectedOutput)
			}
		}
	}

}

func TestApplyUpdates(t *testing.T) {
	cases := []struct {
		name           string
		expectedOutput string
		oldInput       *User
		newInput       *Updates
	}{
		{
			"Update first and last name",
			"New Name",
			&User{
				ID:        0,
				Email:     "my@email.com",
				PassHash:  []byte("password"),
				UserName:  "myUsername",
				FirstName: "Alex",
				LastName:  "Wong",
				PhotoURL:  "myphoto.com",
			},
			&Updates{
				FirstName: "New",
				LastName:  "Name",
			},
		},
		{
			"Empty first name update",
			"updates not valid",
			&User{
				ID:        0,
				Email:     "my@email.com",
				PassHash:  []byte("password"),
				UserName:  "myUsername",
				FirstName: "Alex",
				LastName:  "Wong",
				PhotoURL:  "myphoto.com",
			},
			&Updates{
				FirstName: "",
				LastName:  "Last",
			},
		},
		{
			"Empty last name update",
			"updates not valid",
			&User{
				ID:        0,
				Email:     "my@email.com",
				PassHash:  []byte("password"),
				UserName:  "myUsername",
				FirstName: "Alex",
				LastName:  "Wong",
				PhotoURL:  "myphoto.com",
			},
			&Updates{
				FirstName: "First",
				LastName:  "",
			},
		},
	}

	for _, c := range cases {
		err := c.oldInput.ApplyUpdates(c.newInput)
		if err != nil {
			if err.Error() != c.expectedOutput {
				t.Errorf("Case: %s\n value: %s did not match expected output\n %v\n", c.name, err, c.expectedOutput)
			}
		} else {
			res := c.oldInput.FullName()

			if res != c.expectedOutput {
				t.Errorf("Case: %s\n value: %s did not match expected output\n %v\n", c.name, res, c.expectedOutput)
			}
		}
	}
}

func TestToUser(t *testing.T) {
	cases := []struct {
		name          string
		expectedURL   string
		expectedError string
		input         *NewUser
	}{
		{
			"test uppercase in email address",
			"https://www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346",
			"",
			&NewUser{
				Email:        "MyEmailAddress@example.com",
				Password:     "123456",
				PasswordConf: "123456",
				UserName:     "myUsername",
				FirstName:    "First",
				LastName:     "name",
			},
		},
		{
			"test space in email address",
			"https://www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346",
			"",
			&NewUser{
				Email:        "myemailaddress@example.com ",
				Password:     "123456",
				PasswordConf: "123456",
				UserName:     "myUsername",
				FirstName:    "First",
				LastName:     "name",
			},
		},
		{
			"test an invalid email address",
			"",
			"This is not a valid email address",
			&NewUser{
				Email:        "realEmail",
				Password:     "123456",
				PasswordConf: "123456",
				UserName:     "myUsername",
				FirstName:    "My",
				LastName:     "name",
			},
		},
	}

	for _, c := range cases {
		output, err := c.input.ToUser()

		if err != nil {
			if err.Error() != c.expectedError {
				t.Errorf("Case: %s\n value: %s did not match expected output\n %v\n", c.name, err, c.expectedError)
			}
		}

		if output != nil {
			// compare photourl field with precomputed url
			if output.PhotoURL != c.expectedURL {
				t.Errorf("Case: %s\n value: %s did not match expected output\n %v\n", c.name, output.PhotoURL, c.expectedURL)
			}
			// let bcrypt compare hash and password, should return nil if hash is valid
			if bcrypt.CompareHashAndPassword(output.PassHash, []byte(c.input.Password)) != nil {
				t.Errorf("Case: %s\n password hash: %s comparison returned an error\n", c.name, output.PassHash)
			}
			// compare other fields: expect them to be the same
			if output.FullName() != (c.input.FirstName + " " + c.input.LastName) {
				t.Errorf("Case: %s\n full name: %s did not match expected output %s\n", c.name, output.FullName(), c.input.FirstName+" "+c.input.LastName)
			}
			if output.UserName != c.input.UserName {
				t.Errorf("Case: %s\n username: %s did not match expected output %s\n", c.name, output.UserName, c.input.UserName)
			}
		}
	}
}
