package web

import "time"

type uiHelmRelease struct {
	Name             string    `json:"name"`
	Namespace        string    `json:"namespace"`
	InstalledVersion string    `json:"installed_version"`
	LatestVersion    string    `json:"latest_version"`
	AppVersion       string    `json:"app_version"`
	Chart            string    `json:"chart"`
	NewestRepo       string    `json:"newest_repo"`
	Updated          time.Time `json:"updated"`
	Deprecated       bool      `json:"deprecated"`
	Outdated         bool
}

type internalSummery struct {
	OutdatedReleases    map[string][]uiHelmRelease
	DeprecatedReleases  map[string][]uiHelmRelease
	MissingRepoReleases map[string][]uiHelmRelease
	GoodReleases        map[string][]uiHelmRelease
}
