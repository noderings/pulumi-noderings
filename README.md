# NodeRings Resource Provider for Pulumi

[![CI](https://github.com/noderings/pulumi-noderings/actions/workflows/ci.yml/badge.svg)](https://github.com/noderings/pulumi-noderings/actions/workflows/ci.yml)
[![npm](https://img.shields.io/npm/v/@noderings/pulumi)](https://www.npmjs.com/package/@noderings/pulumi)
[![PyPI](https://img.shields.io/pypi/v/pulumi-noderings)](https://pypi.org/project/pulumi-noderings/)
[![NuGet](https://img.shields.io/nuget/v/NodeRings.Pulumi)](https://www.nuget.org/packages/NodeRings.Pulumi)
[![Go Reference](https://pkg.go.dev/badge/github.com/noderings/pulumi-noderings/sdk/go/noderings.svg)](https://pkg.go.dev/github.com/noderings/pulumi-noderings/sdk/go/noderings)
[![Pulumi Registry](https://img.shields.io/badge/pulumi-noderings-8A3391.svg?logo=pulumi)](https://www.pulumi.com/registry/packages/noderings/)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

The NodeRings provider lets you manage [NodeRings](https://noderings.com) resources with Pulumi. NodeRings is an infrastructure marketplace: independent hosting providers sell compute, GPU, storage, network and more, and you buy and manage them through one account and one API.

New to it? The step-by-step guide [How to use Pulumi on NodeRings to order a VPS](https://noderings.com/blog/pulumi-vps) walks through a first run.

With this provider you can:

- order a virtual machine by size and city, billed by the hour, and get its IP address back in the same program
- keep DNS zones, records, GeoDNS answers and DNSSEC in code
- put firewalls in front of your servers, matched by address, address group or tag
- manage SSH keys, roles, groups and secrets for your organization

It is built from the [NodeRings Terraform provider](https://github.com/noderings/terraform-provider-noderings) with the [Pulumi Terraform bridge](https://github.com/pulumi/pulumi-terraform-bridge), so both behave the same.

## Installing

| Language | Package | Install |
|---|---|---|
| TypeScript / JavaScript | [`@noderings/pulumi`](https://www.npmjs.com/package/@noderings/pulumi) | `npm install @noderings/pulumi` |
| Python | [`pulumi-noderings`](https://pypi.org/project/pulumi-noderings/) | `pip install pulumi-noderings` |
| Go | [`github.com/noderings/pulumi-noderings/sdk/go/noderings`](https://pkg.go.dev/github.com/noderings/pulumi-noderings/sdk/go/noderings) | `go get github.com/noderings/pulumi-noderings/sdk/go/noderings` |
| .NET | [`NodeRings.Pulumi`](https://www.nuget.org/packages/NodeRings.Pulumi) | `dotnet add package NodeRings.Pulumi` |
| YAML | built in | reference `noderings:VirtualMachine` and the other types directly |

Pulumi downloads the provider plugin from this repository's [releases](https://github.com/noderings/pulumi-noderings/releases) the first time a program uses it.

## Configuration

Create a service account in the NodeRings console (Organization, then Access Control, then Service accounts), give it a role, and generate a token. Then either export it:

```bash
export NODERINGS_TOKEN=<your token>
```

or store it encrypted in your stack:

```bash
pulumi config set noderings:token --secret
```

| Option | Environment variable | Description |
|---|---|---|
| `noderings:token` | `NODERINGS_TOKEN` | Service account token. Required. |
| `noderings:organizationId` | `NODERINGS_ORGANIZATION_ID` | Organization to act in. Defaults to the token's organization. |
| `noderings:endpoint` | `NODERINGS_API_URL` | API base URL. Defaults to the public API. |
| `noderings:maxRetries` | | Retries for rate-limited and transient failures. Defaults to `25`. |
| `noderings:requestsPerSecond` | | Client-side limit per API operation. Defaults to `2`. |
| `noderings:httpTimeout` | | Timeout for one HTTP request. Defaults to `60s`. |
| `noderings:skipOrganizationCheck` | | Skip the token check at startup. Defaults to `false`. |

## Example

Order the cheapest hourly 1 vCPU, 1 GB server in Ashburn and point a DNS name at it.

```typescript
import * as noderings from "@noderings/pulumi";

const offer = noderings.getMarketplaceOfferOutput({
    planCode: "SC-GC-1C-1G-20",
    country: "United States",
    city: "Ashburn",
    billingType: "usage",
});

const key = new noderings.SshKey("laptop", {
    name: "laptop",
    publicKey: "ssh-ed25519 AAAA... you@laptop",
});

const web = new noderings.VirtualMachine("web", {
    hostname: "web-1",
    planId: offer.planId,
    regionId: offer.regionId,
    distro: "ubuntu_24_04",
    sshKeyIds: [key.id],
});

const zone = new noderings.DnsZone("site", { name: "example.com" });

new noderings.DnsRecord("www", {
    zoneId: zone.id,
    name: "www",
    type: "A",
    ttl: 300,
    records: [web.ipv4],
});

export const ip = web.ipv4;
export const nameservers = zone.nameservers;
```

The same program in Python, Go, C# and YAML is in the [documentation](docs/_index.md). Complete, tested programs are in [`examples/`](examples):

| Example | Language | What it does | Cost |
|---|---|---|---|
| [`dns-ts`](examples/dns-ts) | TypeScript | Zone with A, MX, TXT, CAA and GeoDNS records, plus DNSSEC | Free |
| [`dns-go`](examples/dns-go) | Go | Zone with an A and an MX record | Free |
| [`vm-py`](examples/vm-py) | Python | Hourly VM with a firewall and a DNS record | A few cents an hour |

## Things to know

- **Billing.** `billingType: "usage"` servers are charged by the hour from prepaid credit and stop billing when deleted. Monthly servers are charged to your card and are cancelled without a refund when deleted, so set `deletionProtection: true` on them.
- **Replacements.** Changing `hostname`, `distro` or `sshKeyIds` replaces the server (new IP, empty disk) unless `reinstallOnChange: true`, which reinstalls the same server and keeps its addresses. Changing `planId` or `regionId` always replaces it.
- **Cheaper offers.** `getMarketplaceOffer` runs on every update. When a cheaper offer appears, `planId` changes and Pulumi plans a replacement. Use `ignoreChanges: ["planId", "regionId"]` or pin `providerId` to stay put.
- **Firewalls.** Rules are evaluated top to bottom and unmatched traffic is accepted, so end a ruleset with a drop rule.
- **Import.** Existing resources can be adopted with `pulumi import`, for example `pulumi import noderings:index/dnsRecord:DnsRecord www <zone id>/www/A`.

## Reference

- [API documentation](https://www.pulumi.com/registry/packages/noderings/api-docs/) on the Pulumi Registry
- [Provider overview, installation and configuration](docs/_index.md)
- [NodeRings documentation](https://docs.noderings.com)

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for building the provider, regenerating the SDKs, and running the unit and live integration tests.

## Support

Report bugs and request features in [GitHub issues](https://github.com/noderings/pulumi-noderings/issues). For questions about your NodeRings account or servers, contact NodeRings support from the console.

## Related

The [Go SDK](https://github.com/noderings/noderings-go) is the official Go client for the NodeRings API.

## License

Apache License 2.0. See [LICENSE](LICENSE) and [COPYRIGHT](COPYRIGHT).
