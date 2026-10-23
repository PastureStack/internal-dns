# Internal DNS

Current source builds use Go 1.27.3 while the Go language minimum remains 1.26.0. The selected `golang.org/x/net` version is v0.60.0; its required `golang.org/x/sys` version is v0.48.0. Ubuntu 26.04 images install the checksum-pinned `ca-certificates` 20260816ubuntu0.26.04.1 package with only officially fixed, already-installed packages advanced by the reviewed security overlay. These source inputs take effect in newly built artifacts; existing published release coordinates below describe their original artifacts.

Internal DNS serves metadata-derived service discovery records, client-specific search paths, recursive forwarding, and DNS response caching for PastureStack environments.

PastureStack is an independent community effort to preserve, audit, and modernize the Rancher 1.6 ecosystem. It is not affiliated with or endorsed by Rancher Labs or SUSE.

**Upstream:** [`rancher/rancher-dns`](https://github.com/rancher/rancher-dns). This GitHub fork retains the upstream Git history, authorship, dates, and license notices unchanged; PastureStack maintenance is consolidated into one commit after the preserved upstream boundary.

## Naming contract

- The internal service domain is `pasture.internal`.
- The metadata record is `metadata.pasture.internal`.
- Kubernetes service records remain under `svc.cluster.local`.
- Metadata-driven mode uses `--metadata-url`, which must include the API version, such as `http://metadata/2016-07-29`.
- Static metadata addresses use `--metadata-answer`.

The previous internal domain is intentionally not a public compatibility name. Upgrade tooling must perform an explicit, tested DNS cutover; see [COMPATIBILITY.md](COMPATIBILITY.md).

## Local checks

```sh
make test
make validate
make build
```

`scripts/build` also cross-compiles `windows/amd64`. The current published Linux
release is [v0.17.13](https://github.com/PastureStack/internal-dns/releases/tag/v0.17.13),
and the Catalog uses `ghcr.io/pasturestack/internal-dns:v0.17.13`.
The Windows Server 2022 image remains
`ghcr.io/pasturestack/internal-dns-windows:v0.17.12`; it was not rebuilt or
revalidated by this Linux release. That Windows coordinate is an exact manifest
alias of the already reviewed build
(`sha256:91d5147cc9ec20a8b9f8cc91745d620cc916b730645b8064a7c49f5f6a0d5c90`),
so the naming correction did not change image bytes.

Linux integration tests require `bats`, `dig`, free local TCP/UDP port 15353, and controlled recursive DNS access:

```sh
make integration-test
```

## Configuration file

JSON and YAML answer files may define client-specific and default records:

The Linux image reads and writes `/etc/internal-dns/answers.json` by default. Set `PLATFORM_DNS_ANSWERS_FILE` or pass an explicit `--answers` argument to use another path.

```yaml
default:
  search:
    - pasture.internal
  recurse:
    - 1.1.1.1
  authorative:
    - pasture.internal
  a:
    api.example.pasture.internal.:
      answer:
        - 10.0.0.20
```

## Security and license

Use only disposable DNS zones, metadata, and credentials in tests. See [SECURITY.md](SECURITY.md), [LICENSE](LICENSE), [ORIGIN.md](ORIGIN.md), and [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).

Security micro-updates are constrained by `security/ubuntu-security-packages.lock` to the signed Ubuntu 26.04 snapshot of 2026-10-23; the general package snapshot remains unchanged. Builds record the resulting package versions in `/usr/share/pasturestack/ubuntu-security-packages.tsv`. Existing runtime and firewall contracts are unchanged.
