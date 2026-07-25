# Third-Party Notices

Current Go dependencies and embedded code:

- `github.com/miekg/dns` v1.1.72 — BSD-3-Clause; see `LICENSES/miekg-dns.txt` and its author file.
- `gopkg.in/yaml.v3` v3.0.1 — MIT and Apache-2.0 portions; see `LICENSES/yaml-v3.txt`.
- The local `cache` package is derived from SkyDNS — MIT; see `cache/LICENSE`, `cache/AUTHORS`, and `cache/CONTRIBUTORS`.

The Linux container installs Ubuntu, CA certificates, curl, OpenSSL, and tini. Their licenses and notices remain governed by the Ubuntu packages and must be captured by the final image software bill of materials.

Historical vendored dependencies were removed from the current tree. Their original licenses and authors remain available in the preserved private source history and must be reviewed again if any code is restored.
