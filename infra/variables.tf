variable "resource_group_name" {
  description = " Name of the Azure resource group"
  type        = string
  default     = "ecom-tf-rg"
}

variable "location" {
  description = "Azure region for infra"
  type        = string
  default     = "Central India"
}

variable "network_interface_name" {
  description = "Name of the network"
  type        = string
  default     = "ecom-tf-nic"
}

variable "public_ip_name" {
  description = "Name of public ip"
  type        = string
  default     = "ecom-tf-public-ip"
}

variable "network_security_group_name" {
  description = "Name of network security group"
  type        = string
  default     = "ecom-tf-nsg"
}