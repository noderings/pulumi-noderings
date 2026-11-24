import * as pulumi from "@pulumi/pulumi";
import * as noderings from "@noderings/pulumi";

const config = new pulumi.Config();
const zoneName = config.require("zone");

const zone = new noderings.DnsZone("site", { name: zoneName });

const apex = new noderings.DnsRecord("apex", {
    zoneId: zone.id,
    name: "@",
    type: "A",
    ttl: 300,
    records: ["192.0.2.10"],
});

new noderings.DnsRecord("mx", {
    zoneId: zone.id,
    name: "@",
    type: "MX",
    ttl: 3600,
    records: ["10 mx1.mail.example.net.", "20 mx2.mail.example.net."],
});

new noderings.DnsRecord("spf", {
    zoneId: zone.id,
    name: "@",
    type: "TXT",
    ttl: 3600,
    records: ["v=spf1 include:mail.example.net -all"],
});

new noderings.DnsRecord("caa", {
    zoneId: zone.id,
    name: "@",
    type: "CAA",
    ttl: 3600,
    records: ['0 issue "letsencrypt.org"'],
});

new noderings.DnsGeoRecord("app", {
    zoneId: zone.id,
    name: "app",
    type: "A",
    ttl: 60,
    default: "192.0.2.10",
    continents: { EU: "198.51.100.10" },
    countries: { JP: "203.0.113.10" },
});

const dnssec = new noderings.DnsZoneDnssec("site", { zoneId: zone.id });

export const nameservers = zone.nameservers;
export const zoneState = zone.state;
export const dsRecord = dnssec.ds;
export const apexRecords = apex.records;
