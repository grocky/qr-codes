resource "cloudflare_pages_project" "wifi_signs" {
  account_id        = var.cloudflare_account_id
  name              = var.pages_project_name
  production_branch = "main"
}

resource "cloudflare_pages_domain" "wifi_signs" {
  account_id   = var.cloudflare_account_id
  project_name = cloudflare_pages_project.wifi_signs.name
  name         = local.site_hostname
}
