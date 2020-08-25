package models

import "strings"

const (
	// ErrNotFound is returned when a DB resource is not found
	ErrNotFound modelError = "models: resource not found"

	// ErrInvalidEmailFormat is returned when the regex validation check for email format fails
	ErrInvalidEmailFormat modelError = "models: invalid email format"

	// ErrEmailRequired is returned when an email is not provided
	ErrEmailRequired modelError = "models: email address is required"

	// ErrEmailAlreadyExists is returned when an email address is already tied to an account
	ErrEmailAlreadyExists modelError = "models: email address provided already exists"

	// ErrInvalidID is returned when an invalid ID
	// is passed into method
	ErrInvalidID privateError = "models: ID provided was invalid"

	// ErrInvalidPassword is returned on failed password and hash match
	ErrInvalidPassword modelError = "models: Password invalid"

	// ErrInvalidPasswordLength is returned on failed password length check
	ErrInvalidPasswordLength modelError = "models: Password must be atleast 8 chars long"

	// ErrPasswordRequired is returned when the supplied password is empty
	ErrPasswordRequired modelError = "models: Password field is required"

	// ErrPasswordNotHashed is returned when a password is not hashed
	ErrPasswordNotHashed privateError = "models: Password is not hashed"

	// ErrRememberTooShort is returned when a remember token has fewer than 32 bytes
	ErrRememberTooShort privateError = "models: Remember token has less than 32 bytes, too short"

	// ErrRememberNotHashed is returned when a remember token is not hashed
	ErrRememberNotHashed privateError = "models: Remember token is not hashed"

	// ErrTitleRequired is returned when a user does not provide a gallery title
	ErrTitleRequired modelError = "models: gallery title is required"

	// ErrUserIDRequired is returned when a user tries to create a gallery and is not signed in
	ErrUserIDRequired privateError = "models: user ID is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
