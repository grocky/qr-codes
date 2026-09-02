# Cloudflare Pages project for the static site in web/.
# Deployments are direct uploads via `make deploy` (wrangler), not git-driven,
# so no source/build configuration is declared here.
resource "cloudflare_pages_project" "wifi_signs" {
  account_id        = var.cloudflare_account_id
  name              = var.pages_project_name
  production_branch = "main"
}

resource "cloudflare_pages_domain" "wifi_signs" {
  account_id   = var.cloudflare_account_id
  project_name = cloudflare_pages_project.wifi_signs.name
  name         = var.domain_name
}
