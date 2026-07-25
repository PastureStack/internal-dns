# Security Policy

This migration candidate has no supported production release.

Report suspected vulnerabilities privately to the PastureStack organization maintainers. Do not publish production DNS data, resolver addresses, metadata responses, customer domains, certificates, credentials, or internal network topology.

Integration tests must use disposable VMs and non-production recursive resolvers. Public disclosure should wait until maintainers confirm a coordinated remediation plan.

The Windows Server 2022 image starts as `ContainerAdministrator` only because
Windows link-local address assignment requires that privilege. It does not
mount the Docker socket or host filesystem. Treat recursive resolver addresses
and metadata responses as environment topology and keep them out of public
logs.
