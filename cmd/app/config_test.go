package app

import (
	"reflect"
	"testing"

	"github.com/mogensen/helm-version-checker/pkg/models"
)

func Test_parseRepoesFromStrings(t *testing.T) {
	type args struct {
		repoes []string
	}
	tests := []struct {
		name string
		args args
		want []models.Repo
	}{
		{
			name: "can handel empty",
			args: args{
				repoes: []string{},
			},
			want: []models.Repo{},
		},
		{
			name: "Full test",
			args: args{
				repoes: []string{
					"HOSTNAME=2f7eaa0b547c",
					"SHLVL=1",
					"HOME=/home/helm-version-checker",
					"TERM=xterm",
					"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
					"PWD=/app",
					"HELM_REPO_OTEEMOCHARTS=oteemocharts:https://oteemo.github.io/charts",
					"HELM_REPO_FALCOSECURITY=falcosecurity:https://falcosecurity.github.io/charts",
					"HELM_REPO_FAIRWINDS=fairwinds-stable:https://charts.fairwinds.com/stable",
				},
			},
			want: []models.Repo{
				{
					Name: "oteemocharts",
					Url:  "https://oteemo.github.io/charts",
				},
				{
					Name: "falcosecurity",
					Url:  "https://falcosecurity.github.io/charts",
				},
				{
					Name: "fairwinds-stable",
					Url:  "https://charts.fairwinds.com/stable",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseRepoesFromStrings(tt.args.repoes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRepoesFromStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}
