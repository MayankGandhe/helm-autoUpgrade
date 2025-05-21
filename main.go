package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/getter"
)

var settings = cli.New()

// upgradeHelmChart upgrades a Helm release with a given chart and values file
func upgradeHelmChart(releaseName, chartPath, valuesFile string) error {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Printf(format, v...)
	}); err != nil {
		return fmt.Errorf("failed to init helm action config: %w", err)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("failed to load chart: %w", err)
	}

	p := values.Options{ValueFiles: []string{valuesFile}}
	vals, err := p.MergeValues(getter.All(settings))
	if err != nil {
		return fmt.Errorf("failed to merge values: %w", err)
	}

	upgrade := action.NewUpgrade(actionConfig)
	_, err = upgrade.Run(releaseName, chart, vals)
	if err == nil {
		return nil // upgrade successful
	}
	if err.Error() != "release: not found" && !isReleaseNotFound(err) {
		return fmt.Errorf("helm upgrade failed: %w", err)
	}
	// If release not found, do install
	install := action.NewInstall(actionConfig)
	install.ReleaseName = releaseName
	install.Namespace = settings.Namespace()
	_, err = install.Run(chart, vals)
	if err != nil {
		return fmt.Errorf("helm install failed: %w", err)
	}
	return nil
}

// Helper to check for not found error
func isReleaseNotFound(err error) bool {
	return err != nil && (err.Error() == "release: not found" ||
		len(err.Error()) >= 18 && err.Error()[:18] == "release: not found")
}

// Struct for upgrade request
// Accept chartURL and valuesURL instead of file paths
type UpgradeRequest struct {
	ReleaseName string `json:"releaseName"`
	ChartURL    string `json:"chartURL"`
	ValuesURL   string `json:"valuesURL"`
}

// Helper to download a file from URL to a temp file
func downloadToTemp(url, prefix string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	f, err := os.CreateTemp("", prefix+"-*")
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", err
	}
	return f.Name(), nil
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
	// Download chart and values file in background, then upgrade
	go func() {
		chartPath, err := downloadToTemp(req.ChartURL, "chart")
		if err != nil {
			fmt.Printf("Failed to download chart: %v\n", err)
			return
		}
		defer os.Remove(chartPath)
		valuesPath, err := downloadToTemp(req.ValuesURL, "values")
		if err != nil {
			fmt.Printf("Failed to download values: %v\n", err)
			return
		}
		defer os.Remove(valuesPath)
		err = upgradeHelmChart(req.ReleaseName, chartPath, valuesPath)
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
