resource "null_resource" "aws_auth" {
  provisioner "local-exec" {
    command = "aws eks update-kubeconfig --name ${var.cluster_name} --region ${var.region}"
  }

  depends_on = [
    module.eks,
  ]
}
