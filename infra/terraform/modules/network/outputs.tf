# -----------------------------------------------------------------------------
# Network Module Outputs
# -----------------------------------------------------------------------------

output "vcn_id" {
  description = "OCID of the created VCN"
  value       = oci_core_vcn.main.id
}

output "nodes_subnet_id" {
  description = "OCID of the nodes subnet"
  value       = oci_core_subnet.nodes.id
}

output "vcn_cidr" {
  description = "CIDR block of the VCN"
  value       = oci_core_vcn.main.cidr_blocks[0]
}

output "internet_gateway_id" {
  description = "OCID of the Internet Gateway"
  value       = oci_core_internet_gateway.main.id
}

