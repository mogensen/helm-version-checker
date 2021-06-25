package metrics

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mogensen/helm-version-checker/pkg/models"
	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics exposes helm release checks as prometheus metrics
type Metrics struct {
	*http.Server

	registry            *prometheus.Registry
	releaseIsLatest     *prometheus.GaugeVec
	releaseIsDeprecated *prometheus.GaugeVec
	log                 *logrus.Entry

	// release cache stores a cache of a helm release info
	releaseCache map[string]models.HelmRelease
	mu           sync.Mutex
}

// New returns a new configured instance of the Metrics server
func New(log *logrus.Entry) *Metrics {
	releaseIsLatest := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "helm_version_checker",
			Name:      "is_latest",
			Help:      "Detailing if the helm release uses the latest version of the helm chart",
		},
		[]string{
			"namespace", "name", "installed_version", "latest_version",
		},
	)
	releaseIsDeprecated := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "helm_version_checker",
			Name:      "is_deprecated",
			Help:      "Detailing if the helm release uses a deprecated helm chart",
		},
		[]string{
			"namespace", "name", "installed_version", "latest_version",
		},
	)

	registry := prometheus.NewRegistry()
	registry.MustRegister(releaseIsLatest)
	registry.MustRegister(releaseIsDeprecated)

	return &Metrics{
		log:                 log,
		registry:            registry,
		releaseIsLatest:     releaseIsLatest,
		releaseIsDeprecated: releaseIsDeprecated,
		releaseCache:        make(map[string]models.HelmRelease),
	}
}

// Run will run the metrics server
func (m *Metrics) Run(servingAddress string) error {
	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}))

	ln, err := net.Listen("tcp", servingAddress)
	if err != nil {
		return err
	}

	m.Server = &http.Server{
		Addr:           ln.Addr().String(),
		ReadTimeout:    8 * time.Second,
		WriteTimeout:   8 * time.Second,
		MaxHeaderBytes: 1 << 15, // 1 MiB
		Handler:        router,
	}

	go func() {
		m.log.Infof("serving metrics on %s/metrics", servingAddress)

		if err := m.Serve(ln); err != nil && !strings.Contains(err.Error(), "Server closed") {
			m.log.Errorf("failed to serve prometheus metrics: %s", err)
			return
		}
	}()

	return nil
}

// AddHelmReleaseInfo registers a new or updates and existing helm release record
func (m *Metrics) AddHelmReleaseInfo(rel models.HelmRelease) {
	// Remove old helm release information if it exists
	id := fmt.Sprintf("%s/%s", rel.Namespace, rel.Name)
	m.RemoveHelmReleaseInfo(id)

	m.mu.Lock()
	defer m.mu.Unlock()

	m.releaseCache[id] = rel

	isLatestF := 0.0
	if rel.InstalledVersion == rel.LatestVersion {
		isLatestF = 1.0
	}

	m.releaseIsLatest.With(
		m.buildLabelsLatest(rel),
	).Set(isLatestF)

	if rel.Deprecated {
		m.releaseIsDeprecated.With(
			m.buildLabelsLatest(rel),
		).Set(1.0)
	}
}

// RemoveHelmReleaseInfo removed an existing helm release record
func (m *Metrics) RemoveHelmReleaseInfo(dns string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.releaseCache[dns]
	if !ok {
		m.log.Debugf("Did not find %s in cache", dns)
		return
	}

	m.releaseIsLatest.Delete(m.buildLabelsLatest(item))
	m.releaseIsDeprecated.Delete(m.buildLabelsLatest(item))

	delete(m.releaseCache, dns)
}

func (m *Metrics) buildLabelsLatest(cer models.HelmRelease) prometheus.Labels {
	return prometheus.Labels{
		"namespace":         cer.Namespace,
		"name":              cer.Name,
		"installed_version": cer.InstalledVersion,
		"latest_version":    cer.LatestVersion,
	}
}

// Shutdown closes the metrics server gracefully
func (m *Metrics) Shutdown() error {
	// If metrics server is not started than exit early
	if m.Server == nil {
		return nil
	}

	m.log.Info("shutting down prometheus metrics server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := m.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("prometheus metrics server shutdown failed: %s", err)
	}

	m.log.Info("prometheus metrics server gracefully stopped")

	return nil
}
