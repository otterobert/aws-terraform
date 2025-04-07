variable "region" {
  description = "The AWS region where resources will be created"
  type        = string
  default     = "us-west-2" # Default region (you can change it or set it when running Terraform)
}
