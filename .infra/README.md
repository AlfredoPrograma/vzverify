# Spinning up infrastructure

## Bootstrap
Must create the S3 bucket dedicated to store the state of the infrastructure remotely; so follow the next steps:

1. Move to `.infra/bootstrap` directory and run `terraform init` to initialize terraform dependencies.
2. Run `terraform apply` and approve it manually (or pass `--auto-approve` flag to avoid confirmation) to create the S3 bucket

## App

Once we set up the remote S3 bucket for state; now we can spin up the corresponding infrastructure:

1. Move to `.infra/app` directory and run `terraform init`. 
2. Fill the prompt with the path to key to target the corresponding environment state (`dev/terraform.tfstate`, `stg/terraform.tfstate` or `prod/terraform.tfstate` as standard values)
3. Run `terraform apply`