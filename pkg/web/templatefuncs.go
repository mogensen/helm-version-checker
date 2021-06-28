package web

import (
	"html/template"
	"time"
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

func count(releases map[string][]uiHelmRelease) int {
	res := 0
	for _, l := range releases {
		res += len(l)
	}

	return res
}
