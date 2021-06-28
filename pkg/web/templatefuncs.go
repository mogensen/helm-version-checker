package web

import (
	"html/template"
	"time"

	"github.com/mogensen/helm-version-checker/pkg/models"
)

func getFunctions() template.FuncMap {
	return template.FuncMap{
		"timeNow": templateTimeNow,
		"count":   count,
	}
}

func templateTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func count(releases map[string][]models.HelmRelease) int {
	res := 0
	for _, l := range releases {
		res += len(l)
	}

	return res
}
