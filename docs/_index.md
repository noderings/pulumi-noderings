---
title: NodeRings
meta_desc: Use the Pulumi NodeRings provider to order and manage compute, networking, DNS, access control and more on the NodeRings infrastructure marketplace.
layout: package
---

The NodeRings provider for Pulumi lets you manage [NodeRings](https://noderings.com) resources as part of your Pulumi programs. NodeRings is an infrastructure marketplace where independent hosting providers sell compute, GPU, storage, network and more through one account and one API. With this provider you can order servers by size and city, keep DNS zones, records, GeoDNS answers and DNSSEC in code, put firewalls in front of servers, and manage SSH keys, roles, groups and secrets.

## Installation

{{< chooser language "typescript,python,go,csharp,java,yaml,hcl" >}}
{{% choosable language typescript %}}

npm:

```bash
npm install @noderings/pulumi
```

Yarn:

```bash
yarn add @noderings/pulumi
```

pnpm:

```bash
pnpm add @noderings/pulumi
```

Bun:

```bash
bun add @noderings/pulumi
```

{{% /choosable %}}
{{% choosable language python %}}

pip:

Add `pulumi-noderings` to `requirements.txt`, then run:

```bash
pip install -r requirements.txt
```

uv:

```bash
uv add pulumi-noderings
```

Poetry:

```bash
poetry add pulumi-noderings
```

{{% /choosable %}}
{{% choosable language go %}}

```bash
go get github.com/noderings/pulumi-noderings/sdk/go/noderings
```

{{% /choosable %}}
{{% choosable language csharp %}}

```bash
dotnet add package NodeRings.Pulumi
```

{{% /choosable %}}
{{% choosable language java %}}

There is no published Java SDK. Generate one into your project from the provider:

```bash
pulumi package add noderings
```

{{% /choosable %}}
{{% choosable language yaml %}}

```bash
pulumi package add noderings
```

{{% /choosable %}}
{{% choosable language hcl %}}

Declare the provider in your program, then run `pulumi install`:

<!-- x-release-please-start-version -->
```hcl
terraform {
  required_providers {
    noderings = {
      source  = "pulumi/noderings"
      version = "0.1.6"
    }
  }
}
```
<!-- x-release-please-end -->

```bash
pulumi install
```

{{% /choosable %}}
{{< /chooser >}}

The provider plugin is downloaded from [GitHub releases](https://github.com/noderings/pulumi-noderings/releases) the first time a program uses it.

## Example Usage

This program creates a DNS zone, which is free, and prints the nameservers to set at your registrar.

{{< chooser language "typescript,python,go,csharp,yaml" >}}
{{% choosable language typescript %}}

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as noderings from "@noderings/pulumi";

const zone = new noderings.DnsZone("example", { name: "example.com" });

export const nameservers = zone.nameservers;
```

{{% /choosable %}}
{{% choosable language python %}}

```python
import pulumi
import pulumi_noderings as noderings

zone = noderings.DnsZone("example", name="example.com")

pulumi.export("nameservers", zone.nameservers)
```

{{% /choosable %}}
{{% choosable language go %}}

```go
package main

import (
	"github.com/noderings/pulumi-noderings/sdk/go/noderings"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		zone, err := noderings.NewDnsZone(ctx, "example", &noderings.DnsZoneArgs{
			Name: pulumi.String("example.com"),
		})
		if err != nil {
			return err
		}
		ctx.Export("nameservers", zone.Nameservers)
		return nil
	})
}
```

{{% /choosable %}}
{{% choosable language csharp %}}

```csharp
using System.Collections.Generic;
using Pulumi;
using NodeRings.Pulumi;

return await Deployment.RunAsync(() =>
{
    var zone = new DnsZone("example", new()
    {
        Name = "example.com",
    });

    return new Dictionary<string, object?>
    {
        ["nameservers"] = zone.Nameservers,
    };
});
```

{{% /choosable %}}
{{% choosable language yaml %}}

```yaml
name: noderings-example
runtime: yaml
resources:
  example:
    type: noderings:DnsZone
    properties:
      name: example.com
outputs:
  nameservers: ${example.nameservers}
```

{{% /choosable %}}
{{< /chooser >}}

Set the token before running `pulumi up`:

```bash
pulumi config set --secret noderings:token <your-service-account-token>
```

## Configuration

Create the token in the [NodeRings console](https://noderings.com): open **Organization**, then **Access Control**, then **Service accounts**, create a service account, give it a role, and generate a token. The console shows the token once. Administrator covers everything; a narrower role is enough if a program only manages DNS, for example.

- `token` (Required, Secret) — Service account token. May also be set with the `NODERINGS_TOKEN` environment variable (`NODERINGS_API_TOKEN` and `NR_API_TOKEN` also work).
- `organizationId` (Optional) — Organization to act in, sent as the `X-Organization-ID` header. Defaults to the token's own organization. May also be set with the `NODERINGS_ORGANIZATION_ID` environment variable.
- `endpoint` (Optional) — API base URL. Defaults to the public NodeRings API. May also be set with the `NODERINGS_API_URL` environment variable.
- `maxRetries` (Optional) — How many times to retry rate-limited and transient failures. Defaults to `25`.
- `requestsPerSecond` (Optional) — Client-side limit on requests per second for each API operation. Defaults to `2`, the API's own limit. Set `0` to disable it.
- `httpTimeout` (Optional) — Timeout for a single HTTP request, as a duration such as `60s`. Defaults to `60s`.
- `skipOrganizationCheck` (Optional) — Skip checking the token and organization when the provider starts. Defaults to `false`.

To work in two organizations from one program, create an explicit `noderings.Provider` for each, with its own `token` or `organizationId`, and pass it to resources with the `provider` resource option.

## Order a virtual machine

A plan code such as `SC-GC-1C-1G-20` (shared compute, 1 vCPU, 1 GB RAM, 20 GB disk) is a size that several hosting providers may sell, in several cities, at different prices. `getMarketplaceOffer` returns the cheapest matching offer, or the best rated with `sort: "rating"`, and you can pin one provider with `providerId` or `storefrontSlug`. `billingType: "usage"` asks for hourly billing from prepaid credit, which needs at least $10 available to order.

This program orders the cheapest hourly offer in Ashburn with your SSH key and points `www` at it:

{{< chooser language "typescript,python,go,csharp,yaml" >}}
{{% choosable language typescript %}}

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as noderings from "@noderings/pulumi";
import * as fs from "fs";
import * as os from "os";

const offer = noderings.getMarketplaceOfferOutput({
    planCode: "SC-GC-1C-1G-20",
    country: "United States",
    city: "Ashburn",
    billingType: "usage",
});

const key = new noderings.SshKey("laptop", {
    name: "laptop",
    publicKey: fs.readFileSync(`${os.homedir()}/.ssh/id_ed25519.pub`, "utf8"),
});

const web = new noderings.VirtualMachine("web", {
    hostname: "web-1",
    planId: offer.planId,
    regionId: offer.regionId,
    distro: "ubuntu_24_04",
    sshKeyIds: [key.id],
    tags: ["web"],
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

{{% /choosable %}}
{{% choosable language python %}}

```python
import os
import pulumi
import pulumi_noderings as noderings

offer = noderings.get_marketplace_offer_output(
    plan_code="SC-GC-1C-1G-20",
    country="United States",
    city="Ashburn",
    billing_type="usage",
)

with open(os.path.expanduser("~/.ssh/id_ed25519.pub")) as f:
    key = noderings.SshKey("laptop", name="laptop", public_key=f.read())

web = noderings.VirtualMachine(
    "web",
    hostname="web-1",
    plan_id=offer.plan_id,
    region_id=offer.region_id,
    distro="ubuntu_24_04",
    ssh_key_ids=[key.id],
    tags=["web"],
)

zone = noderings.DnsZone("site", name="example.com")

noderings.DnsRecord("www", zone_id=zone.id, name="www", type="A", ttl=300, records=[web.ipv4])

pulumi.export("ip", web.ipv4)
pulumi.export("nameservers", zone.nameservers)
```

{{% /choosable %}}
{{% choosable language go %}}

```go
package main

import (
	"os"
	"path/filepath"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/noderings/pulumi-noderings/sdk/go/noderings"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		offer := noderings.GetMarketplaceOfferOutput(ctx, noderings.GetMarketplaceOfferOutputArgs{
			PlanCode:    pulumi.String("SC-GC-1C-1G-20"),
			Country:     pulumi.String("United States"),
			City:        pulumi.String("Ashburn"),
			BillingType: pulumi.String("usage"),
		})

		home, _ := os.UserHomeDir()
		pub, err := os.ReadFile(filepath.Join(home, ".ssh", "id_ed25519.pub"))
		if err != nil {
			return err
		}
		key, err := noderings.NewSshKey(ctx, "laptop", &noderings.SshKeyArgs{
			Name:      pulumi.String("laptop"),
			PublicKey: pulumi.String(string(pub)),
		})
		if err != nil {
			return err
		}

		web, err := noderings.NewVirtualMachine(ctx, "web", &noderings.VirtualMachineArgs{
			Hostname:  pulumi.String("web-1"),
			PlanId:    offer.PlanId(),
			RegionId:  offer.RegionId(),
			Distro:    pulumi.String("ubuntu_24_04"),
			SshKeyIds: pulumi.StringArray{key.ID()},
			Tags:      pulumi.StringArray{pulumi.String("web")},
		})
		if err != nil {
			return err
		}

		zone, err := noderings.NewDnsZone(ctx, "site", &noderings.DnsZoneArgs{
			Name: pulumi.String("example.com"),
		})
		if err != nil {
			return err
		}

		_, err = noderings.NewDnsRecord(ctx, "www", &noderings.DnsRecordArgs{
			ZoneId:  zone.ID(),
			Name:    pulumi.String("www"),
			Type:    pulumi.String("A"),
			Ttl:     pulumi.Int(300),
			Records: pulumi.StringArray{web.Ipv4},
		})
		if err != nil {
			return err
		}

		ctx.Export("ip", web.Ipv4)
		ctx.Export("nameservers", zone.Nameservers)
		return nil
	})
}
```

{{% /choosable %}}
{{% choosable language csharp %}}

```csharp
using System;
using System.Collections.Generic;
using System.IO;
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
        PublicKey = File.ReadAllText(Path.Combine(
            Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), ".ssh", "id_ed25519.pub")),
    });

    var web = new VirtualMachine("web", new()
    {
        Hostname = "web-1",
        PlanId = offer.Apply(o => o.PlanId),
        RegionId = offer.Apply(o => o.RegionId),
        Distro = "ubuntu_24_04",
        SshKeyIds = new[] { key.Id },
        Tags = new[] { "web" },
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

{{% /choosable %}}
{{% choosable language yaml %}}

```yaml
name: noderings-web
runtime: yaml
config:
  sshPublicKey:
    type: string
variables:
  offer:
    fn::invoke:
      function: noderings:getMarketplaceOffer
      arguments:
        planCode: SC-GC-1C-1G-20
        country: United States
        city: Ashburn
        billingType: usage
resources:
  laptop:
    type: noderings:SshKey
    properties:
      name: laptop
      publicKey: ${sshPublicKey}
  web:
    type: noderings:VirtualMachine
    properties:
      hostname: web-1
      planId: ${offer.planId}
      regionId: ${offer.regionId}
      distro: ubuntu_24_04
      sshKeyIds:
        - ${laptop.id}
      tags:
        - web
  site:
    type: noderings:DnsZone
    properties:
      name: example.com
  www:
    type: noderings:DnsRecord
    properties:
      zoneId: ${site.id}
      name: www
      type: A
      ttl: 300
      records:
        - ${web.ipv4}
outputs:
  ip: ${web.ipv4}
  nameservers: ${site.nameservers}
```

{{% /choosable %}}
{{< /chooser >}}

The YAML program reads the public key from config: `pulumi config set sshPublicKey "$(cat ~/.ssh/id_ed25519.pub)"`.

The offer lookup runs on every `pulumi up`. If a cheaper offer appears later, `planId` changes and Pulumi plans to replace the server. Add `ignoreChanges: ["planId", "regionId"]` to the server's resource options, or pin the provider, to keep it where it is.

## What changes in place

| Change | Effect |
|---|---|
| `tags`, `powerState`, firewall rules, address groups | Updated in place |
| `hostname`, `distro`, `sshKeyIds` | New server, unless `reinstallOnChange: true` reinstalls the same one |
| `planId`, `regionId` | New server. There is no resize |

A replacement deletes the old server, and with it its IP addresses and disk. Set `deletionProtection: true` on anything billed monthly: deleting a monthly subscription cancels it with no refund for the rest of the period. Usage-billed servers stop billing when they are deleted.

## Resources

| Area | Resources |
|---|---|
| Servers | `VirtualMachine`, `VirtualMachineBackup`, `VirtualMachineSnapshot`, `VirtualMachineAdditionalIp` |
| DNS | `DnsZone`, `DnsRecord`, `DnsGeoRecord`, `DnsZoneDnssec` |
| Firewall | `FirewallRuleset`, `FirewallRulesetRules`, `FirewallRulesetAttachment`, `FirewallAddressGroup` |
| Access | `SshKey`, `Role`, `RoleAssignment`, `Group`, `GroupMembership`, `Secret`, `EnvironmentVariable` |

Functions such as `getMarketplaceOffer`, `getMarketplaceOffers`, `getDistros`, `getVirtualMachines` and `getCreditBalance` read the catalog and your account. Every resource can be adopted with `pulumi import`, for example `pulumi import noderings:index/dnsRecord:DnsRecord www <zone id>/www/A`.

## Rate limits

The NodeRings API allows about two requests per second for each operation. The provider paces itself and retries rate-limited calls, so a large stack takes longer and still finishes.
