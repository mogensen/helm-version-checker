package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/mogensen/helm-version-checker/pkg/metrics"
	"github.com/mogensen/helm-version-checker/pkg/models"
	"github.com/sirupsen/logrus"
)

// Controller probes certificates and registers the result in the metrics server
type Controller struct {
	log *logrus.Entry

	metrics      *metrics.Metrics
	helmReleases map[string]models.HelmRelease
	interval     time.Duration
	helm         helmService
}

// New returns a new configured instance of the Controller struct
func New(interval time.Duration, servingAddress string, log *logrus.Entry) *Controller {
	metrics := metrics.New(log)
	if err := metrics.Run(servingAddress); err != nil {
		log.Errorf("failed to start metrics server: %s", err)
		return nil
	}
	return &Controller{
		helmReleases: map[string]models.HelmRelease{},
		metrics:      metrics,
		interval:     interval,
		log:          log,
		helm: &helmServiceInst{
			log: log,
		},
	}
}

// HelmRels exposes helm release info to external services
func (c *Controller) HelmRels() []models.HelmRelease {
	r := []models.HelmRelease{}
	for _, rel := range c.helmReleases {
		r = append(r, rel)
	}

	return r
}

// Run starts the main loop that will call ProbeAll regularly.
func (c *Controller) Run(ctx context.Context) error {

	// First init helm repoes
	err := c.init(ctx)
	if err != nil {
		return err
	}

	// Start by probing all certificates before starting the ticker
	c.probeAll(ctx)

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		//select as usual
		select {
		case <-ctx.Done():
			c.log.Info("Stopping controller..")
			return nil
		case <-ticker.C:
			//give priority to a possible concurrent Done() event non-blocking way
			select {
			case <-ctx.Done():
				return nil
			default:
			}
			c.probeAll(ctx)
		}
	}
}

// probeAll triggers the Probe function for each registered service in the manager.
// Everything is done asynchronously.
func (c *Controller) probeAll(ctx context.Context) {
	c.log.Debug("Probing all")

	res, err := c.helm.probe()
	if err != nil {
		// If we get an error, we just ignore this probing
		return
	}

	for _, rel := range res.Releases {
		id := fmt.Sprintf("%s/%s", rel.Namespace, rel.Name)
		c.log.Debugf("Found %s", id)

		c.helmReleases[id] = rel
		c.metrics.AddHelmReleaseInfo(rel)
	}
}

func (c *Controller) init(ctx context.Context) error {
	c.log.Debug("Initial helm repo add")
	return c.helm.init()
}

// Shutdown closes the metrics server gracefully
func (c *Controller) Shutdown() error {
	// If metrics server is not started than exit early
	if c.metrics == nil {
		return nil
	}

	c.log.Info("shutting down metrics server...")

	if err := c.metrics.Shutdown(); err != nil {
		return fmt.Errorf("metrics server shutdown failed:%s", err)
	}

	c.log.Info("metrics server gracefully stopped")

	return nil
}
