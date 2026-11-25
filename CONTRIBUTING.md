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

Go (see `provider/go.mod`), the [Pulumi CLI](https://www.pulumi.com/docs/install/) at the version in `.pulumi-version` (code generation output depends on it), Node.js 22 with Yarn 1, Python 3.9+ with `pip install build`, and the .NET 8 SDK.

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

Releases are automated; nobody pushes tags or builds by hand. Use [Conventional Commits](https://www.conventionalcommits.org/) for commit messages and pull request titles (pull requests are squash-merged, so the title becomes the commit):

| Prefix | Changelog section | Version bump (before 1.0) |
|---|---|---|
| `feat:` | Features | minor |
| `fix:` | Bug Fixes | patch |
| `perf:` | Performance | patch |
| `deps:` | Dependencies | patch |
| `docs:` | Documentation | patch |
| `chore:`, `ci:`, `test:`, `refactor:`, `build:` | not listed | none |

1. [release-please](https://github.com/googleapis/release-please) keeps one open pull request, `chore: release X.Y.Z`, that updates `CHANGELOG.md`, `.release-please-manifest.json` and the version shown in `docs/`.
2. Merging it pushes the `vX.Y.Z` tag and a draft GitHub release with the changelog.
3. The tag starts `.github/workflows/release.yml`:
   - GoReleaser builds the plugin for every platform, uploads it and `schema.json` to the draft, and publishes the release, with a build provenance attestation.
   - The Go SDK is tagged `sdk/vX.Y.Z`.
   - The npm, PyPI and NuGet SDKs are built with that version and published, each with provenance.
4. The Pulumi Registry picks the release up on its twice-daily check and opens its own pull request with the new docs.

A registry that is not set up is skipped with a warning, so the rest still publish. To republish a tag, run **Release** from the Actions tab and pick the tag under "Use workflow from"; versions already on a registry are skipped.

### One-time registry setup

The publish jobs run in the GitHub environments `npm`, `pypi` and `nuget`, which only accept `v*` tags. Each registry trusts `release.yml` in this repository, so no long-lived tokens are needed after the first npm publish.

| Registry | Setup | Then set (Settings → Secrets and variables → Actions) |
|---|---|---|
| PyPI | Account → Publishing → Add a new pending publisher: project `pulumi-noderings`, owner `noderings`, repository `pulumi-noderings`, workflow `release.yml`, environment `pypi` | Variable `PYPI_TRUSTED_PUBLISHING` = `true` |
| NuGet | Account → Trusted Publishing → Create: repository owner `noderings`, repository `pulumi-noderings`, workflow file `release.yml`, environment `nuget` | Variable `NUGET_USER` = the nuget.org user name that owns the policy. (Or secret `NUGET_API_KEY` instead.) |
| npm | Organization `noderings`. npm can only add a trusted publisher to a package that already exists, so the first version is published with a short-lived granular token. Then on the package's Settings → Trusted publishing → GitHub Actions: organization `noderings`, repository `pulumi-noderings`, workflow `release.yml`, environment `npm`, allow `npm publish` | First release: secret `NPM_TOKEN`. Afterwards: delete `NPM_TOKEN` and the token on npmjs.com, set variable `NPM_TRUSTED_PUBLISHING` = `true` |

`NPM_TOKEN` wins over trusted publishing while it is set. A new npm trusted publisher must complete its first publish within two days or it expires, so add it right before a release.

## Code of conduct

See [CODE-OF-CONDUCT.md](CODE-OF-CONDUCT.md).
