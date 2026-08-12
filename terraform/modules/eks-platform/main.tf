variable "cluster_name" { type = string }

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 21.0"

  name               = var.cluster_name
  kubernetes_version = "1.33"

  endpoint_private_access = true
  endpoint_public_access  = false
  enable_irsa = true

  eks_managed_node_groups = {
    system = {
      min_size       = 2
      desired_size   = 2
      max_size       = 5
      instance_types = ["m7i.large"]
      labels         = { workload = "system" }
    }

    application = {
      min_size       = 3
      desired_size   = 3
      max_size       = 30
      instance_types = ["m7i.xlarge"]
      labels         = { workload = "application" }
    }
  }
}
