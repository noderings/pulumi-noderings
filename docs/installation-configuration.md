---
title: NodeRings Installation & Configuration
meta_desc: How to install the Pulumi NodeRings provider and configure it with a NodeRings service account token.
layout: package
---

## Installation

The NodeRings provider is available as a package in these Pulumi languages:

* JavaScript/TypeScript: [`@noderings/pulumi`](https://www.npmjs.com/package/@noderings/pulumi)
* Python: [`pulumi-noderings`](https://pypi.org/project/pulumi-noderings/)
* Go: [`github.com/noderings/pulumi-noderings/sdk/go/noderings`](https://pkg.go.dev/github.com/noderings/pulumi-noderings/sdk/go/noderings)
* .NET: [`NodeRings.Pulumi`](https://www.nuget.org/packages/NodeRings.Pulumi)
* YAML and Java: generate a local SDK with `pulumi package add noderings`

{{< chooser language "typescript,python,go,csharp,java,yaml" >}}
{{% choosable language typescript %}}

```bash
npm install @noderings/pulumi
```

{{% /choosable %}}
{{% choosable language python %}}

```bash
pip install pulumi-noderings
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

There is no published Java SDK. Generate one into your project:

```bash
pulumi package add noderings
```

{{% /choosable %}}
{{% choosable language yaml %}}

```bash
pulumi package add noderings
```

{{% /choosable %}}
{{< /chooser >}}

The provider plugin is downloaded from the [GitHub releases](https://github.com/noderings/pulumi-noderings/releases) of `noderings/pulumi-noderings` the first time a program uses it. To install it ahead of time, for example on a build agent without internet access at run time:

<!-- x-release-please-start-version -->
```bash
pulumi plugin install resource noderings 0.1.3 --server github://api.github.com/noderings
```
<!-- x-release-please-end -->

## Create a service account token

The provider authenticates with a NodeRings service account token:

1. Sign in to the [NodeRings console](https://noderings.com).
2. Open **Organization**, then **Access Control**, then **Service accounts**.
3. Create a service account and give it a role. Administrator covers everything; a narrower role is enough if a program only manages DNS, for example.
4. Generate a token. The console shows it once.

## Configuration

Set the token as a secret in your stack:

```bash
pulumi config set --secret noderings:token <your-service-account-token>
```

Or use an environment variable, which the provider reads when `token` is not set:

```bash
export NODERINGS_TOKEN=<your-service-account-token>
```

| Name | Required | Secret | Description |
|---|---|---|---|
| `token` | Yes | Yes | Service account token. Environment variable: `NODERINGS_TOKEN` (`NODERINGS_API_TOKEN` and `NR_API_TOKEN` also work). |
| `organizationId` | No | No | Organization to act in, sent as the `X-Organization-ID` header. Defaults to the token's own organization. Environment variable: `NODERINGS_ORGANIZATION_ID`. |
| `endpoint` | No | No | API base URL. Defaults to `https://api.noderings.com`. Environment variable: `NODERINGS_API_URL`. |
| `maxRetries` | No | No | How many times to retry rate-limited and transient failures. Defaults to `25`. |
| `requestsPerSecond` | No | No | Client-side limit on requests per second for each API operation. Defaults to `2`, the API's own limit. Set `0` to disable it. |
| `httpTimeout` | No | No | Timeout for a single HTTP request, as a duration such as `60s`. Defaults to `60s`. |
| `skipOrganizationCheck` | No | No | Skip checking the token and organization when the provider starts. Defaults to `false`. |

### Several organizations in one program

Create an explicit provider for each organization, each with its own `token` or `organizationId`, and pass it to resources with the `provider` resource option.

{{< chooser language "typescript,python,go,csharp,yaml" >}}
{{% choosable language typescript %}}

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as noderings from "@noderings/pulumi";

const config = new pulumi.Config();
const staging = new noderings.Provider("staging", {
    token: config.requireSecret("stagingToken"),
});

new noderings.DnsZone("example", { name: "staging.example.com" }, { provider: staging });
```

{{% /choosable %}}
{{% choosable language python %}}

```python
import pulumi
import pulumi_noderings as noderings

config = pulumi.Config()
staging = noderings.Provider("staging", token=config.require_secret("stagingToken"))

noderings.DnsZone("example", name="staging.example.com",
                  opts=pulumi.ResourceOptions(provider=staging))
```

{{% /choosable %}}
{{% choosable language go %}}

```go
package main

import (
	"github.com/noderings/pulumi-noderings/sdk/go/noderings"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		staging, err := noderings.NewProvider(ctx, "staging", &noderings.ProviderArgs{
			Token: cfg.RequireSecret("stagingToken"),
		})
		if err != nil {
			return err
		}
		_, err = noderings.NewDnsZone(ctx, "example", &noderings.DnsZoneArgs{
			Name: pulumi.String("staging.example.com"),
		}, pulumi.Provider(staging))
		return err
	})
}
```

{{% /choosable %}}
{{% choosable language csharp %}}

```csharp
using Pulumi;
using NodeRings.Pulumi;

return await Deployment.RunAsync(() =>
{
    var config = new Config();
    var staging = new Provider("staging", new ProviderArgs
    {
        Token = config.RequireSecret("stagingToken"),
    });

    new DnsZone("example", new DnsZoneArgs
    {
        Name = "staging.example.com",
    }, new CustomResourceOptions { Provider = staging });
});
```

{{% /choosable %}}
{{% choosable language yaml %}}

```yaml
name: noderings-two-orgs
runtime: yaml
config:
  stagingToken:
    type: string
    secret: true
resources:
  staging:
    type: pulumi:providers:noderings
    properties:
      token: ${stagingToken}
  example:
    type: noderings:DnsZone
    properties:
      name: staging.example.com
    options:
      provider: ${staging}
```

{{% /choosable %}}
{{< /chooser >}}
