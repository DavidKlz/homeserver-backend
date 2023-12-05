package types

import "github.com/google/uuid"

type MetaInfo struct {
	ID   uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name string    `json:"name" gorm:"index:idx_meta_info,unqiue"`
	Type string    `json:"type" gorm:"index:idx_meta_info,unqiue"`
}
