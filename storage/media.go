package storage

import (
	"fmt"

	"github.com/DavidKlz/homeserver-backend/types"
	"github.com/google/uuid"
)

func (s *PostgresStore) SaveMedia(m *types.Media) error {
	return s.DB.Save(m).Error
}

func (s *PostgresStore) SaveMediaType(mt *types.MediaType) error {
	return s.DB.Save(mt).Error
}

func (s *PostgresStore) DeleteMedia(id *uuid.UUID) error {
	return s.DB.Delete(&types.Media{}, id).Error
}

func (s *PostgresStore) DeleteMediaType(id string) error {
	return s.DB.Delete(&types.MediaType{}, id).Error
}

func (s *PostgresStore) GetMedia(id *uuid.UUID) (*types.Media, error) {
	media := &types.Media{}
	res := s.DB.Find(&media, id)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("No Media found for the given ID")
	}

	return media, nil
}

func (s *PostgresStore) GetMediaType(name string) (*types.MediaType, error) {
	mt := &types.MediaType{}
	res := s.DB.Find(&mt, "Name = ?", name)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("No Media type found for the given ID")
	}

	return mt, nil
}

func (s *PostgresStore) GetAllMedia() ([]types.Media, error) {
	media := []types.Media{}
	res := s.DB.Find(&media)
	if res.Error != nil {
		return nil, res.Error
	}

	return media, nil
}

func (s *PostgresStore) GetAllMediaTypes() ([]types.MediaType, error) {
	mt := []types.MediaType{}
	res := s.DB.Find(&mt)
	if res.Error != nil {
		return nil, res.Error
	}

	return mt, nil
}
