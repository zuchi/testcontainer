package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"testcontainer/use_case"
)

type UserRepository struct {
	conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{conn: conn}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user use_case.User) error {
	_, err := ur.conn.Exec(ctx, "insert into usertable(email, password) values ($1, $2)", user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*use_case.User, error) {
	u := use_case.User{}
	err := ur.conn.QueryRow(ctx, "select * from usertable where upper(email) = upper($1)", email).Scan(&u.Email, &u.Password)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if u.Email == "" {
		return nil, nil
	}

	return &u, nil
}
