terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket       = "grocky-tfstate"
    key          = "wifi-signs.grocky.net/terraform.tfstate"
    region       = "us-east-1"
    use_lockfile = true
    encrypt      = true
  }
}

provider "aws" {
  region = "us-east-1"

  default_tags {
    tags = {
      Project   = "qr-codes"
      ManagedBy = "terraform"
    }
  }
}

provider "cloudflare" {
  # Uses the CLOUDFLARE_API_TOKEN environment variable.
  # The token needs the "Cloudflare Pages: Edit" account permission.
}

# ====== Variables ======

variable "cloudflare_account_id" {
  description = "Cloudflare account that owns the Pages project"
  type        = string
}

variable "domain_name" {
  description = "Public hostname of the site"
  type        = string
  default     = "wifi-signs.grocky.net"
}

variable "pages_project_name" {
  description = "Cloudflare Pages project name (deploy target for wrangler)"
  type        = string
  default     = "wifi-signs"
}

# ====== Data Sources ======

# The grocky.net zone lives in Route53, owned by ddns-service
# (~/Projects/grocky/ddns-service/terraform). Look it up by name rather
# than coupling to that repo's state.
data "aws_route53_zone" "root" {
  name = "grocky.net."
}
