package main

import (
	"context"
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
)

var settings = cli.New()

// upgradeHelmChart upgrades a Helm release with a given chart and values file
func upgradeHelmChart(releaseName, chartPath, valuesFile string) error {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), fmt.Printf); err != nil {
		return fmt.Errorf("failed to init helm action config: %w", err)
	}

	upgrade := action.NewUpgrade(actionConfig)
	chart, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("failed to load chart: %w", err)
	}

	p := values.Options{ValueFiles: []string{valuesFile}}
	vals, err := p.MergeValues(settings)
	if err != nil {
		return fmt.Errorf("failed to merge values: %w", err)
	}

	_, err = upgrade.Run(releaseName, chart, vals)
	if err != nil {
		return fmt.Errorf("helm upgrade failed: %w", err)
	}
	return nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: ./helm-upgrade <release-name> <chart-path> <values.yaml>")
		os.Exit(1)
	}
	releaseName := os.Args[1]
	chartPath := os.Args[2]
	valuesFile := os.Args[3]

	err := upgradeHelmChart(releaseName, chartPath, valuesFile)
	if err != nil {
		fmt.Printf("Helm upgrade failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Helm upgrade successful!")
}
