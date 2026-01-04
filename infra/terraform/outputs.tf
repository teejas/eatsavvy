# -----------------------------------------------------------------------------
# OCI Infrastructure Outputs
# -----------------------------------------------------------------------------

# -----------------------------------------------------------------------------
# OKE Cluster Outputs
# -----------------------------------------------------------------------------

output "cluster_id" {
  description = "OCID of the OKE cluster"
  value       = module.oke.cluster_id
}

output "cluster_name" {
  description = "Name of the OKE cluster"
  value       = module.oke.cluster_name
}

output "cluster_kubernetes_version" {
  description = "Kubernetes version of the OKE cluster"
  value       = module.oke.cluster_kubernetes_version
}

# -----------------------------------------------------------------------------
# Network Outputs
# -----------------------------------------------------------------------------

output "vcn_id" {
  description = "OCID of the VCN"
  value       = module.network.vcn_id
}

output "nodes_subnet_id" {
  description = "OCID of the nodes subnet"
  value       = module.network.nodes_subnet_id
}

# -----------------------------------------------------------------------------
# Kubeconfig Helper
# -----------------------------------------------------------------------------

output "kubeconfig_command" {
  description = "Command to configure kubectl for this cluster"
  value       = "oci ce cluster create-kubeconfig --cluster-id ${module.oke.cluster_id} --file $HOME/.kube/eatsavvy.config --region ${var.region} --token-version 2.0.0"
}

output "region" {
  description = "OCI region"
  value       = var.region
}

