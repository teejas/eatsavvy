# -----------------------------------------------------------------------------
# OCI Infrastructure - VCN and OKE Cluster
# -----------------------------------------------------------------------------

# -----------------------------------------------------------------------------
# Network Module - VCN, Internet Gateway, Subnet
# -----------------------------------------------------------------------------

module "network" {
  source = "./modules/network"

  compartment_id    = var.compartment_id
  cluster_name      = var.cluster_name
  vcn_cidr          = var.vcn_cidr
  vcn_dns_label     = var.vcn_dns_label
  nodes_subnet_cidr = var.nodes_subnet_cidr
}

# -----------------------------------------------------------------------------
# OKE Module - Kubernetes Cluster and Node Pool
# -----------------------------------------------------------------------------

module "oke" {
  source = "./modules/oke"

  compartment_id      = var.compartment_id
  cluster_name        = var.cluster_name
  vcn_id              = module.network.vcn_id
  nodes_subnet_id     = module.network.nodes_subnet_id
  k8s_version         = var.k8s_version
  availability_domain = var.availability_domain
  node_shape          = var.node_shape
  node_ocpus          = var.node_ocpus
  node_memory_gb      = var.node_memory_gb
  node_count          = var.node_count
  node_image_id       = var.node_image_id

  depends_on = [module.network]
}

