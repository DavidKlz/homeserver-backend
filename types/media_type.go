package types

const (
	ANIMATION = "Animation"
	VIDEO     = "Video"
	IMAGE     = "Image"
)

type MediaType struct {
	Name string `gorm:"primaryKey" json:"name"`
}
