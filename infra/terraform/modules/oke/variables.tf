# -----------------------------------------------------------------------------
# OKE Module Variables
# -----------------------------------------------------------------------------

variable "compartment_id" {
  description = "OCID of the compartment where the OKE cluster will be created"
  type        = string
}

variable "cluster_name" {
  description = "Name of the OKE cluster"
  type        = string
}

variable "vcn_id" {
  description = "OCID of the VCN for the cluster"
  type        = string
}

variable "nodes_subnet_id" {
  description = "OCID of the subnet for OKE nodes"
  type        = string
}

variable "k8s_version" {
  description = "Kubernetes version for the cluster"
  type        = string
  default     = "v1.34.1"
}

variable "availability_domain" {
  description = "Availability domain for the node pool"
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
  description = "OCID of the node image. If empty, latest Oracle Linux 8 image is used."
  type        = string
  default     = ""
}

variable "pods_cidr" {
  description = "CIDR block for Kubernetes pods"
  type        = string
  default     = "10.244.0.0/16"
}

variable "services_cidr" {
  description = "CIDR block for Kubernetes services"
  type        = string
  default     = "10.96.0.0/16"
}

