# wifi-signs.grocky.net → the Pages project's *.pages.dev hostname.
# The zone lives on Cloudflare, so the record is proxied (orange-cloud) and
# Cloudflare serves the site at the edge once cloudflare_pages_domain
# attaches the hostname to the project.
resource "cloudflare_dns_record" "wifi_signs" {
  zone_id = var.cloudflare_zone_id
  name    = var.site_subdomain
  type    = "CNAME"
  content = "${cloudflare_pages_project.wifi_signs.name}.pages.dev"
  proxied = true
  ttl     = 1 # automatic; required for proxied records
}
