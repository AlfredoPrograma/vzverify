terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.14.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}


# Create remote S3 bucket

resource "aws_s3_bucket" "remote_state_bucket" {
  bucket = "vzverify-remote-state"

  tags = {
    terraform = true
  }
}

resource "aws_s3_bucket_versioning" "remote_state_bucket_versioning" {
  bucket = aws_s3_bucket.remote_state_bucket.id

  versioning_configuration {
    status = "Enabled"
  }
}
