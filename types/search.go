package types

type Search struct {
	MetaInfo map[string][]string `json:"metaInfo"`
	Type     string              `json:"type"`
	Favorite bool                `json:"favorite"`
}
