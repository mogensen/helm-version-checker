package web

import (
	"html/template"
	"time"
)

func getFunctions() template.FuncMap {
	return template.FuncMap{
		"timeNow": templateTimeNow,
	}
}

func templateTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
