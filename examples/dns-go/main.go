// Command noderings-dns-go creates a DNS zone with an A and an MX record.
package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/noderings/pulumi-noderings/sdk/go/noderings"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")

		zone, err := noderings.NewDnsZone(ctx, "site", &noderings.DnsZoneArgs{
			Name: pulumi.String(cfg.Require("zone")),
		})
		if err != nil {
			return err
		}

		_, err = noderings.NewDnsRecord(ctx, "apex", &noderings.DnsRecordArgs{
			ZoneId:  zone.ID(),
			Name:    pulumi.String("@"),
			Type:    pulumi.String("A"),
			Ttl:     pulumi.Int(300),
			Records: pulumi.StringArray{pulumi.String("192.0.2.10")},
		})
		if err != nil {
			return err
		}

		_, err = noderings.NewDnsRecord(ctx, "mx", &noderings.DnsRecordArgs{
			ZoneId: zone.ID(),
			Name:   pulumi.String("@"),
			Type:   pulumi.String("MX"),
			Ttl:    pulumi.Int(3600),
			Records: pulumi.StringArray{
				pulumi.String("10 mx1.mail.example.net."),
				pulumi.String("20 mx2.mail.example.net."),
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("nameservers", zone.Nameservers)
		return nil
	})
}
