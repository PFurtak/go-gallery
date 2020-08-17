package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a DB resource is not found
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID
	// is passed into method
	ErrInvalidID = errors.New("models: ID provided was invalid")
)

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

type UserService struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}

// ByID is used to look up user by provided ID
// 1 - user, nil
// 2 - nil, ErrNotFound
// nil, otherError
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail is used to look up a user by provided email
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
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

// Create will create the provided user
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// Update will update the provided user with provided data
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// Delete will remove user from DB with the provided ID
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// DestructiveReset drops user table and rebuilds
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate will attempt to migrate the users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// Close ends user service DB connection
func (us *UserService) Close() error {
	return us.db.Close()
}
