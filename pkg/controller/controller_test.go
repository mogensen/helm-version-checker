package controller

import (
	"context"
	"testing"
	"time"

	"github.com/mogensen/helm-version-checker/pkg/metrics"
	"github.com/mogensen/helm-version-checker/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestController_probeAll_canAddFromBothHelmAndWhatsUp(t *testing.T) {
	c := testController()
	assert.Empty(t, c.helmReleases)

	c.probeAll(context.Background())

	assert.EqualValues(t, 2, len(c.helmReleases))
}

func TestController_probeAll_canRemoveUninstalledHelmReleases(t *testing.T) {
	c := testController()
	c.helmReleases["default/removedRelease"] = models.HelmRelease{
		Id:               "default/removedRelease",
		Name:             "releaseName",
		Namespace:        "default",
		InstalledVersion: "1.0.0",
		LatestVersion:    "2.0.0",
	}

	assert.EqualValues(t, 1, len(c.HelmRels()))

	c.probeAll(context.Background())

	assert.EqualValues(t, 2, len(c.HelmRels()))
}

func testController() *Controller {
	log := logrus.NewEntry(logrus.New())
	return &Controller{
		log:          log,
		metrics:      metrics.New(log),
		helmReleases: make(map[string]models.HelmRelease),
		interval:     time.Microsecond,
		helm:         mockProber{},
	}
}

func TestController_Run_StopsWhenContextIsCanceled(t *testing.T) {
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

func (h mockProber) list() ([]*models.HelmRelease, error) {
	return []*models.HelmRelease{
		{
			Id:               "default/upToDateRelease",
			Name:             "upToDateRelease",
			Namespace:        "default",
			InstalledVersion: "1.0.0",
			LatestVersion:    "---",
		},
	}, nil
}

func (h mockProber) probe() ([]*models.HelmRelease, error) {
	time.Sleep(time.Millisecond * 100)
	return []*models.HelmRelease{
		{
			Id:               "default/releaseName",
			Name:             "releaseName",
			Namespace:        "default",
			InstalledVersion: "1.0.0",
			LatestVersion:    "2.0.0",
		},
	}, nil
}
