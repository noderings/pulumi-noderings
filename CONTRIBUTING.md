# Contributing

Thanks for helping. Bug reports, fixes and new examples are all welcome.

## How this repository works

This provider is generated from the [NodeRings Terraform provider](https://github.com/noderings/terraform-provider-noderings) with the [Pulumi Terraform bridge](https://github.com/pulumi/pulumi-terraform-bridge). Resource behavior (API calls, validation, retries, import) lives in the Terraform provider. A bug in what a resource *does* usually belongs there; a bug in names, SDK packaging or docs belongs here.

| Path | What it is |
|---|---|
| `provider/resources.go` | The mapping from Terraform to Pulumi: names, config, SDK package settings |
| `provider/cmd/pulumi-tfgen-noderings` | Generates the Pulumi schema |
| `provider/cmd/pulumi-resource-noderings` | The provider plugin Pulumi runs |
| `provider/cmd/pulumi-resource-noderings/schema.json` | The generated schema (committed) |
| `sdk/` | Generated SDKs for Node.js, Python, Go and .NET (committed) |
| `examples/` | Example programs, also used as integration tests |

## Prerequisites

Go (see `provider/go.mod`), the [Pulumi CLI](https://www.pulumi.com/docs/install/), Node.js 22 with Yarn 1, Python 3.9+ with `pip install build`, and the .NET 8 SDK.

## Making a change

1. Change `provider/resources.go`, or bump the Terraform provider in `provider/go.mod`.
2. Regenerate the schema and SDKs and commit the result:

   ```bash
   make tfgen
   make generate_sdks
   ```

3. Run the unit tests (no account needed):

   ```bash
   make test_provider
   ```

4. Open a pull request. CI checks that the committed schema and SDKs match what the generator produces.

## Integration tests

The programs in `examples/` are deployed against the real NodeRings API, checked, and destroyed. Each test also checks that a second `pulumi preview` shows no changes.

```bash
export NODERINGS_TOKEN=...            # a service account token
make test                             # DNS tests only, no charge
NODERINGS_TEST_BILLED=1 make test     # also orders an hourly VM for a few minutes
```

| Variable | Purpose |
|---|---|
| `NODERINGS_TOKEN` | Required. Tests are skipped without it. |
| `NODERINGS_TEST_BILLED=1` | Also run tests that order a billed VM (a few cents per run). |
| `NODERINGS_TEST_SSH_PUBLIC_KEY` | Public key for the test VM. Defaults to `~/.ssh/id_ed25519.pub`. |
| `NODERINGS_TEST_DOMAIN` | Parent domain for the throwaway test zones. Defaults to `example.com`. |

Test zones and servers get random names, and every test destroys what it created, also when it fails.

## Releasing

Push a tag such as `v0.2.0`. The release workflow builds the plugin for every platform, attaches it to a GitHub release, and publishes the SDKs to npm, PyPI and NuGet. The Go SDK is released by tagging `sdk/v0.2.0` on the same commit, which the workflow does too.

## Code of conduct

See [CODE-OF-CONDUCT.md](CODE-OF-CONDUCT.md).
