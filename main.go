package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

// Struct for upgrade request
type UpgradeRequest struct {
	ReleaseName string `json:"releaseName"`
	ChartPath   string `json:"chartPath"`
	ValuesFile  string `json:"valuesFile"`
}

func upgradeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	var req UpgradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"acknowledged"}`))
	// Run upgrade in background
	go func() {
		err := upgradeHelmChart(req.ReleaseName, req.ChartPath, req.ValuesFile)
		if err != nil {
			fmt.Printf("Helm upgrade failed: %v\n", err)
		} else {
			fmt.Println("Helm upgrade successful!")
		}
	}()
}

func main() {
	http.HandleFunc("/upgrade", upgradeHandler)
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
