package controller

import (
	"encoding/json"
	"os/exec"

	"github.com/mogensen/helm-version-checker/pkg/models"
	"github.com/sirupsen/logrus"
)

type helmService interface {
	init() error
	probe() (*models.WhatupResult, error)
}

type helmServiceInst struct {
	log *logrus.Entry
}

func (h helmServiceInst) init() error {

	repoes := make(map[string]string)

	// TODO Make config!
	repoes["fairwinds-stable"] = "https://charts.fairwinds.com/stable"
	repoes["stable"] = "https://charts.helm.sh/stable"
	repoes["cert-checker"] = "https://mogensen.github.io/cert-checker"
	repoes["prometheus-community"] = "https://prometheus-community.github.io/helm-charts"

	for repoName, repoUrl := range repoes {
		cmd := exec.Command("helm", "repo", "add", repoName, repoUrl)
		stdout, err := cmd.Output()
		h.log.Debugf("helm repo add output: %s", string(stdout))
		if err != nil {
			h.log.Errorf("Error adding repo: %v", err.Error())
			return err
		}
	}
	return nil
}

func (h helmServiceInst) probe() (*models.WhatupResult, error) {
	prg := "helm"

	cmd := exec.Command(prg, "whatup", "-A", "-a", "--ignore-deprecation=false", "--ignore-repo=false", "-o", "json")
	stdout, err := cmd.Output()
	h.log.Debugf("json: %s", string(stdout))

	if err != nil {
		h.log.Errorf("Error probing: %v", err.Error())
		return nil, err
	}

	var res models.WhatupResult
	err = json.Unmarshal(stdout, &res)
	if err != nil {
		h.log.Errorf("Error unmarshaling json: %v", err.Error())
		return nil, err
	}
	return &res, nil
}
