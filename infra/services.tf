
resource "kubernetes_service" "db-service" {
  # this is the only K8s resource we will deploy with terraform
  # to makes the database available to pods with NS database
  # all other deployments ,services ..etc should be deployed using kubectl

  metadata {
    name = "database"
    labels = {
      app = "database"
    }
  }

  spec {
    selector = {
      app = "database"
    }
    external_name = module.db.this_db_instance_address

    type = "ExternalName"
  }
}
