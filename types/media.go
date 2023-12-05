package types

import "github.com/google/uuid"

type Media struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name            string    `json:"name"`
	Favorite        bool      `json:"favorite"`
	Type            string    `json:"type"`
	PathToFile      string
	PathToThumbnail string
}

type MediaReturn struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Favorite bool      `json:"favorite"`
	Type     string    `json:"type"`
}
