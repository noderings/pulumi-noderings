# NodeRings.Pulumi

The NodeRings provider for [Pulumi](https://www.pulumi.com) lets you manage [NodeRings](https://noderings.com) resources as part of your Pulumi programs. NodeRings is an infrastructure marketplace where independent hosting providers sell compute, GPU, storage, network and more through one account and one API.

With this package you can order a virtual machine by size and city, keep DNS zones, records, GeoDNS answers and DNSSEC in code, put firewalls in front of your servers, and manage SSH keys, roles, groups and secrets.

## Install

```bash
dotnet add package NodeRings.Pulumi
```

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
| `noderings:token` | `NODERINGS_TOKEN` | Service account token. Required. `NODERINGS_API_TOKEN` and `NR_API_TOKEN` also work. |
| `noderings:organizationId` | `NODERINGS_ORGANIZATION_ID` | Organization to act in. Defaults to the token's organization. |
| `noderings:endpoint` | `NODERINGS_API_URL` | API base URL. Defaults to the public API. |
| `noderings:maxRetries` | | Retries for rate-limited and transient failures. Defaults to `25`. |
| `noderings:requestsPerSecond` | | Client-side limit per API operation. Defaults to `2`. |
| `noderings:httpTimeout` | | Timeout for one HTTP request. Defaults to `60s`. |
| `noderings:skipOrganizationCheck` | | Skip the token check at startup. Defaults to `false`. |

## Example

Order the cheapest hourly 1 vCPU, 1 GB server in Ashburn and point a DNS name at it.

```csharp
using System.Collections.Generic;
using Pulumi;
using NodeRings.Pulumi;

return await Deployment.RunAsync(() =>
{
    var offer = GetMarketplaceOffer.Invoke(new()
    {
        PlanCode = "SC-GC-1C-1G-20",
        Country = "United States",
        City = "Ashburn",
        BillingType = "usage",
    });

    var key = new SshKey("laptop", new()
    {
        Name = "laptop",
        PublicKey = "ssh-ed25519 AAAA... you@laptop",
    });

    var web = new VirtualMachine("web", new()
    {
        Hostname = "web-1",
        PlanId = offer.Apply(o => o.PlanId),
        RegionId = offer.Apply(o => o.RegionId),
        Distro = "ubuntu_24_04",
        SshKeyIds = new[] { key.Id },
    });

    var zone = new DnsZone("site", new() { Name = "example.com" });

    new DnsRecord("www", new()
    {
        ZoneId = zone.Id,
        Name = "www",
        Type = "A",
        Ttl = 300,
        Records = new[] { web.Ipv4 },
    });

    return new Dictionary<string, object?>
    {
        ["ip"] = web.Ipv4,
        ["nameservers"] = zone.Nameservers,
    };
});
```

## Links

- [NodeRings](https://noderings.com)
- [Source repository](https://github.com/noderings/pulumi-noderings)
- [API documentation](https://www.pulumi.com/registry/packages/noderings/api-docs/) on the Pulumi Registry
- [NodeRings documentation](https://docs.noderings.com)

## License

Apache License 2.0. See [LICENSE](https://github.com/noderings/pulumi-noderings/blob/main/LICENSE).
