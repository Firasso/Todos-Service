resource "aws_security_group" "sec_grp_rds" {
  name_prefix = "${module.eks.cluster_id}-"
  vpc_id      = module.vpc.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group_rule" "allow-workers-nodes-communications" {
  description              = "Allow worker nodes to communicate with database"
  from_port                = 5432
  protocol                 = "tcp"
  security_group_id        = aws_security_group.sec_grp_rds.id
  source_security_group_id = module.eks.worker_security_group_id
  to_port                  = 5432
  type                     = "ingress"
}



module "db" {
  source = "terraform-aws-modules/rds/aws"

  identifier = "demodb-postgres"

  engine            = "postgres"
  engine_version    = "9.6.9"
  instance_class    = "db.t2.micro"
  allocated_storage = 5
  storage_encrypted = false
  apply_immediately = true

  # kms_key_id        = "arm:aws:kms:<region>:<account id>:key/<kms key id>"
  name = "todos"

  # NOTE: Do NOT use 'user' as the value for 'username' as it throws:
  # "Error creating DB Instance: InvalidParameterValue: MasterUsername
  # user cannot be used as it is a reserved word used by the engine"
  username = "demouser"
  password = "YourPwdShouldBeLongAndSecure!"
  port     = "5432"


  maintenance_window = "Mon:00:00-Mon:03:00"
  backup_window      = "03:00-06:00"

  # disable backups to create DB faster
  backup_retention_period = 0

  tags = {
    Owner       = "user"
    Environment = "dev"
  }

  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]

  vpc_security_group_ids = ["${aws_security_group.sec_grp_rds.id}"]

  # DB subnet group
  subnet_ids = module.vpc.database_subnets

  # DB parameter group
  family = "postgres9.6"

  # DB option group
  major_engine_version = "9.6"

  # Snapshot name upon DB deletion
  final_snapshot_identifier = "demo"

  # Database Deletion Protection
  deletion_protection = false
}
