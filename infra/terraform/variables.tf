# -----------------------------------------------------------------------------
# OCI Authentication Variables
# -----------------------------------------------------------------------------

variable "tenancy_ocid" {
  description = "OCID of the OCI tenancy"
  type        = string
}

variable "user_ocid" {
  description = "OCID of the OCI user for API authentication"
  type        = string
}

variable "fingerprint" {
  description = "Fingerprint of the API signing key"
  type        = string
}

variable "private_key_path" {
  description = "Path to the private key file for OCI API authentication"
  type        = string
}

variable "region" {
  description = "OCI region (e.g., us-ashburn-1)"
  type        = string
}

variable "compartment_id" {
  description = "OCID of the compartment where resources will be created"
  type        = string
}

# -----------------------------------------------------------------------------
# Network Variables
# -----------------------------------------------------------------------------

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

# -----------------------------------------------------------------------------
# OKE Cluster Variables
# -----------------------------------------------------------------------------

variable "cluster_name" {
  description = "Name of the OKE cluster"
  type        = string
  default     = "eatsavvy-cluster"
}

variable "k8s_version" {
  description = "Kubernetes version for the OKE cluster"
  type        = string
  default     = "v1.34.1"
}

variable "availability_domain" {
  description = "Availability domain for the node pool (e.g., 'Uocm:US-ASHBURN-AD-1')"
  type        = string
}

variable "node_shape" {
  description = "Shape for the worker nodes"
  type        = string
  default     = "VM.Standard.E4.Flex"
}

variable "node_ocpus" {
  description = "Number of OCPUs for flex shapes"
  type        = number
  default     = 1
}

variable "node_memory_gb" {
  description = "Memory in GB for flex shapes"
  type        = number
  default     = 16
}

variable "node_count" {
  description = "Number of worker nodes in the node pool"
  type        = number
  default     = 2
}

variable "node_image_id" {
  description = "OCID of the node image (Oracle Linux). If not provided, latest Oracle Linux image will be used."
  type        = string
  default     = ""
}

