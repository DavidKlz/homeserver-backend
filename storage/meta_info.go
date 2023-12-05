package storage

import (
	"fmt"

	"github.com/DavidKlz/homeserver-backend/types"
	"github.com/google/uuid"
)

func (s *PostgresStore) SaveMetaInfo(mi *types.MetaInfo) error {
	meta := &types.MetaInfo{}
	res := s.DB.Find(&meta, "name = ? AND type = ?", mi.Name, mi.Type)

	if res.RowsAffected != 0 {
		mi.ID = meta.ID
	}

	return s.DB.Save(mi).Error
}

func (s *PostgresStore) SaveMediaToMetaInfo(m2mi *types.MediaToMetaInfo) error {
	return s.DB.Save(m2mi).Error
}

func (s *PostgresStore) DeleteMetaInfo(id *uuid.UUID) error {
	return s.DB.Delete(&types.MetaInfo{}, id).Error
}

func (s *PostgresStore) DeleteMediaToMetaInfo(id *uuid.UUID) error {
	return s.DB.Delete(&types.MediaToMetaInfo{}, id).Error
}

func (s *PostgresStore) GetMetaInfo(id *uuid.UUID) (*types.MetaInfo, error) {
	mi := &types.MetaInfo{}
	res := s.DB.Find(&mi, id)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("No Meta Info found for the given ID")
	}

	return mi, nil
}

func (s *PostgresStore) GetMediaToMetaInfo(id *uuid.UUID) (*types.MediaToMetaInfo, error) {
	mi := &types.MediaToMetaInfo{}
	res := s.DB.Find(&mi, id)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("No Meta Info found for the given ID")
	}

	return mi, nil
}

func (s *PostgresStore) GetAllMetaInfo() ([]types.MetaInfo, error) {
	mi := []types.MetaInfo{}
	res := s.DB.Find(&mi)
	if res.Error != nil {
		return nil, res.Error
	}

	return mi, nil
}

func (s *PostgresStore) GetAllMediaToMetaInfos() ([]types.MediaToMetaInfo, error) {
	mi := []types.MediaToMetaInfo{}
	res := s.DB.Find(&mi)
	if res.Error != nil {
		return nil, res.Error
	}

	return mi, nil
}

func (s *PostgresStore) GetMetaInfosOfMedia(id *uuid.UUID) ([]types.MetaInfo, error) {
	m2mi := []types.MediaToMetaInfo{}
	res := s.DB.Find(&m2mi, "media_id = ?", id)
	if res.Error != nil {
		return nil, res.Error
	}

	mi := []types.MetaInfo{}

	for _, elem := range m2mi {
		mi = append(mi, elem.MetaInfo)
	}

	return mi, nil
}

func (s *PostgresStore) GetMediaOfMetaInfo(id *uuid.UUID) ([]types.Media, error) {
	m2mi := []types.MediaToMetaInfo{}
	res := s.DB.Find(&m2mi, "meta_info_id = ?", id)
	if res.Error != nil {
		return nil, res.Error
	}

	m := []types.Media{}

	for _, elem := range m2mi {
		m = append(m, elem.Media)
	}

	return m, nil
}
