package use_case

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type LoginUc struct {
	userRepository UserRepository
}

func NewLoginUc(userRepository UserRepository) *LoginUc {
	return &LoginUc{userRepository: userRepository}
}

func (l *LoginUc) Login(ctx context.Context, email, password string) (*User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email or password is empty")
	}

	user, err := l.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("password for email %s isn't valid", email)
	}

	return user, nil
}
