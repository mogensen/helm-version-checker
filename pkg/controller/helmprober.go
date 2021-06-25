package controller

import (
	"encoding/json"
	"os/exec"

	"github.com/mogensen/helm-version-checker/pkg/models"
	"github.com/sirupsen/logrus"
)

type prober interface {
	Probe() (*models.WhatupResult, error)
}

type helmProber struct {
	log *logrus.Entry
}

func (h helmProber) Probe() (*models.WhatupResult, error) {
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
