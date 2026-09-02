resource "cloudflare_dns_record" "wifi_signs" {
  zone_id = var.cloudflare_zone_id
  name    = var.site_subdomain
  type    = "CNAME"
  content = cloudflare_pages_project.wifi_signs.subdomain
  proxied = true
  ttl     = 1
}
