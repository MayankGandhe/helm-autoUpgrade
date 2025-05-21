# Usage Example: helm-autoUpgrade API

1. **Deploy the Service, Pod, and RBAC**
   - Apply the `podAndService.yml` manifest in this folder to your Kubernetes cluster:
   ```sh
   kubectl apply -f podAndService.yml
   ```
   This command will install the Service, Pod, Role, and RoleBinding, attaching the necessary permissions to the default service account.

2. **Trigger a Helm Upgrade/Install via API**
   - Once the service and pod are running, you can exec into any pod in the cluster and run the following curl command:
   ```sh
   curl -X POST http://helm-upgrade-svc:8080/upgrade \
     -H "Content-Type: application/json" \
     -d '{
       "releaseName": "nginx-remote",
       "chartURL": "https://charts.bitnami.com/bitnami/nginx-13.2.2.tgz",
       "valuesURL": "https://raw.githubusercontent.com/MayankGandhe/helm-autoUpgrade/refs/heads/main/test/new-value.yaml"
     }'
   ```
   - This will trigger a Helm upgrade (or install if the release does not exist). The release will be named `nginx-remote` and will use the provided chart and values file.
   - The values file can override settings such as the pod name using `nameOverride` or similar Helm values.

3. **Integration with Application Pod**
   - In a real deployment, your application pod can call this API service directly. The API will immediately acknowledge the request, allowing your application to proceed with its workflow while the upgrade/install runs in the background.

---

**Notes:**
- The API expects URLs for both the Helm chart (`chartURL`) and the values file (`valuesURL`).
- The service will handle both upgrades and fresh installs automatically.
- Ensure the pod running this service has the necessary RBAC permissions (see `podAndService.yml`).
