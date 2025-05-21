# helm-autoUpgrade

This repository provides a containerized tool to programmatically upgrade a Helm release using the Helm Go SDK. It supports custom `values.yaml` files and is suitable for automation in CI/CD pipelines.

## Features
- Upgrade any Helm release with a specified chart and values file
- Uses the official Helm Go SDK (no shelling out)
- Ready-to-use Docker image
- GitHub Actions workflow for CI and ECR push
- Kubernetes RBAC manifest for required permissions

## Usage

### Prerequisites
- Kubernetes cluster access (with permissions to upgrade releases)
- Helm 3 compatible chart
- Docker, AWS CLI (for ECR), and kubectl (for deployment)

### Build Locally
```sh
git clone <repo-url>
cd helm-autoUpgrade
go mod tidy
go build -o helm-upgrade main.go
```

### Run Locally
```sh
./helm-upgrade <release-name> <chart-path> <values.yaml>
```

### Build and Push Docker Image
```sh
docker build -t <your-ecr-repo>:<tag> .
docker push <your-ecr-repo>:<tag>
```

### GitHub Actions CI
- The workflow in `.github/workflows/ci.yml` builds and pushes the image to ECR on every push to `main`.
- Set up AWS credentials and ECR repository as required.

### Kubernetes RBAC
- Apply `k8s-role.yaml` to grant the pod permissions to perform Helm upgrades:
```sh
kubectl apply -f k8s-role.yaml
```

### Deploy as a Pod (example)
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: helm-upgrade
spec:
  serviceAccountName: default
  containers:
    - name: helm-upgrade
      image: <your-ecr-repo>:<tag>
      args: ["<release-name>", "<chart-path>", "<values.yaml>"]
      volumeMounts:
        - name: chart-volume
          mountPath: /charts
        - name: values-volume
          mountPath: /values
  volumes:
    - name: chart-volume
      configMap:
        name: your-chart-cm
    - name: values-volume
      configMap:
        name: your-values-cm
```

## Notes
- The tool uses the Helm Go SDK for all operations.
- Ensure the pod has access to the Kubernetes API and the necessary RBAC permissions.
- Customize the deployment and RBAC as needed for your environment.
