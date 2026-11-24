---
title: NodeRings Installation & Configuration
meta_desc: Information on how to install the NodeRings provider for Pulumi and configure its credentials.
layout: package
---

## Installation

The NodeRings provider is available as a package in these Pulumi languages:

* JavaScript/TypeScript: [`@noderings/pulumi`](https://www.npmjs.com/package/@noderings/pulumi)
* Python: [`pulumi-noderings`](https://pypi.org/project/pulumi-noderings/)
* Go: [`github.com/noderings/pulumi-noderings/sdk/go/noderings`](https://pkg.go.dev/github.com/noderings/pulumi-noderings/sdk/go/noderings)
* .NET: [`NodeRings.Pulumi`](https://www.nuget.org/packages/NodeRings.Pulumi)

```bash
npm install @noderings/pulumi
pip install pulumi-noderings
go get github.com/noderings/pulumi-noderings/sdk/go/noderings
dotnet add package NodeRings.Pulumi
```

The provider plugin is downloaded automatically from the [GitHub releases](https://github.com/noderings/pulumi-noderings/releases) the first time a program uses it. To install it by hand:

```bash
pulumi plugin install resource noderings --server github://api.github.com/noderings
```

## Credentials

The provider signs in with a NodeRings service account token.

1. Sign in to the [NodeRings console](https://noderings.com).
2. Open **Organization**, then **Access Control**, then **Service accounts**.
3. Create a service account and give it a role. Administrator covers everything; a narrower role is enough if the program only manages DNS, for example.
4. Generate a token and copy it. The console shows it once.

Ordering usage-billed servers draws from prepaid credit, and an order needs at least $10 available.

## Configuration

Set the token as an environment variable:

```bash
export NODERINGS_TOKEN=<your token>
```

or store it encrypted in the stack configuration:

```bash
pulumi config set noderings:token --secret
```

| Option | Environment variable | Description |
|---|---|---|
| `noderings:token` | `NODERINGS_TOKEN` | Service account token. Required. |
| `noderings:organizationId` | `NODERINGS_ORGANIZATION_ID` | Organization to act in. Defaults to the token's own organization. |
| `noderings:endpoint` | `NODERINGS_API_URL` | API base URL. Defaults to the public NodeRings API. |
| `noderings:maxRetries` | | Retries for rate-limited and transient failures. Defaults to `25`. |
| `noderings:requestsPerSecond` | | Client-side limit per API operation. Defaults to `2`, the API's own limit. `0` disables it. |
| `noderings:httpTimeout` | | Timeout for one HTTP request, such as `60s`. Defaults to `60s`. |
| `noderings:skipOrganizationCheck` | | Skip checking the token and organization at startup. Defaults to `false`. |

To manage two organizations from one program, create an explicit provider for each and pass it with the `provider` resource option.

## Rate limits

The NodeRings API allows about two requests per second for each operation. The provider paces itself and retries rate-limited calls, so a large stack takes longer and still finishes.
