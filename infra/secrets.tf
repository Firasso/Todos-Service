resource "null_resource" "env-secrets" {
  provisioner "local-exec" {
    command = "cd .. && make secrets"
  }

  depends_on = [
    module.eks,
    null_resource.aws_auth
  ]
}
