# -----------------------------------------------------------------------------
# Network Module - VCN, Internet Gateway, Route Table, Subnet
# Creates the networking foundation for OKE cluster
# -----------------------------------------------------------------------------

terraform {
  required_providers {
    oci = {
      source = "oracle/oci"
    }
  }
}

# -----------------------------------------------------------------------------
# VCN - Virtual Cloud Network
# -----------------------------------------------------------------------------

resource "oci_core_vcn" "main" {
  compartment_id = var.compartment_id
  cidr_blocks    = [var.vcn_cidr]
  display_name   = "${var.cluster_name}-vcn"
  dns_label      = var.vcn_dns_label
}

# -----------------------------------------------------------------------------
# Internet Gateway - For outbound internet access
# -----------------------------------------------------------------------------

resource "oci_core_internet_gateway" "main" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.main.id
  display_name   = "${var.cluster_name}-igw"
  enabled        = true
}

# -----------------------------------------------------------------------------
# Route Table - Routes 0.0.0.0/0 via Internet Gateway
# -----------------------------------------------------------------------------

resource "oci_core_route_table" "public" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.main.id
  display_name   = "${var.cluster_name}-public-rt"

  route_rules {
    destination       = "0.0.0.0/0"
    destination_type  = "CIDR_BLOCK"
    network_entity_id = oci_core_internet_gateway.main.id
  }
}

# -----------------------------------------------------------------------------
# Security List - Allow necessary traffic for OKE nodes
# -----------------------------------------------------------------------------

resource "oci_core_security_list" "nodes" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.main.id
  display_name   = "${var.cluster_name}-nodes-sl"

  # Allow all egress traffic
  egress_security_rules {
    destination      = "0.0.0.0/0"
    protocol         = "all"
    destination_type = "CIDR_BLOCK"
    stateless        = false
  }

  # Allow all traffic within the VCN (for pod-to-pod, node-to-node communication)
  ingress_security_rules {
    source      = var.vcn_cidr
    protocol    = "all"
    source_type = "CIDR_BLOCK"
    stateless   = false
  }

  # Allow ICMP for path discovery
  ingress_security_rules {
    source      = "0.0.0.0/0"
    protocol    = "1" # ICMP
    source_type = "CIDR_BLOCK"
    stateless   = false

    icmp_options {
      type = 3
      code = 4
    }
  }

  # Allow SSH (optional, for node debugging)
  ingress_security_rules {
    source      = "0.0.0.0/0"
    protocol    = "6" # TCP
    source_type = "CIDR_BLOCK"
    stateless   = false

    tcp_options {
      min = 22
      max = 22
    }
  }

  # Allow NodePort range (not using LoadBalancer, but useful for debugging)
  ingress_security_rules {
    source      = "0.0.0.0/0"
    protocol    = "6" # TCP
    source_type = "CIDR_BLOCK"
    stateless   = false

    tcp_options {
      min = 30000
      max = 32767
    }
  }

  # Allow Kubernetes API access (port 6443)
  ingress_security_rules {
    source      = "0.0.0.0/0"
    protocol    = "6" # TCP
    source_type = "CIDR_BLOCK"
    stateless   = false

    tcp_options {
      min = 6443
      max = 6443
    }
  }
}

# -----------------------------------------------------------------------------
# Subnet - Public subnet for OKE nodes
# -----------------------------------------------------------------------------

resource "oci_core_subnet" "nodes" {
  compartment_id             = var.compartment_id
  vcn_id                     = oci_core_vcn.main.id
  cidr_block                 = var.nodes_subnet_cidr
  display_name               = "${var.cluster_name}-nodes-subnet"
  dns_label                  = "nodes"
  route_table_id             = oci_core_route_table.public.id
  security_list_ids          = [oci_core_security_list.nodes.id]
  prohibit_public_ip_on_vnic = false # Public subnet - nodes get public IPs
}

