package web

import (
	"github.com/mogensen/helm-version-checker/pkg/models"
)

type internalSummery struct {
	OutdatedReleases    map[string][]models.HelmRelease
	DeprecatedReleases  map[string][]models.HelmRelease
	MissingRepoReleases map[string][]models.HelmRelease
	GoodReleases        map[string][]models.HelmRelease
}
