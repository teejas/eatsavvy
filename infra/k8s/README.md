# EatSavvy Kubernetes Manifests

Kubernetes manifests for deploying the EatSavvy application stack.

## Safety Check

The `deploy.sh` script includes a safety check that **only allows deployment to contexts containing "eatsavvy"** in the name. This prevents accidental deployment to the wrong cluster.

## Prerequisites

1. OKE cluster created (via Terraform in `../terraform/oci/`)
2. kubectl configured with the cluster
3. Context renamed to include "eatsavvy"

## Quick Start

```bash
# 1. After creating OKE cluster, configure kubectl
oci ce cluster create-kubeconfig \
  --cluster-id <cluster_id> \
  --file ~/.kube/eatsavvy.config \
  --region us-sanjose-1 \
  --token-version 2.0.0

# 2. Rename context to include "eatsavvy" (REQUIRED for safety check)
KUBECONFIG=~/.kube/eatsavvy.config kubectl config rename-context \
  $(KUBECONFIG=~/.kube/eatsavvy.config kubectl config current-context) \
  eatsavvy-cluster

# 3. Create secrets file (gitignored - won't be committed)
cp 01-secrets.yaml.example 01-secrets.yaml
# Edit 01-secrets.yaml with your credentials

# 4. Deploy
OCI_CLI_PROFILE={YOUR_PROFILE} ./deploy.sh apply
```

## Files

| File | Description |
|------|-------------|
| `deploy.sh` | Deployment script with safety check |
| `00-namespace.yaml` | Namespace |
| `01-secrets.yaml.example` | Secrets template (copy to 01-secrets.yaml) |
| `01-secrets.yaml` | Your secrets (gitignored) |
| `02-rabbitmq.yaml` | RabbitMQ Deployment, Service, PVC |
| `03-api.yaml` | API Deployment and Service |
| `04-worker.yaml` | Worker Deployment |
| `05-cloudflared.yaml` | Cloudflare Tunnel Deployment |

## Deploy Script Usage

```bash
# Deploy all manifests
./deploy.sh apply

# Delete all resources
./deploy.sh delete

# Check deployment status
./deploy.sh status

# Tail logs (default: api)
./deploy.sh logs
./deploy.sh logs worker
./deploy.sh logs rabbitmq
./deploy.sh logs cloudflared

# Run safety check only
./deploy.sh check
```

## Manual Deployment

If you prefer to apply manifests individually:

```bash
export KUBECONFIG=~/.kube/eatsavvy.config

# Verify context (must contain "eatsavvy")
kubectl config current-context

# Apply in order
kubectl apply -f 00-namespace.yaml
kubectl apply -f 01-secrets.yaml
kubectl apply -f 02-rabbitmq.yaml
kubectl apply -f 03-api.yaml
kubectl apply -f 04-worker.yaml
kubectl apply -f 05-cloudflared.yaml
```

## Secrets

The `01-secrets.yaml` file is **gitignored** to prevent committing credentials.

To set up:
```bash
cp 01-secrets.yaml.example 01-secrets.yaml
# Edit 01-secrets.yaml with your actual values
```

Required secrets:
- **OCIR credentials** - For pulling images from Oracle Container Registry
- **RabbitMQ credentials** - Username/password for RabbitMQ
- **Cloudflare tunnel token** - Your Cloudflare Tunnel token

## Troubleshooting

### Safety check failed

```
SAFETY CHECK FAILED!
Context 'context-xyz' does not contain 'eatsavvy'
```

Rename your context:
```bash
KUBECONFIG=~/.kube/eatsavvy.config kubectl config rename-context \
  context-xyz eatsavvy-cluster
```

### 01-secrets.yaml not found

Create it from the example:
```bash
cp 01-secrets.yaml.example 01-secrets.yaml
```

### Image pull errors

Check your OCIR credentials in `01-secrets.yaml`.

### RabbitMQ PVC pending

Verify the storage class exists:
```bash
kubectl get storageclass
```

### Cloudflared not connecting

Check the tunnel token in `01-secrets.yaml` and verify logs:
```bash
./deploy.sh logs cloudflared
```
