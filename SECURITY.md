# Security policy

## Reporting a vulnerability

Please do not open a public issue for security problems. Email **partners@noderings.com** with:

- a description of the issue and its impact
- steps to reproduce, or a proof of concept
- the package and SDK versions, and the Pulumi CLI version (`pulumi version`)

We acknowledge reports within three working days and keep you updated until the issue is fixed. Credit is given in the release notes unless you ask us not to.

## Supported versions

Security fixes are made for the latest minor release.

## Handling credentials

The provider reads the NodeRings service account token from the `noderings:token` stack configuration or the `NODERINGS_TOKEN` environment variable. Set it as a secret with `pulumi config set noderings:token --secret`, so it is encrypted in the stack configuration and state, and give the service account the narrowest NodeRings role that covers the resources you manage.
