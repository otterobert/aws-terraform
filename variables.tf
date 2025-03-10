# variables.tf

variable "cluster_name" {
  description = "The name of the EKS cluster"
  type        = string
  default     = "david-eks-terraform-cluster" # You can set a default value or leave it blank for user input
}

variable "region" {
  description = "The AWS region where resources will be created"
  type        = string
  default     = "us-west-2" # Default region (you can change it or set it when running Terraform)
}
