package models

import "time"

// HelmRelease stores info about helm releases
type HelmRelease struct {
	Id               string    `json:"id"`
	Name             string    `json:"name"`
	Namespace        string    `json:"namespace"`
	InstalledVersion string    `json:"installed_version"`
	LatestVersion    string    `json:"latest_version"`
	AppVersion       string    `json:"app_version"`
	Chart            string    `json:"chart"`
	NewestRepo       string    `json:"newest_repo"`
	Updated          time.Time `json:"updated"`
	Outdated         bool      `json:"outdated"`
	Deprecated       bool      `json:"deprecated"`
}
