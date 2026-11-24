---
title: NodeRings
meta_desc: Provides an overview of the NodeRings provider for Pulumi.
layout: package
---

The NodeRings provider for Pulumi manages resources on [NodeRings](https://noderings.com), a cloud marketplace where independent hosting providers sell virtual machines and game servers through one account and one API.

With it you can order a virtual machine by size and city and get its address back, keep DNS zones, records, GeoDNS answers and DNSSEC in code, put firewalls in front of servers, run game servers with their backups and databases, and manage service accounts, roles and groups.

The provider needs a NodeRings service account token. See [Installation & Configuration](./installation-configuration/).

## Example

This program finds the cheapest hourly offer for a 1 vCPU, 1 GB server in Ashburn, orders it with your SSH key, and points `www` at it.

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

## How ordering works

A plan code such as `SC-GC-1C-1G-20` (shared compute, 1 vCPU, 1 GB RAM, 20 GB disk) is a size that several hosting providers may sell, in several cities, at different prices. `getMarketplaceOffer` returns the cheapest matching offer, or the best rated with `sort: "rating"`, and you can pin one provider with `providerId` or `storefrontSlug`. `billingType: "usage"` asks for hourly billing from prepaid credit. Without it you can get a monthly plan charged to your card.

The offer lookup runs on every `pulumi up`. If a cheaper offer appears later, `planId` changes and Pulumi plans to replace the server. Add `ignoreChanges: ["planId", "regionId"]` to the server's resource options, or pin the provider, to keep it where it is.

## What changes in place

| Change | Effect |
|---|---|
| `tags`, `powerState`, firewall rules, address groups | Updated in place |
| `hostname`, `distro`, `sshKeyIds` | New server, unless `reinstallOnChange: true` reinstalls the same one |
| `planId`, `regionId` | New server. There is no resize |

A replacement deletes the old server, and with it its IP addresses and disk. Set `deletionProtection: true` on anything you pay for monthly: deleting a monthly subscription cancels it with no refund for the rest of the period.

## Resources

| Area | Resources |
|---|---|
| Servers | `VirtualMachine`, `VirtualMachineBackup`, `VirtualMachineSnapshot`, `VirtualMachineAdditionalIp` |
| DNS | `DnsZone`, `DnsRecord`, `DnsGeoRecord`, `DnsZoneDnssec` |
| Firewall | `FirewallRuleset`, `FirewallRulesetRules`, `FirewallRulesetAttachment`, `FirewallAddressGroup` |
| Game servers | `GameServer`, `GameServerAllocation`, `GameServerBackup`, `GameServerDatabase`, `GameServerSchedule` |
| Access | `SshKey`, `Role`, `RoleAssignment`, `Group`, `GroupMembership`, `Secret`, `EnvironmentVariable` |

Functions such as `getMarketplaceOffer`, `getMarketplaceOffers`, `getDistros`, `getVirtualMachines` and `getCreditBalance` read the catalog and your account. The full list is in the [API documentation](./api-docs/).
