package app

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	env "github.com/caarlos0/env/v6"
	"github.com/mogensen/helm-version-checker/pkg/controller"
	"github.com/mogensen/helm-version-checker/pkg/web"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	helpOutput = "Helm Chart version monitoring utility for watching updated and deprecated helm releases and reporting the result as metrics."
)

// NewCommand sets up the helm-version-checker command and all dependencies
func NewCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "helm-version-checker",
		Short: helpOutput,
		Long:  helpOutput,
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg := config{}
			if err := env.Parse(&cfg); err != nil {
				fmt.Printf("%+v\n", err)
			}

			nlog := logrus.New()
			nlog.SetOutput(os.Stdout)
			nlog.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

			log := logrus.NewEntry(nlog)

			logLevel, err := logrus.ParseLevel(cfg.LogLevel)
			if err != nil {
				return fmt.Errorf("failed to parse  loglevel %q: %s",
					cfg.LogLevel, err)
			}
			nlog.SetLevel(logLevel)

			log.Debugf("Config: %+v\n", cfg)

			// create a WaitGroup
			wg := new(sync.WaitGroup)
			wg.Add(2)

			// Metrics

			metricsAddress := fmt.Sprintf("%s:%d", "0.0.0.0", cfg.MetricsPort)
			c := controller.New(time.Minute, metricsAddress, log)

			go func() {
				<-ctx.Done()
				if err := c.Shutdown(); err != nil {
					log.Error(err)
				}
			}()

			go func() {
				c.Run(ctx)
				wg.Done()
			}()

			// Web UI

			webAddress := fmt.Sprintf("%s:%d", "0.0.0.0", cfg.WebPort)
			ui := web.New(c, webAddress, log)

			go func() {
				<-ctx.Done()
				if err := ui.Shutdown(); err != nil {
					log.Error(err)
				}
			}()

			go func() {
				ui.Run(ctx)
				wg.Done()
			}()

			// wait until WaitGroup is done
			wg.Wait()
			log.Infof("Everything is successfully stopped")

			return nil
		},
	}

	return cmd
}
