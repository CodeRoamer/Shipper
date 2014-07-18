package models

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrUserAlreadyExist = errors.New("User already exist")
	ErrUserNotExist     = errors.New("User does not exist")
	ErrUserNameIllegal  = errors.New("User name contains illegal characters")
)

// User represents the object of individual and member of organization.
type User struct {
	Id              int64
	LowerName       string `xorm:"unique not null"`
	Name            string `xorm:"unique not null"`
	FullName        string
	Email           string `xorm:"unique not null"`
	Password        string `xorm:"not null"`
	LoginName       string
	Avatar          string `xorm:"varchar(2048) not null"`
	Location        string
	Website         string
	IsActive        bool
	IsAdmin         bool
	Created         time.Time `xorm:"created"`
	Updated         time.Time `xorm:"updated"`
}

// HomeLink returns the user home page link.
func (u *User) HomeLink() string {
	return "/user/" + u.Name
}

// IsUserExist checks if given user name exist,
// the user name should be noncased unique.
func IsUserExist(name string) (bool, error) {
	if len(name) == 0 {
		return false, nil
	}
	return x.Get(&User{LowerName: strings.ToLower(name)})
}

// CreateUser creates record of a new user.
func CreateUser(u *User) (*User, error) {
	isExist, err := IsUserExist(u.Name)
	if err != nil {
		return nil, err
	} else if isExist {
		return nil, ErrUserAlreadyExist
	}

	u.LowerName = strings.ToLower(u.Name)

	_, err = x.Insert(u)

	return u, err
}

// GetUserByName returns the user object by given name if exists.
func GetUserByName(name string) (*User, error) {
	if len(name) == 0 {
		return nil, ErrUserNotExist
	}
	user := &User{LowerName: strings.ToLower(name)}
	has, err := x.Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist
	}
	return user, nil
}
