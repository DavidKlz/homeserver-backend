package storage

import (
	"fmt"

	"github.com/DavidKlz/homeserver-backend/types"
	"github.com/google/uuid"
)

func (s *PostgresStore) SaveUser(u *types.User) error {
	return s.DB.Save(u).Error
}

func (s *PostgresStore) DeleteUser(id *uuid.UUID) error {
	return s.DB.Delete(&types.User{}, id).Error
}

func (s *PostgresStore) GetUser(id *uuid.UUID) (*types.User, error) {
	user := &types.User{}
	res := s.DB.Find(&user, id)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("No User found for the given ID")
	}

	return user, nil
}

func (s *PostgresStore) GetAllUser() ([]types.User, error) {
	users := []types.User{}
	res := s.DB.Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

func (s *PostgresStore) FindUser(user, pass string) (string, error) {
	u := &types.User{}
	res := s.DB.Find(&u, "username = ? AND password = ?", user, pass)
	if res.Error != nil {
		return "", res.Error
	}

	if res.RowsAffected == 0 {
		return "", fmt.Errorf("no user found with name: %s", user)
	}

	return u.Token, res.Error
}

func (s *PostgresStore) FindUserByToken(token string) (*types.User, error) {
	u := &types.User{}
	res := s.DB.Find(&u, "token = ?", token)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("no user found with jwt token: %s", token)
	}

	return u, res.Error
}
