# Description

Demo GO TODO list API

## Development

to run the development environment in a container with hot reload you can run :

```bash
make up
```

once you are done you can clean the container up by running

```bash
make down
```

## Infrastructure provisioning

you can use the development container which comes with aws , terraform , kubectl and other tools
to run the container

```shell
make infra
```

the command will run a terminal inside the container and mount the working directory

do the following steps in the terminal of the container

### Login to Amazon

```shell
# Access your "My Security Credentials" section in your profile.
# Create an access key

aws configure

Default region name: eu-west-1
Default output format: json
```

the container will save the credentials in the infra/.aws folder so you will only need to do this once

### Terraform

navigate into the infra folder and initialize terraform

```shell
cd infra && terraform init

terraform plan
terraform apply # takes around 20-25 minutes
```

### Lets deployed

this currently only deploys the API not the other services that are included in the dev container

navigate into k8s folder and deploy to kubernetes
note that you do not need to authenticate or add kube config .. it is done as a step of terraform automation

```shell
cd k8s && kubectl apply -f .

kubectl get pods
kubectl get svc
```

### Clean up

```shell
kubectl delete all --all
terraform destroy
```
