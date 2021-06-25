package controller

import (
	"context"
	"testing"
	"time"

	"github.com/mogensen/helm-version-checker/pkg/metrics"
	"github.com/mogensen/helm-version-checker/pkg/models"
	"github.com/sirupsen/logrus"
)

func TestController_Run_StopsWhenContextIsCanceled(t *testing.T) {
	type fields struct {
		log          *logrus.Entry
		metrics      *metrics.Metrics
		helmReleases []models.HelmRelease
		interval     time.Duration
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		interval time.Duration
	}{
		{interval: time.Millisecond * 100},
		{interval: time.Hour * 1},
		{interval: time.Second * 2},
	}
	for _, tt := range tests {
		t.Run("TestController_Run_StopsWhenContextIsCanceled "+tt.interval.String(), func(t *testing.T) {

			log := logrus.NewEntry(logrus.New())
			c := New(tt.interval, "0.0.0.0:0", log)
			c.helm = mockProber{}

			timeout := time.After(2 * time.Second)
			done := make(chan bool)
			go func() {
				ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
				if err := c.Run(ctx); err != nil {
					t.Errorf("Controller.Run() error = %v", err)
				}
				done <- true
				c.Shutdown()
			}()

			select {
			case <-timeout:
				t.Fatal("Test didn't finish in time")
			case <-done:
			}

		})
	}
}

// mockProber used in unit tests, to decouple tests from
// actual helm and kubectl invocations.
type mockProber struct{}

func (h mockProber) init() error {
	return nil
}

func (h mockProber) probe() (*models.WhatupResult, error) {
	time.Sleep(time.Millisecond * 100)
	return &models.WhatupResult{
		Releases: []models.HelmRelease{
			{
				Name:      "releaseName",
				Namespace: "default",
			},
		},
	}, nil
}
