provider "aws" {
  version = ">= 2.28.1"
  region  = var.region
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
  token                  = data.aws_eks_cluster_auth.cluster.token
  load_config_file       = false
  version                = "~> 1.11"
}

data "http" "workstation_external_ip" {
  url = "http://ipv4.icanhazip.com"
}

locals {
  workstation_external_cidr = "${chomp(data.http.workstation_external_ip.body)}/32"
}

data "aws_eks_cluster" "cluster" {
  name = module.eks.cluster_id
}

data "aws_eks_cluster_auth" "cluster" {
  name = module.eks.cluster_id
}

data "aws_availability_zones" "available" {
}


module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "2.6.0"

  name             = "test-vpc"
  cidr             = "10.0.0.0/16"
  azs              = data.aws_availability_zones.available.names
  private_subnets  = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets   = ["10.0.4.0/24", "10.0.5.0/24", "10.0.6.0/24"]
  database_subnets = ["10.0.7.0/24", "10.0.8.0/24", "10.0.9.0/24"]

  enable_nat_gateway           = true
  single_nat_gateway           = true
  enable_dns_hostnames         = true
  create_database_subnet_group = false

  public_subnet_tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                    = "1"
  }

  private_subnet_tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
    "kubernetes.io/role/internal-elb"           = "1"
  }
}


module "eks" {
  source          = "terraform-aws-modules/eks/aws"
  version         = "12.2.0"
  cluster_name    = var.cluster_name
  cluster_version = "1.17"

  vpc_id                          = module.vpc.vpc_id
  subnets                         = flatten([module.vpc.public_subnets, module.vpc.private_subnets])
  cluster_create_timeout          = "1h"
  cluster_endpoint_private_access = true

  worker_groups_launch_template = [
    {
      name          = "worker-group-1"
      instance_type = "t2.small"
      key_name      = aws_key_pair.mykey.key_name

      asg_desired_capacity = 1
      asg_max_size         = 2
      asg_min_size         = 1
      autoscaling_enabled  = true
      public_ip            = true

    }
  ]

  worker_additional_security_group_ids = [aws_security_group.all_worker_mgmt.id]
  map_roles                            = var.map_roles
  map_users                            = var.map_users
  map_accounts                         = var.map_accounts

  config_output_path = "./.kube/"
}





