# -----------------------------------------------------------------------------
# OKE Module Outputs
# -----------------------------------------------------------------------------

output "cluster_id" {
  description = "OCID of the OKE cluster"
  value       = oci_containerengine_cluster.main.id
}

output "cluster_name" {
  description = "Name of the OKE cluster"
  value       = oci_containerengine_cluster.main.name
}

output "cluster_kubernetes_version" {
  description = "Kubernetes version of the cluster"
  value       = oci_containerengine_cluster.main.kubernetes_version
}

output "node_pool_id" {
  description = "OCID of the node pool"
  value       = oci_containerengine_node_pool.main.id
}

output "cluster_endpoints" {
  description = "Cluster endpoint information"
  value = {
    public_endpoint = oci_containerengine_cluster.main.endpoints[0].public_endpoint
  }
}

