package web

import (
	"embed"
	"html/template"
	"io"
	"sort"

	"github.com/mogensen/helm-version-checker/pkg/models"
)

//go:embed views/*
var views embed.FS

// templateHTML generates an html representation for the given helm releases, and writes the result to the io.writer
func templateHTML(releases []models.HelmRelease, w io.Writer) error {

	sum := internalSummery{
		OutdatedReleases:    make(map[string][]uiHelmRelease),
		DeprecatedReleases:  make(map[string][]uiHelmRelease),
		MissingRepoReleases: make(map[string][]uiHelmRelease),
		GoodReleases:        make(map[string][]uiHelmRelease),
	}

	for _, c := range releases {
		uiC := uiHelmRelease{
			Name:             c.Name,
			Namespace:        c.Namespace,
			Deprecated:       c.Deprecated,
			InstalledVersion: c.InstalledVersion,
			LatestVersion:    c.LatestVersion,
			NewestRepo:       c.NewestRepo,
			Outdated:         c.Outdated,
			Chart:            c.Chart,
		}

		if uiC.Deprecated {
			sum.DeprecatedReleases[uiC.Namespace] = append(sum.DeprecatedReleases[uiC.Namespace], uiC)
		} else if uiC.NewestRepo == "---" {
			sum.MissingRepoReleases[uiC.Namespace] = append(sum.MissingRepoReleases[uiC.Namespace], uiC)
		} else if uiC.Outdated {
			sum.OutdatedReleases[uiC.Namespace] = append(sum.OutdatedReleases[uiC.Namespace], uiC)
		} else {
			sum.GoodReleases[uiC.Namespace] = append(sum.GoodReleases[uiC.Namespace], uiC)
		}
	}

	for i, v := range sum.DeprecatedReleases {
		sort.Sort(ByName(v))
		sum.DeprecatedReleases[i] = v
	}

	for i, v := range sum.MissingRepoReleases {
		sort.Sort(ByName(v))
		sum.MissingRepoReleases[i] = v
	}

	for i, v := range sum.OutdatedReleases {
		sort.Sort(ByName(v))
		sum.OutdatedReleases[i] = v
	}

	for i, v := range sum.GoodReleases {
		sort.Sort(ByName(v))
		sum.GoodReleases[i] = v
	}

	t := template.Must(template.New("index.html").Funcs(getFunctions()).ParseFS(views, "views/*"))
	err := t.Execute(w, sum)
	if err != nil {
		return err
	}
	return nil
}
