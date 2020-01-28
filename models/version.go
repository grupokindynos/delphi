package models

type VersionResponse struct {
	LatestVersion int `json:"latest_version"`
	MinVersion    int `json:"min_version"`
	SystemVersion int `json:"system_version"`
}
