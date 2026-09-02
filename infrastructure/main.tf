terraform {
  required_version = ">= 1.5"

  required_providers {
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

provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

variable "cloudflare_api_token" {
  type        = string
  description = "Cloudflare API token with Pages edit + DNS edit permissions."
  sensitive   = true
}

variable "cloudflare_account_id" {
  type        = string
  description = "Cloudflare account ID owning the zone and Pages project."
}

variable "cloudflare_zone_id" {
  type        = string
  description = "Cloudflare zone ID for grocky.net."
}

variable "cloudflare_zone_name" {
  type        = string
  description = "Apex domain managed in Cloudflare."
  default     = "grocky.net"
}

variable "site_subdomain" {
  type        = string
  description = "Subdomain serving the site (joined with cloudflare_zone_name)."
  default     = "wifi-signs"
}

variable "pages_project_name" {
  type        = string
  description = "Cloudflare Pages project name (deploy target for wrangler)."
  default     = "wifi-signs"
}

locals {
  site_hostname = "${var.site_subdomain}.${var.cloudflare_zone_name}"
}
