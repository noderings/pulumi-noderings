# Examples

Each folder is a complete Pulumi program. They double as the integration tests in `examples_*_test.go`, which deploy them against the live NodeRings API, check the outputs, confirm a second preview shows no changes, and destroy them.

| Folder | Language | What it creates | Cost |
|---|---|---|---|
| `dns-ts` | TypeScript | Zone with A, MX, TXT, CAA and GeoDNS records, plus DNSSEC | Free |
| `dns-go` | Go | Zone with an A and an MX record | Free |
| `vm-py` | Python | Hourly VM, firewall, and a DNS record pointing at it | A few cents an hour |

## Run one yourself

```bash
export NODERINGS_TOKEN=<service account token>
cd dns-ts
npm install && npm install @noderings/pulumi     # or: pip install pulumi-noderings / go get ...
pulumi stack init dev
pulumi config set zone example.com
pulumi up
pulumi destroy
```

The test manifests leave out the NodeRings SDK because the test harness links the local build. Add it when you copy an example into your own project: `@noderings/pulumi` for TypeScript, `pulumi-noderings` for Python, `github.com/noderings/pulumi-noderings/sdk` for Go.

`vm-py` orders a real server billed by the hour. It also needs `pulumi config set sshPublicKey "$(cat ~/.ssh/id_ed25519.pub)"`, and the server stops billing when you run `pulumi destroy`.
