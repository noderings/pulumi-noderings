# Changelog

## [0.1.2](https://github.com/noderings/pulumi-noderings/compare/v0.1.1...v0.1.2) (2026-11-26)


### Documentation

* add community health files ([de3e2f0](https://github.com/noderings/pulumi-noderings/commit/de3e2f0ff48af125ed08d0fe06fd34a11032117e))
* add README badges and a link to the how-to guide ([9371366](https://github.com/noderings/pulumi-noderings/commit/9371366e3cb68410780f3249ad8c455431c9713c))
* describe NodeRings as a marketplace for any infrastructure ([#1](https://github.com/noderings/pulumi-noderings/issues/1)) ([032d029](https://github.com/noderings/pulumi-noderings/commit/032d029b91e2f7e52aebef6c6e16eb3b5dcaf07e))

## 0.1.1

- The provider identifies itself to the NodeRings API as `pulumi-noderings/<version>` instead of the Terraform provider's User-Agent.
- The overview page follows the Pulumi Registry layout: installation for every language, a minimal example, and the configuration reference.
- The schema logo is now an SVG wordmark.

## 0.1.0

First release, built from the NodeRings Terraform provider.

- 23 resources: virtual machines with backups, snapshots and extra IPs; DNS zones, records, GeoDNS records and DNSSEC; firewall rulesets, rules, attachments and address groups; game servers with allocations, schedules, backups and databases; SSH keys, roles, role assignments, groups, group memberships, secrets and environment variables.
- 35 functions for the marketplace catalog, distributions, instant apps and account lookups.
- SDKs for TypeScript/JavaScript, Python, Go and .NET; YAML works without an SDK.
