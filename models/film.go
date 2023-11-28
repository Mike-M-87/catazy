package models

type Film struct {
	Title    string `gorm:"primaryKey" json:"title"`
	Year     string `json:"year"`
	Quality  string `json:"quality"`
	Image    string `json:"image"`
	WatchUrl string `json:"url"`
	Duration string `json:"duration"`
	Watched  bool   `json:"watched"`
}