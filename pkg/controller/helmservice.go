package controller

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mogensen/helm-version-checker/pkg/models"
	"github.com/sirupsen/logrus"
)

type helmService interface {
	init() error
	probe() ([]*models.HelmRelease, error)
	list() ([]*models.HelmRelease, error)
}

type helmServiceInst struct {
	log    *logrus.Entry
	repoes []models.Repo
}

func (h helmServiceInst) init() error {

	for _, repo := range h.repoes {
		cmd := exec.Command("helm", "repo", "add", repo.Name, repo.Url)
		stdout, err := cmd.Output()
		h.log.Debugf("helm repo add output: %s", string(stdout))
		if err != nil {
			h.log.Errorf("Error adding repo: %v", err.Error())
			return err
		}
	}
	return nil
}

func (h helmServiceInst) list() ([]*models.HelmRelease, error) {

	cmd := exec.Command("helm", "ls", "-A", "-a", "-o", "json")
	stdout, err := cmd.Output()
	h.log.Tracef("json: %s", string(stdout))

	if err != nil {
		h.log.Errorf("Error probing: %v", err.Error())
		return nil, err
	}

	var cliRes []*release
	err = json.Unmarshal(stdout, &cliRes)
	if err != nil {
		h.log.Errorf("Error unmarshaling json: %v", err.Error())
		return nil, err
	}

	var res []*models.HelmRelease
	for _, rel := range cliRes {

		installedVersion := rel.Chart[strings.LastIndex(rel.Chart, "-")+1:]
		chart := rel.Chart[:strings.LastIndex(rel.Chart, "-")]

		hRel := models.HelmRelease{
			Id:               fmt.Sprintf("%s/%s", rel.Namespace, rel.Name),
			Name:             rel.Name,
			Namespace:        rel.Namespace,
			Chart:            chart,
			InstalledVersion: installedVersion,
			AppVersion:       rel.AppVersion,
			NewestRepo:       "---",
			LatestVersion:    "---",
			Outdated:         false,
		}

		res = append(res, &hRel)
	}

	return res, nil
}

func (h helmServiceInst) probe() ([]*models.HelmRelease, error) {

	cmd := exec.Command("helm", "whatup", "-A", "-a", "--ignore-deprecation=false", "--ignore-repo=false", "-o", "json")
	stdout, err := cmd.Output()
	h.log.Tracef("json: %s", string(stdout))

	if err != nil {
		h.log.Errorf("Error probing: %v", err.Error())
		return nil, err
	}

	var cliRes whatupResult
	err = json.Unmarshal(stdout, &cliRes)
	if err != nil {
		h.log.Errorf("Error unmarshaling json: %v", err.Error())
		return nil, err
	}

	var res []*models.HelmRelease
	for _, rel := range cliRes.Releases {

		hRel := models.HelmRelease{
			Id:               fmt.Sprintf("%s/%s", rel.Namespace, rel.Name),
			Name:             rel.Name,
			Namespace:        rel.Namespace,
			Chart:            rel.Chart,
			InstalledVersion: rel.InstalledVersion,
			AppVersion:       rel.AppVersion,
			LatestVersion:    rel.LatestVersion,
			NewestRepo:       rel.NewestRepo,
			Outdated:         rel.InstalledVersion != rel.LatestVersion,
		}

		res = append(res, &hRel)
	}

	return res, nil
}
