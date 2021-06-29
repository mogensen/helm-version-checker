package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/mogensen/helm-version-checker/pkg/models"
)

type config struct {
	MetricsPort int `env:"METRICS_PORT" envDefault:"8080"`
	WebPort     int `env:"WEB_PORT" envDefault:"8081"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
	Repos    []models.Repo
}

func getRepoesFromEnv() []models.Repo {
	return parseRepoesFromStrings(os.Environ())
}

func parseRepoesFromStrings(repoes []string) []models.Repo {
	res := []models.Repo{}

	for _, element := range repoes {
		variable := strings.Split(element, "=")
		if strings.HasPrefix(variable[0], "HELM_REPO_") {
			fmt.Println(variable[0], "=>", variable[1])

			name := variable[1][:strings.Index(variable[1], ":")]
			url := variable[1][strings.Index(variable[1], ":")+1:]

			res = append(res, models.Repo{
				Name: name,
				Url:  url,
			})
		}
	}
	return res
}
