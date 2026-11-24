# Changelog

## 0.1.1

- The provider identifies itself to the NodeRings API as `pulumi-noderings/<version>` instead of the Terraform provider's User-Agent.
- The overview page follows the Pulumi Registry layout: installation for every language, a minimal example, and the configuration reference.
- The schema logo is now an SVG wordmark.

## 0.1.0

First release, built from the NodeRings Terraform provider.

- 23 resources: virtual machines with backups, snapshots and extra IPs; DNS zones, records, GeoDNS records and DNSSEC; firewall rulesets, rules, attachments and address groups; game servers with allocations, schedules, backups and databases; SSH keys, roles, role assignments, groups, group memberships, secrets and environment variables.
- 35 functions for the marketplace catalog, distributions, instant apps and account lookups.
- SDKs for TypeScript/JavaScript, Python, Go and .NET; YAML works without an SDK.
