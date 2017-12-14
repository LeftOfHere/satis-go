package db

// SatisDB needs a comment
type SatisDB struct {
	Name         string            `json:"name"`
	Homepage     string            `json:"homepage"`
	Repositories []SatisRepository `json:"repositories"`
	RequireAll   bool              `json:"require-all"`
}

// SatisRepository needs a comment
type SatisRepository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}
