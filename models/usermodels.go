package models

import (
	"errors"
	"regexp"
	"strings"

	"github.com/Users/patrickfurtak/desktop/go-gallery/hash"
	"github.com/Users/patrickfurtak/desktop/go-gallery/rand"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound is returned when a DB resource is not found
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID
	// is passed into method
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned on failed password and hash match
	ErrInvalidPassword = errors.New("models: Password invalid")

	// ErrInvalidPasswordLength is returned on failed password length check
	ErrInvalidPasswordLength = errors.New("models: Password must be atleast 8 chars long")

	// ErrPasswordRequired is returned when the supplied password is empty
	ErrPasswordRequired = errors.New("models: Password field is required")

	// ErrPasswordNotHashed is returned when a password is not hashed
	ErrPasswordNotHashed = errors.New("models: Password is not hashed")

	// ErrRememberTooShort is returned when a remember token has fewer than 32 bytes
	ErrRememberTooShort = errors.New("models: Remember token has less than 32 bytes, too short")

	// ErrRememberNotHashed is returned when a remember token is not hashed
	ErrRememberNotHashed = errors.New("models: Remember token is not hashed")
)

const userPwPepper = "sadjfhusdfjhsdfbchfdsssswqdnfgchdnsdfhdskjdbfuv"
const hmacKey = "notsosecret"

// UserDB is used to interact with the users database
type UserDB interface {
	// Methods for querying single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Close is used to close the DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacKey)
	uv := newUserValidator(ug, hmac)
	return &userService{
		UserDB: uv,
	}, nil
}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &userGorm{
		db: db,
	}, nil
}

var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

// Update will update the provided user with provided data
func (uv *userValidator) Update(user *User) error {
	if err := runUserValFuncs(
		user,
		uv.pwLengthCheck,
		uv.bcryptPassword,
		uv.pwHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailExistCheck); err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// Create will create the provided user
func (uv *userValidator) Create(user *User) error {
	if err := runUserValFuncs(
		user,
		uv.pwRequired,
		uv.pwLengthCheck,
		uv.bcryptPassword,
		uv.pwHashRequired,
		uv.setRememberIfUnset,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailExistCheck); err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

func (uv *userValidator) emailExistCheck(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		return nil
	}

	if err != nil {
		return err
	}
	if user.ID != existing.ID {
		return errors.New("models: this email has already been used to sign up")
	}
	return nil
}

type userValFunc func(*User) error

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

// bcryptPassword will hash a users password with a predefined pepper
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(pwBytes), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) pwLengthCheck(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrInvalidPasswordLength
	}
	return nil
}

func (uv *userValidator) pwRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) pwHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordNotHashed
	}
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrRememberNotHashed
	}
	return nil
}

func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

func (uv *userValidator) idValidate(user *User) error {
	if user.ID <= 0 {
		return ErrInvalidID
	}
	return nil
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if !uv.emailRegex.MatchString(user.Email) {
		return errors.New("Email address is not valid")
	}
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return errors.New("Email address is required")
	}
	return nil
}

// ByEmail will normalize the email address before writing and reading to the DB
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	err := runUserValFuncs(&user, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// Delete will remove user from DB with the provided ID
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFuncs(&user, uv.idValidate)
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return nil
	}
	if n < 32 {
		return ErrRememberTooShort
	}
	return nil
}

type userService struct {
	UserDB
}

// UserService is a set of methods used to manipulate and work with user model
type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

// ByRemember looks up by user remember token. Returns the user.
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// User is the struct that dictates the database fields
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// ByID is used to look up user by provided ID
// 1 - user, nil
// 2 - nil, ErrNotFound
// nil, otherError
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail is used to look up a user by provided email
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// Authenticate is used to authenticate with a provided email and password
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}
	return foundUser, nil
}

// first will query using the provided gorm.db and will get the first item
// returned and place it into dst. Helper function to keep code DRY.
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// Update will update the provided user with provided data
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete will remove user from DB with the provided ID
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// DestructiveReset drops user table and rebuilds
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// AutoMigrate will attempt to migrate the users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// Close ends user service DB connection
func (ug *userGorm) Close() error {
	return ug.db.Close()
}
