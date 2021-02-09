package validators

import (
	"errors"
	"microservices/pb"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

var (
	ErrInvalidUserId      = errors.New("invalid userId")
	ErrEmptyName          = errors.New("name can't be empty")
	ErrEmptyEmail         = errors.New("email can't be empty")
	ErrEmptyPassword      = errors.New("password can't be empty")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrSignInFailed       = errors.New("signin failed")
)

func ValidateSignUp(user *pb.User) error {
	if !bson.IsObjectIdHex(user.Id) {
		return ErrInvalidUserId
	}
	if user.Email == "" {
		return ErrEmptyEmail
	}
	if user.Name == "" {
		return ErrEmptyName
	}
	if user.Password == "" {
		return ErrEmptyPassword
	}
	return nil
}

func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
