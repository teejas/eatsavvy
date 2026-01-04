# EatSavvy Terraform Infrastructure

Terraform configuration for deploying OCI infrastructure (VCN and OKE cluster).

**Note:** Kubernetes applications are deployed separately using manifests in `./k8s/`.

## Architecture

```
┌───────────────────────────────────────────────────────────────┐
│                         OCI Tenancy                           │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │                    VCN (10.0.0.0/16)                    │  │
│  │  ┌───────────────────────────────────────────────────┐  │  │
│  │  │             Public Subnet (10.0.1.0/24)           │  │  │
│  │  │  ┌─────────────────────────────────────────────┐  │  │  │
│  │  │  │                OKE Cluster                  │  │  │  │
│  │  │  │  ┌───────────────────────────────────────┐  │  │  │  │
│  │  │  │  │         Node Pool (2 nodes)           │  │  │  │  │
│  │  │  │  │                                       │  │  │  │  │
│  │  │  │  │  ┌──────────┐       ┌──────────────┐  │  │  │  │  │
│  │  │  │  │  │   API    │       │   RabbitMQ   │  │  │  │  │  │
│  │  │  │  │  │  :8080   │──────►│ :5672/:15672 │  │  │  │  │  │
│  │  │  │  │  └────▲─────┘       └──────┬───────┘  │  │  │  │  │
│  │  │  │  │       │                    │          │  │  │  │  │
│  │  │  │  │       │                    ▼          │  │  │  │  │
│  │  │  │  │  ┌────┴─────┐       ┌──────────────┐  │  │  │  │  │
│  │  │  │  │  │Cloudflared│       │    Worker   │  │  │  │  │  │
│  │  │  │  │  │ (Tunnel) │       └──────┬───────┘  │  │  │  │  │
│  │  │  │  │  └────▲─────┘              │          │  │  │  │  │
│  │  │  │  └───────┼────────────────────┼──────────┘  │  │  │  │
│  │  │  └──────────┼────────────────────┼─────────────┘  │  │  │
│  │  └─────────────┼────────────────────┼────────────────┘  │  │
│  └────────────────┼────────────────────┼───────────────────┘  │
└───────────────────┼────────────────────┼──────────────────────┘
                    │                    │
                    │                    │
┌───────────────────┴───┐       ┌────────┴────────┐
│     Cloudflare        │       │  External APIs  │
│      Network          │       │ ─────────────── │
└───────────▲───────────┘       │ • Google Places │
            │                   │ • VAPI          │
┌───────────┴───────────┐       │ • Supabase      │
│       Internet        │       └─────────────────┘
│       (Users)         │
└───────────────────────┘
```

**Data Flow:**
1. Users access the app via Cloudflare
2. Cloudflare Tunnel routes traffic to API (no public IP needed)
3. API queues enrichment jobs to RabbitMQ
4. Worker consumes jobs, calls external APIs, stores results in Supabase

## Directory Structure

```
infra/
├── terraform/           # OCI Infrastructure (Terraform)
│   ├── main.tf
│   ├── providers.tf
│   ├── variables.tf
│   ├── outputs.tf
│   ├── terraform.tfvars
│   └── modules/
│       ├── network/     # VCN, IGW, Route Table, Subnet
│       └── oke/         # OKE Cluster, Node Pool
└── k8s/                 # Kubernetes Applications (Manifests)
    ├── deploy.sh        # Deployment script with safety check
    └── *.yaml           # K8s manifests
```

## Prerequisites

1. **OCI CLI** installed and configured
   ```bash
   brew install oci-cli
   oci setup config
   ```

2. **Terraform** >= 1.3.0
   ```bash
   brew install terraform
   ```

## Quick Start

```bash
cd infra/terraform

# Initialize
terraform init

# Review plan
terraform plan

# Apply
terraform apply
```

## After Cluster Creation

After the OKE cluster is created:

```bash
# 1. Configure kubectl (command shown in terraform output)
oci ce cluster create-kubeconfig \
  --cluster-id <cluster_id> \
  --file ~/.kube/eatsavvy.config \
  --region us-sanjose-1 \
  --token-version 2.0.0

# 2. Rename context to include "eatsavvy" (required for deploy script safety check)
KUBECONFIG=~/.kube/eatsavvy.config kubectl config rename-context \
  $(KUBECONFIG=~/.kube/eatsavvy.config kubectl config current-context) \
  eatsavvy-cluster

# 3. Deploy K8s applications
cd ../k8s
./deploy.sh apply
```

## Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `tenancy_ocid` | Yes | OCID of your OCI tenancy |
| `user_ocid` | Yes | OCID of the OCI user |
| `fingerprint` | Yes | Fingerprint of the API signing key |
| `private_key_path` | Yes | Path to the OCI API private key |
| `region` | Yes | OCI region (e.g., us-sanjose-1) |
| `compartment_id` | Yes | OCID of the compartment |
| `availability_domain` | Yes | AD for the node pool |
| `cluster_name` | No | Name of the cluster (default: eatsavvy-cluster) |
| `node_count` | No | Number of worker nodes (default: 2) |

## Outputs

| Output | Description |
|--------|-------------|
| `cluster_id` | OKE cluster OCID |
| `cluster_name` | Name of the cluster |
| `kubeconfig_command` | Command to configure kubectl |

## Destroying

```bash
# First, delete K8s resources
cd ../k8s
./deploy.sh delete

# Then, destroy OCI infrastructure
cd ../terraform
terraform destroy
```
