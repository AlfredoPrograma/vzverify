terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.14.0"
    }
  }

  backend "s3" {
    bucket = "vzverify-remote-state"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "us-east-1"
}

# S3 bucket for IDS

resource "aws_s3_bucket" "identities_docs" {
  bucket = "vzverify-identities-docs"

  tags = {
    terraform = true
  }
}

resource "aws_s3_bucket_versioning" "identities_docs_versioning" {
  bucket = aws_s3_bucket.identities_docs.id

  versioning_configuration {
    status = "Enabled"
  }
}
