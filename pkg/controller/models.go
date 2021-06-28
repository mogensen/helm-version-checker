package controller

import "github.com/mogensen/helm-version-checker/pkg/models"

// whatupResult is used for unmarshalling results from `helm whatup`
type whatupResult struct {
	Releases []models.HelmRelease `json:"releases"`
}

type release struct {
	// Name is the name of the release
	Name string `json:"name,omitempty"`
	// Chart is the chart that was released.
	Chart string `json:"chart,omitempty"`
	// AppVersion is the AppVersion that was released.
	AppVersion string `json:"app_version,omitempty"`
	// Namespace is the kubernetes namespace of the release.
	Namespace string `json:"namespace,omitempty"`
}
