package postgres

import (
	"context"
	"fmt"

	"github.com/algakz/banner_service/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type User struct {
	ID             int    `db:"id"`
	Username       string `db:"username"`
	HashedPassword string `db:"hashed_password"`
	Role           string `db:"role"`
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (username, hashed_password, access_level) VALUES ($1, $2, 2)"
	commandTag, err := r.db.Exec(ctx, query, user.Username, user.Password)
	if err != nil {
		logrus.Errorf("error executing insert query: %s", err.Error())
		return err
	}
	if commandTag.RowsAffected() != 1 {
		logrus.Errorf("expected one row to be affected, got %d", commandTag.RowsAffected())
		return err
	}
	return nil
}

func (r UserRepository) GetUser(
	ctx context.Context,
	username, password string,
) (*models.User, error) {
	query := "SELECT id, username, access_level FROM users WHERE username=$1 AND hashed_password=$2"
	var user User
	var access_level int

	err := r.db.QueryRow(ctx, query, username, password).
		Scan(&user.ID, &user.Username, &access_level)
	if err != nil {
		logrus.Errorf("error occured while Scanning row from db: %s", err.Error())
		return nil, err
	}

	getRoleQuery := "SELECT name as role FROM access_level WHERE id=$1"
	err = r.db.QueryRow(ctx, getRoleQuery, access_level).Scan(&user.Role)
	if err != nil {
		logrus.Errorf("error occured while Scanning row from db: %s", err.Error())
		return nil, err
	}
	logrus.Debugf("user from db: %s", user.Username)

	return DBUserToModelUser(user), nil
}

func DBUserToModelUser(u User) *models.User {
	return &models.User{
		Id:       fmt.Sprint(u.ID),
		Username: u.Username,
		Role:     u.Role,
	}
}
