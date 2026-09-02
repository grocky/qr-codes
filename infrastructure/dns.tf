# wifi-signs.grocky.net -> the Pages project's *.pages.dev hostname.
# Cloudflare validates subdomain custom domains over the CNAME itself,
# so external (Route53) DNS works without moving the zone.
resource "aws_route53_record" "wifi_signs" {
  zone_id = data.aws_route53_zone.root.zone_id
  name    = var.domain_name
  type    = "CNAME"
  ttl     = 300
  records = ["${cloudflare_pages_project.wifi_signs.name}.pages.dev"]
}
