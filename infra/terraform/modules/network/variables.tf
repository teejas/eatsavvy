# -----------------------------------------------------------------------------
# Network Module Variables
# -----------------------------------------------------------------------------

variable "compartment_id" {
  description = "OCID of the compartment where network resources will be created"
  type        = string
}

variable "cluster_name" {
  description = "Name prefix for network resources"
  type        = string
}

variable "vcn_cidr" {
  description = "CIDR block for the VCN"
  type        = string
  default     = "10.0.0.0/16"
}

variable "vcn_dns_label" {
  description = "DNS label for the VCN (alphanumeric, max 15 chars)"
  type        = string
  default     = "eatsavvyvcn"
}

variable "nodes_subnet_cidr" {
  description = "CIDR block for the OKE nodes subnet"
  type        = string
  default     = "10.0.20.0/24"
}

