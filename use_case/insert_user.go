package use_case

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type InsertUc struct {
	userRepo UserRepository
}

func NewInsertUC(repository UserRepository) *InsertUc {
	return &InsertUc{
		userRepo: repository,
	}
}

func (i *InsertUc) InsertNewUser(ctx context.Context, user User) error {
	if user.Email == "" || user.Password == "" {
		return errors.New("all field is required to save the user")
	}

	uFound, err := i.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return err
	}

	if uFound != nil {
		return fmt.Errorf("there is an email %s stored in the database", user.Email)
	}

	pwd, err := cryptUserPassword(user.Password)
	user.Password = pwd

	err = i.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func cryptUserPassword(password string) (string, error) {
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(fromPassword), nil
}
