# -----------------------------------------------------------------------------
# OKE Module - Kubernetes Cluster and Node Pool
# Creates an OKE cluster with a managed node pool
# -----------------------------------------------------------------------------

terraform {
  required_providers {
    oci = {
      source = "oracle/oci"
    }
  }
}

# -----------------------------------------------------------------------------
# Data source to get the latest Oracle Linux image for nodes
# -----------------------------------------------------------------------------

data "oci_core_images" "oracle_linux" {
  compartment_id           = var.compartment_id
  operating_system         = "Oracle Linux"
  operating_system_version = "8"
  shape                    = var.node_shape
  sort_by                  = "TIMECREATED"
  sort_order               = "DESC"

  # Only fetch if node_image_id is not provided
  count = var.node_image_id == "" ? 1 : 0
}

locals {
  # Use provided image ID or latest Oracle Linux image
  node_image_id = var.node_image_id != "" ? var.node_image_id : data.oci_core_images.oracle_linux[0].images[0].id
}

# -----------------------------------------------------------------------------
# OKE Cluster
# -----------------------------------------------------------------------------

resource "oci_containerengine_cluster" "main" {
  compartment_id     = var.compartment_id
  kubernetes_version = var.k8s_version
  name               = var.cluster_name
  vcn_id             = var.vcn_id

  # Cluster endpoint configuration - public for simplicity
  endpoint_config {
    is_public_ip_enabled = true
    subnet_id            = var.nodes_subnet_id
  }

  # Cluster options
  options {
    # Enable Kubernetes dashboard (deprecated but still available)
    add_ons {
      is_kubernetes_dashboard_enabled = false
      is_tiller_enabled               = false
    }

    # Network configuration
    kubernetes_network_config {
      pods_cidr     = var.pods_cidr
      services_cidr = var.services_cidr
    }

    # Persistent volume configuration
    persistent_volume_config {
      defined_tags  = {}
      freeform_tags = {}
    }

    # Service LB configuration (not using LoadBalancers - using Cloudflare Tunnel instead)
    service_lb_config {
      defined_tags  = {}
      freeform_tags = {}
    }
  }

  # Cluster type - BASIC is sufficient for most use cases
  type = "BASIC_CLUSTER"
}

# -----------------------------------------------------------------------------
# Node Pool
# -----------------------------------------------------------------------------

resource "oci_containerengine_node_pool" "main" {
  cluster_id         = oci_containerengine_cluster.main.id
  compartment_id     = var.compartment_id
  kubernetes_version = var.k8s_version
  name               = "${var.cluster_name}-pool"
  node_shape         = var.node_shape

  # Flex shape configuration
  node_shape_config {
    ocpus         = var.node_ocpus
    memory_in_gbs = var.node_memory_gb
  }

  # Node source (OS image)
  node_source_details {
    image_id    = local.node_image_id
    source_type = "IMAGE"
  }

  # Node placement configuration
  node_config_details {
    size = var.node_count

    placement_configs {
      availability_domain = var.availability_domain
      subnet_id           = var.nodes_subnet_id
    }
  }

  # Initial node labels (optional)
  initial_node_labels {
    key   = "app"
    value = "eatsavvy"
  }
}
