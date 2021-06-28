package web

import "github.com/mogensen/helm-version-checker/pkg/models"

// ByName sorts HelmRelease'es by Repository and then image name
type ByName []models.HelmRelease

func (a ByName) Len() int { return len(a) }
func (a ByName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}
func (a ByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
