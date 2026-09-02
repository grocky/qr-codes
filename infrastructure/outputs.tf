output "url" {
  description = "Public URL of the site"
  value       = "https://${var.domain_name}"
}

output "pages_url" {
  description = "Direct Cloudflare Pages URL"
  value       = "https://${cloudflare_pages_project.wifi_signs.name}.pages.dev"
}

output "pages_project" {
  description = "Pages project name for wrangler deploys"
  value       = cloudflare_pages_project.wifi_signs.name
}
