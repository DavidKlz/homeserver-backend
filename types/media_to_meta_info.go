package types

import "github.com/google/uuid"

type MediaToMetaInfo struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	MediaID    uuid.UUID `gorm:"index:idx_media_to_meta_info,unique"`
	Media      Media     `gorm:"foreignKey:id;references:media_id"`
	MetaInfoID uuid.UUID `gorm:"index:idx_media_to_meta_info,unique"`
	MetaInfo   MetaInfo  `gorm:"foreignKey:id;references:meta_info_id"`
}
