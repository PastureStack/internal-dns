# Compatibility

The current POC preserves A, AAAA, CNAME, PTR, TXT, recursive lookup, search suffix, authoritative suffix, EDNS0, TCP fallback, response truncation, TTL, cache, metadata-region, service-link, sidekick, external-service, and Kubernetes service behaviors covered by the historical test suite.

The public internal domain changes to `pasture.internal`, and the metadata record changes to `metadata.pasture.internal`. This is a deliberate breaking change. A production migration needs a measured dual-resolution or coordinated cutover plan in the Catalog, metadata producer, containers, host resolver configuration, and every consumer.

Metadata-driven mode depends on these `2016-07-29` resources:

- `/version` with long polling;
- `/services` and `/containers`;
- `/self/host`;
- `/region_name` and `/environments`;
- `/stacks/{stack}/services/{service}`.

Answer files emit the canonical `authoritative` key. The historical misspelling `authorative` remains accepted on input so existing operator files can be upgraded without an abrupt parsing failure.

The Windows variant targets the Windows Server 2022 (`ltsc2022`) container ABI.
Its Nano Server startup command assigns `169.254.169.251/32` through `netsh`
and then starts the same cross-compiled DNS service. The container therefore
requires Windows container administrator privileges for link-local address
configuration. Linux and Windows link-local addresses, recursive resolver
selection, DNS port binding, and rollback must be verified in disposable hosts
before production release.
