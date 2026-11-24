"""Order an hourly VM, put a firewall in front of it, and point a DNS name at it."""

import pulumi
import pulumi_noderings as noderings

config = pulumi.Config()
zone_name = config.require("zone")
public_key = config.require("sshPublicKey")
stack = pulumi.get_stack()

# Cheapest hourly offer for 1 vCPU, 1 GB RAM and 20 GB disk in Ashburn whose
# provider supports the guest firewall.
offer = noderings.get_marketplace_offer(
    plan_code="SC-GC-1C-1G-20",
    country="United States",
    city="Ashburn",
    billing_type="usage",
    require_capabilities=["guest_firewall"],
)

key = noderings.SshKey("laptop", name=f"pulumi-{stack}", public_key=public_key)

vm = noderings.VirtualMachine(
    "web",
    hostname=f"pulumi-{stack}",
    plan_id=offer.plan_id,
    region_id=offer.region_id,
    distro="ubuntu_24_04",
    ssh_key_ids=[key.id],
    tags=["pulumi-test"],
)

ruleset = noderings.FirewallRuleset("web", name=f"pulumi-{stack}")

noderings.FirewallRulesetRules(
    "web",
    ruleset_id=ruleset.id,
    rules=[
        noderings.FirewallRulesetRulesRuleArgs(
            direction="ingress", action="accept", protocol="tcp", destination_ports="22",
            source=noderings.FirewallRulesetRulesRuleSourceArgs(cidrs=["0.0.0.0/0"]),
        ),
        noderings.FirewallRulesetRulesRuleArgs(
            direction="ingress", action="accept", protocol="tcp", destination_ports="443",
            source=noderings.FirewallRulesetRulesRuleSourceArgs(cidrs=["0.0.0.0/0"]),
        ),
    ],
)

noderings.FirewallRulesetAttachment("web", virtual_machine_id=vm.id, ruleset_id=ruleset.id)

zone = noderings.DnsZone("site", name=zone_name)
noderings.DnsRecord("www", zone_id=zone.id, name="www", type="A", ttl=300, records=[vm.ipv4])

pulumi.export("provider", offer.provider_name)
pulumi.export("ipv4", vm.ipv4)
pulumi.export("ssh", vm.ipv4.apply(lambda ip: f"ssh root@{ip}"))
